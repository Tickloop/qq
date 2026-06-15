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


func TestOpenRouterModelList(t *testing.T) {
	ctx := context.Background()
	modelList, err := chat.OpenRouterListModels(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range modelList {
		t.Logf("%s (%s)", model.Name, model.ID)
	}
}