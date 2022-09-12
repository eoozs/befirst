package storage

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type FSStorage struct {
	file io.ReadWriteSeeker
}

func NewFSStorage(file io.ReadWriteSeeker) *FSStorage {
	return &FSStorage{
		file: file,
	}
}

func (s *FSStorage) PutTelegramChatID(chatID int64) error {
	_, err := s.file.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("unable to seek: %v", err)
	}
	if _, err = fmt.Fprintln(s.file, strconv.FormatInt(chatID, 10)); err != nil {
		return fmt.Errorf("unable to write: %v", err)
	}

	return nil
}

func (s *FSStorage) GetTelegramChatIDs() ([]int64, error) {
	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("unable to seek: %v", err)
	}

	b, err := ioutil.ReadAll(s.file)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %v", err)
	}

	chatIDs := make([]int64, 0)
	for _, v := range strings.Fields(string(b)) {
		chatID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse chatID (%s): %v", v, err)
		}
		chatIDs = append(chatIDs, chatID)
	}

	return chatIDs, nil
}
