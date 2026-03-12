package credentials

import "errors"

// ErrNotConfigured is returned when no credentials are stored for the requested profile.
var ErrNotConfigured = errors.New("credentials not configured")

// Credentials holds the API key, token, and auth mode for a profile.
type Credentials struct {
	APIKey   string `json:"apiKey"`
	Token    string `json:"token"`
	AuthMode string `json:"authMode"`
}

// Store defines the interface for credential persistence.
type Store interface {
	Get(profile string) (Credentials, error)
	Set(profile string, creds Credentials) error
	Delete(profile string) error
}

// MemoryStore is an in-memory Store implementation, primarily for testing.
type MemoryStore struct {
	creds map[string]Credentials
}

// NewMemoryStore creates a new in-memory credential store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{creds: make(map[string]Credentials)}
}

func (m *MemoryStore) Get(profile string) (Credentials, error) {
	cred, ok := m.creds[profile]
	if !ok {
		return Credentials{}, ErrNotConfigured
	}
	return cred, nil
}

func (m *MemoryStore) Set(profile string, creds Credentials) error {
	m.creds[profile] = creds
	return nil
}

func (m *MemoryStore) Delete(profile string) error {
	delete(m.creds, profile)
	return nil
}
