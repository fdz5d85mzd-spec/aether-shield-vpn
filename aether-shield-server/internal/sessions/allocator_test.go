package sessions

import "testing"

func TestAllocatorStableLease(t *testing.T) {
	a := NewAllocator()
	first, err := a.Allocate("client-a")
	if err != nil {
		t.Fatal(err)
	}
	second, err := a.Allocate("client-a")
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatalf("lease changed: %#v %#v", first, second)
	}
}

func TestAllocatorDoesNotShareAddress(t *testing.T) {
	a := NewAllocator()
	x, _ := a.Allocate("client-a")
	y, _ := a.Allocate("client-b")
	if x.Address == y.Address {
		t.Fatal("duplicate address")
	}
}

func TestReleaseMakesAddressReusable(t *testing.T) {
	a := NewAllocator()
	x, _ := a.Allocate("client-a")
	a.Release("client-a")
	y, _ := a.Allocate("client-b")
	if x.Address != y.Address {
		t.Fatalf("expected reusable address, got %s then %s", x.Address, y.Address)
	}
}
