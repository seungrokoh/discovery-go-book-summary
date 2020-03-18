package practice

import "testing"

func TestFibonacci(t *testing.T) {
	want := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55}
	for i, v, fib := 0, 0, Fibonacci(); i < len(want); i, v = i+1, fib() {
		if v != want[i] {
			t.Errorf("got %v, want %v", v, want[i])
		}
	}
}
