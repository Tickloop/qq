package chat

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

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
	return chatDir, nil
}

func getChatFile() (string, error) {
	newChatUUID, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	chatFile := newChatUUID.String() + ".json"
	return chatFile, nil
}

func ListChats() ([]string, error) {
	chatDir, err := getChatDir()
	if err != nil {
		return nil, err
	}
	matches, err := 
}

func NewMessage(role, message string) Message {
	return Message{
		Role:    role,
		Content: message,
	}
}

type ContextWindow struct {
	Messages []Message
}

func NewContextWindow() ContextWindow {
	return ContextWindow{
		Messages: make([]Message, 0),
	}
}

func (c *ContextWindow) AddMessage(msg Message) {
	c.Messages = append(c.Messages, msg)
}

func (c *ContextWindow) Reset() {
	c.Messages = []Message{}
}

func (c *ContextWindow) Save() error {
	chatDir, err := getChatDir()
	chatFileName, err := getChatFile()
	chatFile := filepath.Join(chatDir, chatFileName)
	if err != nil {
		return err
	}

	data, err := json.Marshal(c.Messages)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(chatDir, 0o700); err != nil {
		return err
	}

	if err := os.WriteFile(chatFile, data, 0o600); err != nil {
		return err
	}

	return nil
}
