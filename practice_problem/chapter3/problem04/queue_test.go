package problem04

import (
	"reflect"
	"testing"
)

func TestQueue(t *testing.T) {

	t.Run("Enqueue Items", func(t *testing.T) {
		items := []int {1, 2, 3, 4, 5}
		queue := Queue{}

		for _, item := range items {
			queue.Enqueue(item)
		}

		got := queue.items
		want := []int{1, 2, 3, 4, 5}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Dequeue Item", func(t *testing.T) {
		items := []int {1, 2, 3, 4, 5}
		queue := Queue{}

		for _, item := range items {
			queue.Enqueue(item)
		}

		queue.Dequeue()

		got := queue.items
		want := []int{2, 3, 4, 5}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
