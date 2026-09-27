package sessions

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

type State string

const (
	StatePending State = "PENDING_NODE_APPLY"
	StateReady   State = "READY"
	StateFailed  State = "FAILED"
)

type Session struct {
	ID              string    `json:"id"`
	CommandID       string    `json:"commandId"`
	NodeID          string    `json:"nodeId"`
	ClientPublicKey string    `json:"-"`
	ClientAddress   string    `json:"clientAddress"`
	Endpoint        string    `json:"endpoint"`
	ServerPublicKey string    `json:"serverPublicKey"`
	State           State     `json:"state"`
	Error           string    `json:"error,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type Store struct {
	mu       sync.RWMutex
	sessions map[string]Session
}

func NewStore() *Store { return &Store{sessions: make(map[string]Session)} }

func newSessionID() (string, error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func (s *Store) Create(session Session) (Session, error) {
	if session.CommandID == "" || session.NodeID == "" || session.ClientPublicKey == "" {
		return Session{}, errors.New("command, node and client key required")
	}
	id, err := newSessionID()
	if err != nil {
		return Session{}, err
	}
	now := time.Now().UTC()
	session.ID = id
	session.State = StatePending
	session.CreatedAt = now
	session.UpdatedAt = now
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[id] = session
	return session, nil
}

func (s *Store) Get(id string) (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.sessions[id]
	return v, ok
}

func (s *Store) SetResult(id string, ready bool, message string) (Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.sessions[id]
	if !ok {
		return Session{}, errors.New("session not found")
	}
	if ready {
		v.State = StateReady
		v.Error = ""
	} else {
		v.State = StateFailed
		v.Error = message
	}
	v.UpdatedAt = time.Now().UTC()
	s.sessions[id] = v
	return v, nil
}
