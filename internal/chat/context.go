package chat

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/tickloop/qq/internal/utils"

	"github.com/google/uuid"
)

var dbg = utils.Dbg

type ContentBlock string

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func getChatDir() (string, error) {
	base, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	chatDir := filepath.Join(base, ".qq/sessions")
	if envChatDir, ok := os.LookupEnv("QQ_CHAT_DIR"); ok && envChatDir != "" {
		if envChatDir, err = utils.ResolveToAbsPath(envChatDir); err == nil {
			chatDir = envChatDir
		}
	}
	return chatDir, nil
}

func NewMessage(role, message string) Message {
	return Message{
		Role:    role,
		Content: message,
	}
}

type ContextWindow struct {
	Id       uuid.UUID
	Messages []Message
}

func NewContextWindow() ContextWindow {
	id, _ := uuid.NewV7()
	return ContextWindow{
		Id:       id,
		Messages: make([]Message, 0),
	}
}

func (c *ContextWindow) saveMessageToChatHistory(msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	data = append(data, '\n')

	chatDir, err := getChatDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(chatDir, 0o700); err != nil {
		return err
	}

	name := filepath.Join(chatDir, c.Id.String()+".jsonl")
	fd, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()

	if _, err := fd.Write(data); err != nil {
		return err
	}
	return nil
}

func (c *ContextWindow) AddMessage(msg Message) {
	c.Messages = append(c.Messages, msg)
	if err := c.saveMessageToChatHistory(msg); err != nil {
		dbg("%s", err)
	}
}

func (c *ContextWindow) Reset() {
	c.Id, _ = uuid.NewV7()
	c.Messages = []Message{}
}
