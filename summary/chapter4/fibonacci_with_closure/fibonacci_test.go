package practice

import (
	"reflect"
	"testing"
)

func TestFibonacci(t *testing.T) {
	fiboGenerator := Fibonacci()

	var got []int
	for i := 0; i < 10; i++ {
		got = append(got, fiboGenerator())
	}
	want := []int {1, 1, 2, 3, 5, 8, 13, 21, 34, 55}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}