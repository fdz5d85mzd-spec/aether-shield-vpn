package nodes

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

type CommandType string

const (
	CommandAddPeer    CommandType = "ADD_PEER"
	CommandRemovePeer CommandType = "REMOVE_PEER"
)

type Command struct {
	ID              string      `json:"id"`
	Type            CommandType `json:"type"`
	ClientPublicKey string      `json:"clientPublicKey"`
	ClientAddress   string      `json:"clientAddress,omitempty"`
	CreatedAt       time.Time   `json:"createdAt"`
}

type CommandQueue struct {
	mu      sync.Mutex
	pending map[string][]Command
}

func NewCommandQueue() *CommandQueue {
	return &CommandQueue{pending: make(map[string][]Command)}
}

func commandID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func (q *CommandQueue) Enqueue(nodeID string, cmd Command) (Command, error) {
	if nodeID == "" || cmd.ClientPublicKey == "" {
		return Command{}, errors.New("node ID and client public key required")
	}
	id, err := commandID()
	if err != nil {
		return Command{}, err
	}
	cmd.ID = id
	cmd.CreatedAt = time.Now().UTC()
	q.mu.Lock()
	defer q.mu.Unlock()
	q.pending[nodeID] = append(q.pending[nodeID], cmd)
	return cmd, nil
}

func (q *CommandQueue) Drain(nodeID string) []Command {
	q.mu.Lock()
	defer q.mu.Unlock()
	items := append([]Command(nil), q.pending[nodeID]...)
	delete(q.pending, nodeID)
	return items
}
