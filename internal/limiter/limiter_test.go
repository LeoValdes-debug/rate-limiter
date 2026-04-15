package limiter

import (
	"testing"
	"time"
)

func TestAllow(t *testing.T) {
	l := New(3, 1)

	for i := 0; i < 3; i++ {
		if !l.Allow("1.2.3.4") {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}

	if l.Allow("1.2.3.4") {
		t.Fatal("4th request should be blocked")
	}
}

func TestRefill(t *testing.T) {
	l := New(2, 10)

	l.Allow("1.2.3.4")
	l.Allow("1.2.3.4")

	if l.Allow("1.2.3.4") {
		t.Fatal("should be blocked before refill")
	}

	time.Sleep(200 * time.Millisecond)

	if !l.Allow("1.2.3.4") {
		t.Fatal("should be allowed after refill")
	}
}

func TestDifferentIPs(t *testing.T) {
	l := New(1, 1)

	if !l.Allow("1.1.1.1") {
		t.Fatal("first ip should pass")
	}
	if !l.Allow("2.2.2.2") {
		t.Fatal("second ip should pass independently")
	}
}
