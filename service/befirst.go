package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/eoozs/befirst/model"
)

type PostsSource interface {
	Next() (*model.Post, error)
	Reset()
}

type Notifier interface {
	SendMessage(msg string) error
}

type BeFirst struct {
	postSources           map[string]PostsSource
	postSourcesLastPostID map[string]string
	notifiers             map[string]Notifier
	interval              time.Duration
}

func NewBeFirst(
	postSources map[string]PostsSource, notifiers map[string]Notifier, pollingInterval time.Duration,
) *BeFirst {
	return &BeFirst{
		postSources:           postSources,
		postSourcesLastPostID: make(map[string]string),
		notifiers:             notifiers,
		interval:              pollingInterval,
	}
}

func (b *BeFirst) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(b.interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("befirst service gracefully stopped")
			return
		case <-ticker.C:
			log.Println("checking for new posts")
			b.run(ctx)
		}
	}
}

func (b *BeFirst) run(ctx context.Context) {
	wg := sync.WaitGroup{}

	for sourceID, source := range b.postSources {
		wg.Add(1)
		go func(sourceID string, source PostsSource) {
			defer wg.Done()
			err := b.runForPostsSource(ctx, sourceID, source)
			if err != nil {
				log.Printf("[error] %s: %s", sourceID, err)
			}
		}(sourceID, source)
	}

	wg.Wait()
}

func (b *BeFirst) runForPostsSource(ctx context.Context, postsSourceID string, ps PostsSource) error {
	mostRecentPostID := ""

	for {
		select {
		case <-ctx.Done():
			ps.Reset()
			return nil
		default:
			post, err := ps.Next()
			if err != nil {
				return err
			}

			if mostRecentPostID == "" {
				mostRecentPostID = post.ID
			}

			lastPostID := b.postSourcesLastPostID[postsSourceID]
			if post.ID == lastPostID || lastPostID == "" {
				b.postSourcesLastPostID[postsSourceID] = mostRecentPostID
				break
			}

			if err := b.notify(*post); err != nil {
				return err
			}
		}
	}
}

func (b *BeFirst) notify(post model.Post) error {
	for notifierID, notifier := range b.notifiers {
		log.Printf("notifying %s for post %s\n", notifierID, post.ID)
		if err := notifier.SendMessage(post.String()); err != nil {
			return err
		}
	}

	return nil
}
