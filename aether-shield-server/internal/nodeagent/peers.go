package nodeagent

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"
)

type PeerManager struct {
	Interface string
	Timeout   time.Duration
}

func (p PeerManager) Add(ctx context.Context, publicKey, address string) error {
	if p.Interface == "" || strings.TrimSpace(publicKey) == "" || strings.TrimSpace(address) == "" {
		return errors.New("interface, public key and address are required")
	}
	timeout := p.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	cmdCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := exec.CommandContext(cmdCtx, "wg", "set", p.Interface, "peer", publicKey, "allowed-ips", address).Run(); err != nil {
		return errors.New("failed to apply WireGuard peer")
	}
	return nil
}

func (p PeerManager) Remove(ctx context.Context, publicKey string) error {
	if p.Interface == "" || strings.TrimSpace(publicKey) == "" {
		return errors.New("interface and public key are required")
	}
	timeout := p.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	cmdCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := exec.CommandContext(cmdCtx, "wg", "set", p.Interface, "peer", publicKey, "remove").Run(); err != nil {
		return errors.New("failed to remove WireGuard peer")
	}
	return nil
}
