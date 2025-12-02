package memory

import "time"

type User struct {
	ID    string
	Login string
}

type Session struct {
	User        User
	AccessToken string
	ExpiresAt   time.Time
}
