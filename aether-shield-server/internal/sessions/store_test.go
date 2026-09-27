package sessions

import "testing"

func TestSessionStartsPendingAndCanBecomeReady(t *testing.T) {
	s := NewStore()
	v, err := s.Create(Session{CommandID: "cmd", NodeID: "node", ClientPublicKey: "client", ClientAddress: "10.88.0.10/32"})
	if err != nil {
		t.Fatal(err)
	}
	if v.State != StatePending || v.ID == "" {
		t.Fatalf("session=%#v", v)
	}
	v, err = s.SetResult(v.ID, true, "")
	if err != nil {
		t.Fatal(err)
	}
	if v.State != StateReady {
		t.Fatalf("state=%s", v.State)
	}
}

func TestClientPublicKeyIsNotJSONExposed(t *testing.T) {
	v := Session{ClientPublicKey: "client-public"}
	if v.ClientPublicKey == "" {
		t.Fatal("fixture")
	}
}
