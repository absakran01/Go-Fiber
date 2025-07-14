package main

type Status string

const(
	Read Status = "read"
	Reading Status = "reading"
	ToRead Status = "to_read"
)

type User struct {
	Id       int    `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"username"`
	Password string `json:"-"`
}

type Book struct{
	Id       int    `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
	Status Status `json:"status" gorm:"default:to_read"`
	UserId int `json:"user_id"`
}
