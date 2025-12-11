package success

import "errors"

type Message struct {
	Data   string `json:"data"`
	SentBy string
}

func NewMessage(data string, sentBy string) *Message {
	return &Message{
		Data:   data,
		SentBy: sentBy,
	}
}

func (m *Message) Validate() error {
	if m.Data == "" {
		return errors.New("message data cannot be empty")
	}

	if m.SentBy == "" {
		return errors.New("message sender cannot be empty")
	}

	return nil
}
