package chat_test

import (
	"context"
	"testing"

	"github.com/tickloop/qq/internal/chat"
)

func TestAWSHello(t *testing.T) {
	ctx := context.Background()
	question := "Hi! What is your name?"
	modelId := "global.anthropic.claude-opus-4-8"
	answer, err := chat.AWSConverse(ctx, question, modelId)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("question: %s", question)
	t.Logf("(%s) answer: %s", modelId, answer)
}


func TestAWSListModels(t *testing.T) {
	ctx := context.Background()
	models, err := chat.AWSListModels(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range models {
		t.Logf("%s (%s)", model.Name, model.ID)
	}
}