package problem02

import (
	"reflect"
	"testing"
)

var nums = []int{3, 1, 4, 2 ,5, 9, 8 ,7, 6}
var sortedNums = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

func TestSelectionSort(t *testing.T) {
	got := SelectionSort(nums)
	want := sortedNums

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
