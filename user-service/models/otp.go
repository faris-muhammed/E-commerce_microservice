package models

import "time"

type OTP struct {
	Id        uint
	Email  string
	OTP       string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type CachedUser struct {
	Name     string
	Email    string
	Password string
	Mobile   string
}