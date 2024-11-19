package model

type Post struct {
	PostID    string `json:"post_id"`
	Email     string `json:"email"`
	ImageURL  string `json:"image_url"`
	IsVisible bool   `json:"is_visible"`
	CreatedAt string `json:"created_at"`
}
