//go:build integration
// +build integration

package uoi

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawler_Next(t *testing.T) {
	c := NewCrawler(http.DefaultClient)

	for i := 0; i < 12; i++ {
		post, err := c.Next()
		assert.NoError(t, err)
		t.Log(post)
	}
}
