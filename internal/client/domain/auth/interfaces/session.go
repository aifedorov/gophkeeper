package interfaces

type Session struct {
	accessToken   string
	encryptionKey string // base64
	userID        string
}

func NewSession(accessToken, encryptionKey, userID string) Session {
	return Session{
		accessToken:   accessToken,
		encryptionKey: encryptionKey,
		userID:        userID,
	}
}

func (s Session) GetAccessToken() string {
	return s.accessToken
}

func (s Session) GetEncryptionKey() string {
	return s.encryptionKey
}

func (s Session) GetUserID() string {
	return s.userID
}
