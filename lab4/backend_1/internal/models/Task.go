package models

type Task struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Status      string `json:"status"`
	UserID      uint   `json:"user_id"`
	Priority    string `json:"priority"`
}
