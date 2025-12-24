package cache

type Secret struct {
	Credentials map[string]int64 `json:"credentials"`
	Cards       map[string]int64 `json:"cards"`
	Files       map[string]int64 `json:"files"`
}

func NewSecret() *Secret {
	return &Secret{
		Credentials: make(map[string]int64),
		Cards:       make(map[string]int64),
		Files:       make(map[string]int64),
	}
}

func (s *Secret) SetCredential(id string, version int64) *Secret {
	s.Credentials[id] = version
	return s
}

func (s *Secret) GetCredentialVersion(id string) (int64, bool) {
	v, ok := s.Credentials[id]
	if !ok {
		return 0, false
	}
	return v, true
}

func (s *Secret) DeleteCredential(id string) *Secret {
	delete(s.Credentials, id)
	return s
}

func (s *Secret) SetCard(id string, version int64) *Secret {
	s.Cards[id] = version
	return s
}

func (s *Secret) GetCardVersion(id string) (int64, bool) {
	v, ok := s.Cards[id]
	if !ok {
		return 0, false
	}
	return v, true
}

func (s *Secret) DeleteCard(id string) *Secret {
	delete(s.Cards, id)
	return s
}

func (s *Secret) SetFileVersion(id string, version int64) *Secret {
	s.Files[id] = version
	return s
}

func (s *Secret) GetFileVersion(id string) (int64, bool) {
	v, ok := s.Files[id]
	if !ok {
		return 0, false
	}
	return v, true
}

func (s *Secret) DeleteFileVersion(id string) *Secret {
	delete(s.Files, id)
	return s
}
