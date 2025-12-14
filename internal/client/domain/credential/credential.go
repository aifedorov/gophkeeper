package credential

type Credential struct {
	ID       string
	Name     string
	Login    string
	Password string
	Notes    string
}

func NewCredential(id, name, login, password, notes string) (*Credential, error) {
	return &Credential{
		ID:       id,
		Name:     name,
		Login:    login,
		Password: password,
		Notes:    notes,
	}, nil
}

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
