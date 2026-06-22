package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// hTTPClient implements Provider using net/http against an OpenAI-compatible API.
type hTTPClient struct {
	baseURL    string
	apiKey     string
	HttpClient *http.Client
}

// NewHTTPClient constructs an HTTPClient.
// baseURL is the API root (e.g., "https://openrouter.ai/api/v1").
// Accepting baseURL enables testing with httptest.
// If httpClient is nil, http.DefaultClient is used.
func newHTTPClient(baseURL, apiKey string, httpClient *http.Client) *hTTPClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &hTTPClient{
		baseURL:    baseURL,
		apiKey:     apiKey,
		HttpClient: httpClient,
	}
}

// Complete sends a non-streaming chat completion request and returns the parsed response.
func (c *hTTPClient) complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	var result CompletionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}

func (c *hTTPClient) listModels(ctx context.Context) ([]Model, error) {
	var models = []Model{}
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/models", http.NoBody)
	if err != nil {
		return models, err
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return models, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models, err
	}

	var modelsResponse ListModelsResponse
	if err := json.Unmarshal(respBody, &modelsResponse); err != nil {
		return models, err
	}
	for _, model := range modelsResponse.Data {
		models = append(models, Model{ID: model.ID, Name: model.Name})
	}

	return models, nil
}

func checkAndGetOpenRouterAPIKey() (string, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("error: OPENROUTER_API_KEY is not set")
	}
	return apiKey, nil
}

func OpenRouterConverse(ctx context.Context, question string, modelId string) (string, error) {
	apiKey, err := checkAndGetOpenRouterAPIKey()
	if err != nil {
		return "", err
	}
	ctxWindow := NewContextWindow()
	ctxWindow.AddMessage(NewMessage("user", question))
	client := newHTTPClient("https://openrouter.ai/api/v1", apiKey, nil)
	req := CompletionRequest{
		Model:    modelId,
		Messages: ctxWindow.Messages,
	}
	resp, err := client.complete(ctx, req)
	if err != nil {
		return "", fmt.Errorf("error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("error: no response from model")
	}

	text := strings.TrimSpace(resp.Choices[0].Message.Content)
	if text == "" {
		return "", fmt.Errorf("error: empty response from model")
	}
	ctxWindow.AddMessage(NewMessage("assistant", text))
	ctxWindow.Save()
	return text, nil
}

func OpenRouterListModels(ctx context.Context) ([]Model, error) {
	client := newHTTPClient("https://openrouter.ai/api/v1", "", nil)
	models, err := client.listModels(ctx)
	if err != nil {
		return models, err
	}
	return models, nil
}
