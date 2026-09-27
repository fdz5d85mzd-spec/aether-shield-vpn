package sessions

import (
	"errors"
	"fmt"
	"sync"
)

type Lease struct {
	ClientPublicKey string `json:"clientPublicKey"`
	Address         string `json:"address"`
}

type Allocator struct {
	mu     sync.Mutex
	base   string
	start  int
	end    int
	leases map[string]Lease
	used   map[string]bool
}

func NewAllocator() *Allocator {
	return &Allocator{
		base:   "10.88.0",
		start:  10,
		end:    250,
		leases: make(map[string]Lease),
		used:   make(map[string]bool),
	}
}

func (a *Allocator) Allocate(publicKey string) (Lease, error) {
	if publicKey == "" {
		return Lease{}, errors.New("client public key required")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if lease, ok := a.leases[publicKey]; ok {
		return lease, nil
	}
	for host := a.start; host <= a.end; host++ {
		address := fmt.Sprintf("%s.%d/32", a.base, host)
		if !a.used[address] {
			lease := Lease{ClientPublicKey: publicKey, Address: address}
			a.leases[publicKey] = lease
			a.used[address] = true
			return lease, nil
		}
	}
	return Lease{}, errors.New("address pool exhausted")
}

func (a *Allocator) Release(publicKey string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if lease, ok := a.leases[publicKey]; ok {
		delete(a.used, lease.Address)
		delete(a.leases, publicKey)
	}
}
