package nodeagent

import (
	"context"
	"errors"
	"testing"
)

type fakeProbe struct{ err error }

func (f fakeProbe) Verify(context.Context, Config) error { return f.err }

func TestProbeContractCanFailClosed(t *testing.T) {
	var probe RuntimeProbe = fakeProbe{err: errors.New("offline")}
	if err := probe.Verify(context.Background(), Config{}); err == nil {
		t.Fatal("expected runtime verification failure")
	}
}
