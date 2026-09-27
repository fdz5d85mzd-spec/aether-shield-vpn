package nodeagent

import (
	"context"
	"testing"
)

func TestPeerManagerRejectsMissingInput(t *testing.T) {
	p := PeerManager{}
	if err := p.Add(context.Background(), "", ""); err == nil {
		t.Fatal("expected validation error")
	}
	if err := p.Remove(context.Background(), ""); err == nil {
		t.Fatal("expected validation error")
	}
}
