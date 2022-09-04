package main

import (
	"context"
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
	svc := service.NewBeFirst(
		map[string]service.PostsSource{
			"uoi": uoi.NewCrawler(),
		},
		map[string]service.Notifier{
			"telegram": telegram.NewClient(),
		},
		10*time.Second,
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
