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

	t.Run("들어있는 원소보다 더 많은 Dequeue 실행", func(t *testing.T) {
		items := []int {1, 2, 3, 4, 5}
		queue := Queue{}

		for _, item := range items {
			queue.Enqueue(item)
		}

		queue.Dequeue()
		queue.Dequeue()
		queue.Dequeue()
		queue.Dequeue()
		queue.Dequeue()
		queue.Dequeue()

		queue.Enqueue(3)

		got := queue.items
		want := []int{3}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
