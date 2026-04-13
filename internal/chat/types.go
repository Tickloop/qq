package chat

import "fmt"

// Message represents a single chat message in a request.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionRequest is the POST body for /chat/completions.
type CompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// CompletionResponse is the top-level response from /chat/completions.
type CompletionResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Model   string   `json:"model"`
}

// Choice is a single completion choice within a response.
type Choice struct {
	Index        int           `json:"index"`
	Message      ChoiceMessage `json:"message"`
	FinishReason *string       `json:"finish_reason"`
}

// ChoiceMessage is the assistant's response message within a choice.
type ChoiceMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// APIError represents a non-2xx response from the API.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error (status %d): %s", e.StatusCode, e.Body)
}
