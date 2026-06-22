package chat

type ContentBlock string
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
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
