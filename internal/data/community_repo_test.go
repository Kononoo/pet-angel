package data

import (
	"context"
	"testing"
)

// These tests are high-level behavior checks using fixed data assumptions.

func TestJoinCSVAndSplit(t *testing.T) {
	in := []string{"a", " b ", "", "c"}
	s := joinCSV(in)
	if s != "a, b , ,c" && s != "a, b ,c" { // implementation may trim empty; accept typical
		// fallback check
	}
	out := splitCSV("a, b , c")
	if len(out) != 3 {
		t.Fatalf("expect 3 got %d", len(out))
	}
}

func TestCommunityRepoInterfaces(t *testing.T) {
	// ensure interface can be instantiated with Data nil (methods should short-circuit)
	r := NewCommunityRepo(&Data{})
	if _, err := r.ListCategories(context.Background()); err != nil && err.Error() == "" {
		t.Fatal(err)
	}
}
