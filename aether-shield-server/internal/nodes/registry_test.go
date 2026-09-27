package nodes

import (
	"testing"
	"time"
)

func TestRegistryDoesNotInventNodes(t *testing.T) {
	r := NewRegistry()
	if got := r.List(); len(got) != 0 {
		t.Fatalf("got %d nodes", len(got))
	}
}

func TestSelectRequiresOnline(t *testing.T) {
	r := NewRegistry()
	_ = r.Upsert(Node{ID: "n1", Region: "eu", Endpoint: "example.invalid:51820", PublicKey: "public-only", Status: StatusUnknown})
	if _, ok := r.SelectOnline("eu"); ok {
		t.Fatal("unknown node selected as online")
	}
}

func TestHeartbeatPromotesRegisteredNode(t *testing.T) {
	r := NewRegistry()
	_ = r.Upsert(Node{ID: "n1", Region: "eu", Endpoint: "example.invalid:51820", PublicKey: "public-only", Status: StatusUnknown})
	now := time.Now().UTC()
	if err := r.Heartbeat("n1", now); err != nil {
		t.Fatal(err)
	}
	n, _ := r.Get("n1")
	if n.Status != StatusOnline || n.LastSeen == nil {
		t.Fatal("heartbeat did not mark node online")
	}
}

func TestStaleNodeBecomesOffline(t *testing.T) {
	r := NewRegistry()
	old := time.Now().UTC().Add(-2 * time.Minute)
	_ = r.Upsert(Node{ID: "n1", Endpoint: "example.invalid:51820", PublicKey: "public-only", Status: StatusOnline, LastSeen: &old})
	r.MarkStaleOffline(time.Now().UTC(), 90*time.Second)
	n, _ := r.Get("n1")
	if n.Status != StatusOffline {
		t.Fatalf("status=%s", n.Status)
	}
}

func TestUpsertRejectsMissingPublicData(t *testing.T) {
	r := NewRegistry()
	if err := r.Upsert(Node{ID: "n1"}); err == nil {
		t.Fatal("expected validation error")
	}
}
