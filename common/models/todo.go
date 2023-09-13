package models

type TodoItem struct {
	Title string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool `json:"isCompleted"`
}