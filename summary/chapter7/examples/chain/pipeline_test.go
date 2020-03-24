package concurrency

import (
	"reflect"
	"testing"
)

func TestChain(t *testing.T) {
	c := make(chan int)
	go func() {
		defer close(c)
		c <- 5
		c <- 3
		c <- 8
	}()

	var got[]int
	want := []int{36, 16, 81}

	for num := range Chain(PlusOne, Square)(c) {
		got = append(got, num)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
