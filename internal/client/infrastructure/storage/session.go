package storage

type Session struct {
	AccessToken   string `json:"access_token"`
	EncryptionKey string `json:"encryption_key"` // base64
	UserID        string `json:"user_id"`
}

func NewSession(accessToken, encryptionKey, userID string) Session {
	return Session{AccessToken: accessToken, EncryptionKey: encryptionKey, UserID: userID}
}

func (s Session) GetAccessToken() string {
	return s.AccessToken
}

func (s Session) GetEncryptionKey() string {
	return s.EncryptionKey
}

func (s Session) GetUserID() string {
	return s.UserID
}
