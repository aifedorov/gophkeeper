package storage

type User struct {
	ID    string
	Login string
}

type Session struct {
	User        User
	AccessToken string
}
