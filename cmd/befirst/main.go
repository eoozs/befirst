package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/eoozs/befirst/crawler/uoi"
	"github.com/eoozs/befirst/notify/telegram"
	"github.com/eoozs/befirst/service"
)

func main() {
	httpClient := http.DefaultClient
	httpClient.Timeout = 10 * time.Second

	syncInterval := 10 * time.Second

	svc := service.NewBeFirst(
		map[string]service.PostsSource{
			"uoi": uoi.NewCrawler(httpClient),
		},
		map[string]service.Notifier{
			"telegram": telegram.NewClient(),
		},
		syncInterval,
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
