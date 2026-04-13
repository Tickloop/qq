package chat

import (
	"encoding/json"
	"testing"
)

func TestCompletionResponse_ValidJSON(t *testing.T) {
	raw := `{"id":"resp-1","choices":[{"index":0,"message":{"role":"assistant","content":"hello"},"finish_reason":"stop"}],"model":"test-model"}`

	var resp CompletionResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if resp.ID != "resp-1" {
		t.Errorf("ID = %q, want %q", resp.ID, "resp-1")
	}
	if resp.Model != "test-model" {
		t.Errorf("Model = %q, want %q", resp.Model, "test-model")
	}
	if len(resp.Choices) != 1 {
		t.Fatalf("len(Choices) = %d, want 1", len(resp.Choices))
	}

	choice := resp.Choices[0]
	if choice.Message.Role != "assistant" {
		t.Errorf("Role = %q, want %q", choice.Message.Role, "assistant")
	}
	if choice.Message.Content != "hello" {
		t.Errorf("Content = %q, want %q", choice.Message.Content, "hello")
	}
	if choice.FinishReason == nil || *choice.FinishReason != "stop" {
		t.Errorf("FinishReason = %v, want %q", choice.FinishReason, "stop")
	}
}

func TestCompletionResponse_EmptyChoices(t *testing.T) {
	raw := `{"id":"resp-2","choices":[],"model":"m"}`

	var resp CompletionResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(resp.Choices) != 0 {
		t.Errorf("len(Choices) = %d, want 0", len(resp.Choices))
	}
}

func TestCompletionResponse_NullContent(t *testing.T) {
	raw := `{"id":"resp-3","choices":[{"message":{"role":"assistant","content":null}}],"model":"m"}`

	var resp CompletionResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if resp.Choices[0].Message.Content != "" {
		t.Errorf("Content = %q, want empty string for null", resp.Choices[0].Message.Content)
	}
}

func TestCompletionResponse_MissingContent(t *testing.T) {
	raw := `{"id":"resp-4","choices":[{"message":{"role":"assistant"}}],"model":"m"}`

	var resp CompletionResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if resp.Choices[0].Message.Content != "" {
		t.Errorf("Content = %q, want empty string for missing field", resp.Choices[0].Message.Content)
	}
}

func TestCompletionResponse_MalformedJSON(t *testing.T) {
	raw := `{not json}`

	var resp CompletionResponse
	if err := json.Unmarshal([]byte(raw), &resp); err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}

func TestAPIError_Error(t *testing.T) {
	e := &APIError{StatusCode: 401, Body: `{"error":{"message":"Invalid API key"}}`}

	want := `api error (status 401): {"error":{"message":"Invalid API key"}}`
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
