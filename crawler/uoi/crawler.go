package uoi

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/eoozs/befirst/model"
	"github.com/eoozs/befirst/pkg/net"
)

var ErrEnd = errors.New("end")

type Crawler struct {
	httpClient        net.HTTPClient
	urlFmt            string
	page              int
	currentPagePosts  []model.Post
	currentPageOffset int
}

func NewCrawler(httpClient net.HTTPClient) *Crawler {
	return &Crawler{
		httpClient:        httpClient,
		urlFmt:            "http://e-enoikiazetai.uoi.gr/show.php?offset={{page}}&dd_type=0&sort=0",
		page:              -1,
		currentPagePosts:  make([]model.Post, 0, 10),
		currentPageOffset: 0,
	}
}

func (c *Crawler) Next() (*model.Post, error) {
	if len(c.currentPagePosts) == 0 || c.currentPageOffset >= len(c.currentPagePosts) {
		if err := c.crawlNextPage(); err != nil {
			return nil, fmt.Errorf("fetch next page: %w", err)
		}
	}

	defer func() { c.currentPageOffset++ }()
	if c.currentPageOffset >= len(c.currentPagePosts) {
		return nil, ErrEnd
	}
	return &c.currentPagePosts[c.currentPageOffset], nil
}

func (c *Crawler) Reset() {
	c.page = -1
	c.currentPagePosts = make([]model.Post, 0, 10)
	c.currentPageOffset = 0
}

func (c *Crawler) crawlNextPage() error {
	c.page++
	url := strings.ReplaceAll(c.urlFmt, "{{page}}", strconv.Itoa(c.page))

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("get %s: %w", url, err)
	}

	if err := c.parsePostsFromRespBody(resp.Body); err != nil {
		return err
	}

	c.currentPageOffset = 0
	return nil
}

func (c *Crawler) parsePostsFromRespBody(r io.Reader) error {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return fmt.Errorf("parse html body: %w", err)
	}

	currentPagePosts := make([]model.Post, 10)

	doc.Find("tbody>tr>td>div>table>tbody>tr>td>font>b").Each(func(i int, selection *goquery.Selection) {
		if i < 0 || i > len(currentPagePosts) {
			return
		}
		t := selection.Text()
		if i%2 == 0 {
			currentPagePosts[i/2].Title = t
			currentPagePosts[i/2].ID = strings.TrimSpace(strings.TrimPrefix(t, "Αριθμ. Αγγελίας:"))
		} else {
			currentPagePosts[(i-1)/2].Title += " " + t // date
		}
	})

	doc.Find("tbody>tr>td  div[align=left].smalls").Each(func(i int, selection *goquery.Selection) {
		if i >= 0 && i < len(currentPagePosts) {
			currentPagePosts[i].Description = strings.Join(strings.Fields(selection.Text()), " ")
		}
	})

	c.currentPagePosts = make([]model.Post, 0, 10)
	for _, p := range currentPagePosts {
		if p.ID != "" {
			c.currentPagePosts = append(c.currentPagePosts, p)
		}
	}

	return err
}
