package chat_test

import (
	"testing"

	"github.com/tickloop/qq/internal/chat"
)

func TestAWSHello(t *testing.T) {
	question := "Hi! What is your name?"
	modelId := "anthropic.claude-3-haiku-20240307-v1:0"
	_, err := chat.AWSConverse(question, modelId)
	if err != nil {
		t.Fatal(err)
	}
}
