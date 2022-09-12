package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/eoozs/befirst/config"
	"github.com/eoozs/befirst/crawler/uoi"
	"github.com/eoozs/befirst/notify/telegram"
	"github.com/eoozs/befirst/pkg/env"
	"github.com/eoozs/befirst/service"
	"github.com/eoozs/befirst/storage"
)

const (
	storageFilePath = "/.befirst.cache"
)

type Config struct {
	HttpClientTimeout time.Duration `json:"httpClientTimeout"`
	SyncInterval      time.Duration `json:"syncInterval"`
}

func main() {
	var cfg Config
	if err := config.LoadFromFile(
		env.GetKey("CONFIG_DIR", "./config"),
		env.GetKey("CONFIG_NAME", "default"),
		&cfg); err != nil {
		log.Fatalf("load config: %s", err)
	}

	httpClient := http.DefaultClient
	httpClient.Timeout = cfg.HttpClientTimeout

	tgClient, err := initializeTelegramClient(httpClient)
	if err != nil {
		log.Fatalf("unable to initialize tg client: %v", err)
	}

	svc := service.NewBeFirst(
		map[string]service.PostsSource{
			"uoi": uoi.NewCrawler(httpClient),
		},
		map[string]service.Notifier{
			"telegram": tgClient,
		},
		cfg.SyncInterval,
	)

	wg := &sync.WaitGroup{}
	ctx, cf := context.WithCancel(context.Background())

	wg.Add(1)
	go svc.Run(ctx, wg)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	cf()
	wg.Wait()
}

func initializeTelegramClient(httpClient *http.Client) (*telegram.Client, error) {
	tgBotToken := os.Getenv("BEFIRST_TG_BOT_TOKEN")
	if len(tgBotToken) == 0 {
		return nil, fmt.Errorf("no token provided")
	}

	file, err := openOrCreateFile(storageFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}

	storageCl := storage.NewFSStorage(file)

	persistedChatIDs, err := storageCl.GetTelegramChatIDs()
	if err != nil {
		return nil, fmt.Errorf("unable to get storage chat IDs: %v", err)
	}

	tgClient := telegram.NewClient(telegram.NewClientInput{
		APIWrapper: telegram.NewRPCApiWrapper(
			"https://api.telegram.org", tgBotToken, httpClient,
		),
		InitialSubscribers: persistedChatIDs,
	})

	tgUpdates, err := tgClient.GetUpdates()
	if err != nil {
		return nil, fmt.Errorf("unable to get tg updates: %v", err)
	}

	for i := range tgUpdates {
		chatID := tgUpdates[i].Message.Chat.ID

		alreadySub, err := tgClient.HasSubscription(chatID)
		if err != nil {
			return nil, fmt.Errorf("unable to determine if chat subscribed: %v", err)
		}

		if alreadySub {
			continue
		}

		err = tgClient.AddSubscription(chatID)
		if err != nil {
			return nil, fmt.Errorf("unable to add chat subscription: %v", err)
		}

		err = storageCl.PutTelegramChatID(chatID)
		if err != nil {
			return nil, fmt.Errorf("unable to persist chat subscription: %v", err)
		}
	}

	return tgClient, nil
}

func openOrCreateFile(filepath string) (*os.File, error) {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_APPEND, os.ModeAppend)
	switch {
	case os.IsNotExist(err):
		file, err = os.Create(filepath)
		if err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	}

	return file, nil
}
