package nodeagent

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"
)

type RuntimeProbe interface {
	Verify(context.Context, Config) error
}

type WireGuardRuntime struct {
	Interface string
	Timeout   time.Duration
}

func (w WireGuardRuntime) Verify(ctx context.Context, config Config) error {
	if w.Interface == "" {
		return errors.New("wireguard interface is required")
	}
	timeout := w.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	probeCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := exec.CommandContext(probeCtx, "ip", "link", "show", "dev", w.Interface).Run(); err != nil {
		return errors.New("wireguard interface unavailable")
	}
	out, err := exec.CommandContext(probeCtx, "wg", "show", w.Interface, "public-key").Output()
	if err != nil {
		return errors.New("wireguard runtime unavailable")
	}
	actual := strings.TrimSpace(string(out))
	if actual == "" || actual != strings.TrimSpace(config.PublicKey) {
		return errors.New("wireguard public key mismatch")
	}
	return nil
}
