// Package interfaces provides authentication domain interfaces for the GophKeeper client.
package interfaces

// Session represents a user session containing authentication and encryption information.
// It includes the JWT access token, base64-encoded encryption key, and user ID.
type Session struct {
	accessToken   string // JWT access token for authenticated requests
	encryptionKey string // Base64-encoded encryption key for data encryption/decryption
	userID        string // Unique user identifier
}

// NewSession creates a new Session with the provided authentication data.
func NewSession(accessToken, encryptionKey, userID string) Session {
	return Session{
		accessToken:   accessToken,
		encryptionKey: encryptionKey,
		userID:        userID,
	}
}

// GetAccessToken returns the JWT access token for authenticated requests.
func (s Session) GetAccessToken() string {
	return s.accessToken
}

// GetEncryptionKey returns the base64-encoded encryption key for data encryption/decryption.
func (s Session) GetEncryptionKey() string {
	return s.encryptionKey
}

// GetUserID returns the unique user identifier.
func (s Session) GetUserID() string {
	return s.userID
}
