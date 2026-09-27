package nodes

import "testing"

func TestCommandQueueIsNodeScoped(t *testing.T) {
	q := NewCommandQueue()
	_, err := q.Enqueue("node-a", Command{Type: CommandAddPeer, ClientPublicKey: "client", ClientAddress: "10.88.0.10/32"})
	if err != nil {
		t.Fatal(err)
	}
	if got := q.Drain("node-b"); len(got) != 0 {
		t.Fatal("command leaked to another node")
	}
	if got := q.Drain("node-a"); len(got) != 1 {
		t.Fatalf("commands=%d", len(got))
	}
}
