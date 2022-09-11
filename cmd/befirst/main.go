package main

import (
	"context"
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

	svc := service.NewBeFirst(
		map[string]service.PostsSource{
			"uoi": uoi.NewCrawler(httpClient),
		},
		map[string]service.Notifier{
			"telegram": telegram.NewClient(),
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
