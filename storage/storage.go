package storage

type Storage interface {
	PutTelegramChatID(int64) error
	GetTelegramChatIDs() ([]int64, error)
}
