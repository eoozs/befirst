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



type Config struct {
	HttpClientTimeout 	time.Duration 	`json:"httpClientTimeout"`
	SyncInterval      	time.Duration 	`json:"syncInterval"`
	StorageFilePath   	string 			`json:"storageFilePath"`
	TelegramBotAPIUrl 	string			`json:"telegramBotAPIUrl"`
	TelegramBotToken 	string			`json:"telegramBotToken"`
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

func initializeTelegramClient(httpClient *http.Client, conf Config) (*telegram.Client, error) {
	if len(conf.TelegramBotToken) == 0 {
		return nil, fmt.Errorf("no token provided")
	}

	file, err := openOrCreateFile(conf.StorageFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}

	storageCl := storage.NewFSStorage(file)

	tgClient := telegram.NewClient(telegram.NewClientInput{
		APIWrapper: telegram.NewRPCApiWrapper(
			conf.TelegramBotAPIUrl, TelegramBotToken, httpClient,
		),
	})

	err = tgClient.Sync()
	if err != nil {
		return nil, fmt.Errorf("unable to call tg sync method: %v", err)
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
