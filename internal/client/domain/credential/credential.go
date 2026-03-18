// Package credential provides credential domain entities for the GophKeeper client.
package credential

// Credential represents a credential entity in the client domain.
// It contains login/password pairs along with metadata.
type Credential struct {
	ID       string // Unique credential identifier
	Name     string // Display name (e.g., "Gmail Account")
	Login    string // Username or email
	Password string // Password
	Notes    string // Optional notes/metadata
	Version  int64  // Version number
}

// NewCredential creates a new Credential entity with the provided data.
// Returns an error if validation fails (currently always returns nil as validation is done in Validate method).
func NewCredential(id, name, login, password, notes string, version int64) (*Credential, error) {
	return &Credential{
		ID:       id,
		Name:     name,
		Login:    login,
		Password: password,
		Notes:    notes,
		Version:  version,
	}, nil
}

// Validate checks that all required fields (ID, Name, Login, Password) are not empty.
// Returns an error if any required field is missing.
func (c *Credential) Validate() error {
	if c.ID == "" {
		return ErrIDRequired
	}
	if c.Name == "" {
		return ErrNameRequired
	}
	if c.Login == "" {
		return ErrLoginRequired
	}
	if c.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}
