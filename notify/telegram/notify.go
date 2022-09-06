package telegram

import "fmt"

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) SendMessage(msg string) error {
	fmt.Printf("TODO: Telegram Send: %s\n", msg)
	return nil
}
