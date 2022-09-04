package uoi

import (
	"github.com/eoozs/befirst/model"
)

type Crawler struct {
}

func NewCrawler() *Crawler {
	return &Crawler{}
}

func (c *Crawler) Next() (*model.Post, error) {
	panic("todo")
}

func (c *Crawler) Reset() {}
