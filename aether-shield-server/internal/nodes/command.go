package nodes

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

type CommandType string
type CommandStatus string

const (
	CommandAddPeer    CommandType = "ADD_PEER"
	CommandRemovePeer CommandType = "REMOVE_PEER"

	CommandPending CommandStatus = "PENDING"
	CommandApplied CommandStatus = "APPLIED"
	CommandFailed  CommandStatus = "FAILED"
)

type Command struct {
	ID              string        `json:"id"`
	Type            CommandType   `json:"type"`
	ClientPublicKey string        `json:"clientPublicKey"`
	ClientAddress   string        `json:"clientAddress,omitempty"`
	CreatedAt       time.Time     `json:"createdAt"`
	Status          CommandStatus `json:"status"`
	Error           string        `json:"error,omitempty"`
}

type CommandQueue struct {
	mu      sync.Mutex
	pending map[string][]string
	records map[string]Command
	owners  map[string]string
}

func NewCommandQueue() *CommandQueue {
	return &CommandQueue{
		pending: make(map[string][]string),
		records: make(map[string]Command),
		owners:  make(map[string]string),
	}
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
	cmd.Status = CommandPending
	q.mu.Lock()
	defer q.mu.Unlock()
	q.records[id] = cmd
	q.owners[id] = nodeID
	q.pending[nodeID] = append(q.pending[nodeID], id)
	return cmd, nil
}

func (q *CommandQueue) Poll(nodeID string) []Command {
	q.mu.Lock()
	defer q.mu.Unlock()
	ids := q.pending[nodeID]
	delete(q.pending, nodeID)
	out := make([]Command, 0, len(ids))
	for _, id := range ids {
		if cmd, ok := q.records[id]; ok {
			out = append(out, cmd)
		}
	}
	return out
}

func (q *CommandQueue) Ack(nodeID, commandID string, applied bool, message string) (Command, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.owners[commandID] != nodeID {
		return Command{}, errors.New("command not owned by node")
	}
	cmd, ok := q.records[commandID]
	if !ok {
		return Command{}, errors.New("command not found")
	}
	if applied {
		cmd.Status = CommandApplied
		cmd.Error = ""
	} else {
		cmd.Status = CommandFailed
		cmd.Error = message
	}
	q.records[commandID] = cmd
	return cmd, nil
}

func (q *CommandQueue) Get(commandID string) (Command, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	cmd, ok := q.records[commandID]
	return cmd, ok
}
