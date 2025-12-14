package credential

type Credential struct {
	id       string
	userID   string
	name     string
	login    string
	password string
	notes    string
}

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

func (c *Credential) GetID() string {
	return c.id
}

func (c *Credential) GetUserID() string {
	return c.userID
}

func (c *Credential) GetName() string {
	return c.name
}

func (c *Credential) GetLogin() string {
	return c.login
}

func (c *Credential) GetPassword() string {
	return c.password
}

func (c *Credential) GetMetadata() string {
	return c.notes
}
