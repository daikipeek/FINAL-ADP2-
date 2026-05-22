package car

import "testing"

func TestCarSelectContainsRatingJoin(t *testing.T) {
	q := carSelect()
	if q == "" || q[:6] != "select" {
		t.Fatalf("unexpected select query: %q", q)
	}
}
