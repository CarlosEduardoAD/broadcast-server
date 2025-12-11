package message

import "errors"

type Message struct {
	Data string `json:"data"`
}

func NewMessage(data string) *Message {
	return &Message{
		Data: data,
	}
}

func (m *Message) Validate() error {
	if m.Data == "" {
		return errors.New("message data cannot be empty")
	}

	return nil
}
