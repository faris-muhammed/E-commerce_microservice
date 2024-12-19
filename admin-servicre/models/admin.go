package models

type Admin struct {
	Id       int
	Username string `gorm:"unique"`
	Password string
}
