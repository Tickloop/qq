package chat

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPClient_Complete_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CompletionResponse{
			ID:    "resp-1",
			Model: "test-model",
			Choices: []Choice{
				{
					Index:   0,
					Message: ChoiceMessage{Role: "assistant", Content: "42"},
				},
			},
		})
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, "test-key", srv.Client())
	resp, err := client.Complete(context.Background(), CompletionRequest{
		Model:    "test-model",
		Messages: []Message{{Role: "user", Content: "what is 6*7"}},
	})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}

	if len(resp.Choices) != 1 {
		t.Fatalf("len(Choices) = %d, want 1", len(resp.Choices))
	}
	if resp.Choices[0].Message.Content != "42" {
		t.Errorf("Content = %q, want %q", resp.Choices[0].Message.Content, "42")
	}
}

func TestHTTPClient_Complete_RequestFormat(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify HTTP method and path
		if r.Method != http.MethodPost {
			t.Errorf("Method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/chat/completions" {
			t.Errorf("Path = %q, want /chat/completions", r.URL.Path)
		}

		// Verify headers
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("Authorization = %q, want %q", got, "Bearer test-key")
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want %q", got, "application/json")
		}

		// Verify body structure
		body, _ := io.ReadAll(r.Body)
		var req CompletionRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("unmarshal request body: %v", err)
		}
		if req.Model != "my-model" {
			t.Errorf("req.Model = %q, want %q", req.Model, "my-model")
		}
		if len(req.Messages) != 1 || req.Messages[0].Content != "hi" {
			t.Errorf("req.Messages = %+v, want [{Role:user Content:hi}]", req.Messages)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CompletionResponse{
			Choices: []Choice{{Message: ChoiceMessage{Content: "hello"}}},
		})
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, "test-key", srv.Client())
	_, err := client.Complete(context.Background(), CompletionRequest{
		Model:    "my-model",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
}

func TestHTTPClient_Complete_APIError401(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":{"message":"Invalid API key"}}`))
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, "bad-key", srv.Client())
	_, err := client.Complete(context.Background(), CompletionRequest{
		Model:    "m",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("StatusCode = %d, want 401", apiErr.StatusCode)
	}
}

func TestHTTPClient_Complete_APIError500(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, "key", srv.Client())
	_, err := client.Complete(context.Background(), CompletionRequest{
		Model:    "m",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("StatusCode = %d, want 500", apiErr.StatusCode)
	}
}

func TestHTTPClient_Complete_MalformedResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{not json}`))
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, "key", srv.Client())
	_, err := client.Complete(context.Background(), CompletionRequest{
		Model:    "m",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}

func TestHTTPClient_Complete_ContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Delay long enough for the context to be cancelled
		time.Sleep(2 * time.Second)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	client := NewHTTPClient(srv.URL, "key", srv.Client())
	_, err := client.Complete(ctx, CompletionRequest{
		Model:    "m",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}

func TestNewHTTPClient_DefaultsToDefaultClient(t *testing.T) {
	client := NewHTTPClient("http://example.com", "key", nil)
	if client.httpClient != http.DefaultClient {
		t.Error("expected http.DefaultClient when nil is passed")
	}
}
