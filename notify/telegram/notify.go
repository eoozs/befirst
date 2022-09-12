package telegram

import "github.com/eoozs/befirst/pkg/set"

type Storage interface {
	PutTelegramChatID(int64) error
	GetTelegramChatIDs() ([]int64, error)
}

type Client struct {
	subscribedChats set.Set[int64]
	apiClient       APIWrapper
	storage 		Storage
}

type NewClientInput struct {
	APIWrapper         APIWrapper
}

func NewClient(apiWrapper APIWrapper, storage Storage) *Client {
	return &Client{
		subscribedChats: 	*set.New[int64](),
		apiClient:       	input.APIWrapper,
		storage:			storage,
	}
}

func (c *Client) Sync() error {
	persistedChatIDs, err := c.storage.GetTelegramChatIDs()
	if err != nil {
		return nil, fmt.Errorf("unable to get storage chat IDs: %v", err)
	}

	tgUpdates, err := c.GetUpdates()
	if err != nil {
		return nil, fmt.Errorf("unable to get tg updates: %v", err)
	}

	for i := range tgUpdates {
		chatID := tgUpdates[i].Message.Chat.ID

		alreadySub, err := c.HasSubscription(chatID)
		if err != nil {
			return nil, fmt.Errorf("unable to determine if chat subscribed: %v", err)
		}

		if alreadySub {
			continue
		}

		err = c.AddSubscription(chatID)
		if err != nil {
			return nil, fmt.Errorf("unable to add chat subscription: %v", err)
		}

		err = c.storage.PutTelegramChatID(chatID)
		if err != nil {
			return nil, fmt.Errorf("unable to persist chat subscription: %v", err)
		}
	}
	return nil
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
