package telegram

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) SendMessage(msg string) error {
	panic("todo")
}
