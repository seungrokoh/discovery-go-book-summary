package problem03

import "testing"

func TestIsContain(t *testing.T) {
	strs := []string {"a", "b", "c", "d"}

	got := IsContain(strs, "d")
	want := 3

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
