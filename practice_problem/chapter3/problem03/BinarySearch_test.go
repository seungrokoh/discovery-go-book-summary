package problem03

import "testing"

func TestIsContain(t *testing.T) {
	strs := []string {"a", "b", "c", "d"}

	got := IsContain(strs, "f")
	want := false

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
