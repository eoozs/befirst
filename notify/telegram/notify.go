package telegram

import "github.com/eoozs/befirst/pkg/set"

type Client struct {
	subscribedChats set.Set[int64]
	apiClient       APIWrapper
}

type NewClientInput struct {
	APIWrapper         APIWrapper
	InitialSubscribers []int64
}

func NewClient(input NewClientInput) *Client {
	initSubsSet := set.New[int64]()
	initSubsSet.Add(input.InitialSubscribers...)

	return &Client{
		subscribedChats: *initSubsSet,
		apiClient:       input.APIWrapper,
	}
}

func (c *Client) HasSubscription(chatID int64) (bool, error) {
	return c.subscribedChats.Contains(chatID), nil
}

func (c *Client) AddSubscription(chatID int64) error {
	c.subscribedChats.Add(chatID)
	return nil
}

func (c *Client) SendMessage(msg string) error {
	for _, chatID := range c.subscribedChats.ToSlice() {
		err := c.apiClient.SendMessage(chatID, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) GetUpdates() ([]Update, error) {
	return c.apiClient.GetUpdates()
}
