package interfaces

type User struct {
	ID    string
	Login string
}

type Session struct {
	User        User
	AccessToken string
}

type Credentials struct {
	Login    string
	Password string
}
