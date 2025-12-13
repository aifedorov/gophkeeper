package interfaces

type Session struct {
	accessToken   string
	encryptionKey string // base64
}

func NewSession(accessToken, encryptionKey string) Session {
	return Session{
		accessToken:   accessToken,
		encryptionKey: encryptionKey,
	}
}

func (s Session) GetAccessToken() string {
	return s.accessToken
}

func (s Session) GetEncryptionKey() string {
	return s.encryptionKey
}
