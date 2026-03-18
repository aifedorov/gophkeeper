package certs

// Provider provides paths to TLS certificates.
type Provider struct {
	certPath string
	keyPath  string
}

// NewProvider creates a new certificate path provider.
func NewProvider(certPath, keyPath string) *Provider {
	return &Provider{
		certPath: certPath,
		keyPath:  keyPath,
	}
}

// CertPath returns the path to the server certificate file.
func (p *Provider) CertPath() string {
	return p.certPath
}

// KeyPath returns the path to the server private key file.
func (p *Provider) KeyPath() string {
	return p.keyPath
}
