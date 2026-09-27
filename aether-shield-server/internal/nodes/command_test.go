package nodes

import "testing"

func TestCommandQueueIsNodeScoped(t *testing.T) {
	q := NewCommandQueue()
	cmd, err := q.Enqueue("node-a", Command{Type: CommandAddPeer, ClientPublicKey: "client", ClientAddress: "10.88.0.10/32"})
	if err != nil {
		t.Fatal(err)
	}
	if got := q.Poll("node-b"); len(got) != 0 {
		t.Fatal("command leaked to another node")
	}
	if got := q.Poll("node-a"); len(got) != 1 {
		t.Fatalf("commands=%d", len(got))
	}
	if _, err := q.Ack("node-b", cmd.ID, true, ""); err == nil {
		t.Fatal("another node acknowledged command")
	}
	applied, err := q.Ack("node-a", cmd.ID, true, "")
	if err != nil {
		t.Fatal(err)
	}
	if applied.Status != CommandApplied {
		t.Fatalf("status=%s", applied.Status)
	}
}

func TestFailedAckRetainsReason(t *testing.T) {
	q := NewCommandQueue()
	cmd, _ := q.Enqueue("node-a", Command{Type: CommandAddPeer, ClientPublicKey: "client"})
	failed, err := q.Ack("node-a", cmd.ID, false, "wg set failed")
	if err != nil {
		t.Fatal(err)
	}
	if failed.Status != CommandFailed || failed.Error == "" {
		t.Fatalf("command=%#v", failed)
	}
}
