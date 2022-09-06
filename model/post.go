package model

import "fmt"

type Post struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

func (p Post) String() string {
	return fmt.Sprintf("[%s] %s %s %s", p.ID, p.Title, p.URL, p.Description)
}
