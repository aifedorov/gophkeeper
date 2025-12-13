package storage

type Session struct {
	AccessToken   string `json:"access_token"`
	EncryptionKey string `json:"encryption_key"` // base64
}

func NewSession(accessToken, encryptionKey string) Session {
	return Session{AccessToken: accessToken, EncryptionKey: encryptionKey}
}

func (s Session) GetAccessToken() string {
	return s.AccessToken
}

func (s Session) GetEncryptionKey() string {
	return s.EncryptionKey
}
