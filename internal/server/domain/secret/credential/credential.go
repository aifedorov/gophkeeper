// Package credential provides credential domain entities.
package credential

// Credential represents a credential entity in the credential domain.
// It contains login/password pairs along with metadata, all of which are encrypted before storage.
type Credential struct {
	id       string
	userID   string
	name     string
	login    string
	password string
	notes    string
}

// NewCredential creates a new Credential entity with the provided data.
// It validates that all required fields (id, name, login, password) are not empty.
// Returns an error if validation fails.
func NewCredential(id, name, login, password, metadata string) (*Credential, error) {
	if id == "" {
		return nil, ErrIDRequired
	}
	if name == "" {
		return nil, ErrNameRequired
	}
	if login == "" {
		return nil, ErrLoginRequired
	}
	if password == "" {
		return nil, ErrPasswordRequired
	}

	return &Credential{
		id:       id,
		name:     name,
		login:    login,
		password: password,
		notes:    metadata,
	}, nil
}

// GetID returns the credential's unique identifier.
func (c *Credential) GetID() string {
	return c.id
}

// GetUserID returns the ID of the user who owns this credential.
func (c *Credential) GetUserID() string {
	return c.userID
}

// GetName returns the credential's display name for this credential.
func (c *Credential) GetName() string {
	return c.name
}

// GetLogin returns the decrypted login/username for this credential.
func (c *Credential) GetLogin() string {
	return c.login
}

// GetPassword returns the decrypted password for this credential.
func (c *Credential) GetPassword() string {
	return c.password
}

// GetMetadata returns the decrypted metadata/notes for this credential.
func (c *Credential) GetMetadata() string {
	return c.notes
}
