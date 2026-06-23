package chat_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/tickloop/qq/internal/chat"
)

func TestOpenRouterConverse(t *testing.T) {
	tmpChatDir, err := os.MkdirTemp(".", "test-*")
	os.Setenv("QQ_CHAT_DIR", tmpChatDir)

	ctx := context.Background()
	question := "hello! What is your name?"
	modelId := "perplexity/sonar"
	answer, err := chat.OpenRouterConverse(ctx, question, modelId)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("question: %s", question)
	t.Logf("(%s) answer: %s", modelId, answer)

	// also test chat persistance
	chats, err := chat.ListChats()
	if err != nil {
		t.Fatal(err)
	}

	if len(chats) == 0 {
		t.Fatal("error: found 0 chats persisted")
	}

	if len(chats[0].Messages) == 0 {
		t.Fatal("error: user message not persisted")
	}

	if chats[0].Messages[0].Role != "user" {
		t.Fatal("error: message role not saved as user")
	}

	if chats[0].Messages[0].Content != question {
		t.Fatalf("error: question not saved as input\nexp:%s\ngot:%s", question, chats[0].Messages[0].Content)
	}

	for _, chat := range chats {
		t.Logf("question: %s", chat.Messages[0].Content)
		t.Logf("id: %s", chat.Id)
	}
	os.RemoveAll(tmpChatDir)
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

func TestAddMessageCreateDir(t *testing.T) {
	testDirPath := "./test-chat-dir-created"
	os.Setenv("QQ_CHAT_DIR", testDirPath)

	// a new dir and a new file should be created
	cw := chat.NewContextWindow()
	cw.AddMessage(chat.NewMessage("user", "hello"))

	if _, err := os.Stat(testDirPath); err != nil {
		t.Log("new dir not created")
		t.Fatal(err)
	}

	fileName := cw.Id.String() + ".jsonl"
	filePath := filepath.Join(testDirPath, fileName)
	if _, err := os.Stat(filePath); err != nil {
		t.Log("new chat file not created")
		t.Fatal(err)
	}

	os.RemoveAll(testDirPath)
}
