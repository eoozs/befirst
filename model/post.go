package model

type Post struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

func (p Post) String() string {
	return "todo"
}
