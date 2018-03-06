package auth

import (
	"errors"
	"fmt"

	"github.com/VolantMQ/vlauth"
)

// Manager auth
type Manager struct {
	p []vlauth.Iface
}

var providers = make(map[string]vlauth.Iface)

// Register auth provider
func Register(name string, i vlauth.Iface) error {
	if name == "" && i == nil {
		return errors.New("invalid args")
	}

	if _, dup := providers[name]; dup {
		return errors.New("already exists")
	}

	providers[name] = i

	return nil
}

// UnRegister authenticator
func UnRegister(name string) {
	delete(providers, name)
}

// NewManager new auth manager
func NewManager(p []string) (*Manager, error) {
	m := Manager{}

	for _, pa := range p {
		pvd, ok := providers[pa]
		if !ok {
			return nil, fmt.Errorf("session: unknown provider %q", pa)
		}

		m.p = append(m.p, pvd)
	}

	return &m, nil
}

// Password authentication
func (m *Manager) Password(user, password string) error {
	for _, p := range m.p {
		if status := p.Password(user, password); status == vlauth.StatusAllow {
			return status
		}
	}

	return vlauth.StatusDeny
}

// ACL check permissions
func (m *Manager) ACL(clientID, user, topic string, access vlauth.AccessType) error {
	for _, p := range m.p {
		if status := p.ACL(clientID, user, topic, access); status == vlauth.StatusAllow {
			return status
		}
	}

	return vlauth.StatusDeny
}
