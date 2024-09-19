package models

type User struct {
	ID   uint   `gorm:"primary_key"`
	Name string `json:"name"`
	Age  int    `json:"age" `
}
