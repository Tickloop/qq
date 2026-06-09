package chat_test

import (
	"context"
	"testing"

	"github.com/tickloop/qq/internal/chat"
)

func TestOpenRouterConverse(t *testing.T) {
	ctx := context.Background()
	question := "hello! What is your name?"
	modelId := "perplexity/sonar"
	answer, err := chat.OpenRouterConverse(ctx, question, modelId)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("question: %s", question)
	t.Logf("(%s) answer: %s", modelId, answer)
}