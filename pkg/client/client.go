package client

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	Dialer Dialer
	Conn   *websocket.Conn
}

func NewClient(d Dialer) *Client {
	return &Client{
		Dialer: d,
		Conn:   nil,
	}

}

func (c *Client) Connect(urlStr string, request_headers http.Header) error {
	conn, _, err := c.Dialer.Dial(urlStr, request_headers)

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
