package client

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Dialer *websocket.Dialer
	Conn   *websocket.Conn
}

func NewClient(conn *websocket.Dialer) *Client {
	return &Client{
		Dialer: conn,
		Conn:   nil,
	}
}

func (c *Client) Connect(urlStr string) error {
	conn, _, err := c.Dialer.Dial(urlStr, nil)

	if err != nil {
		return err
	}

	c.Conn = conn

	return nil
}

func (c *Client) SendMessage(messageType int, message []byte) error {
	if c.Conn != nil {
		err := c.Conn.WriteMessage(messageType, message)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) Close() error {
	if c.Conn != nil {
		err := c.Conn.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
