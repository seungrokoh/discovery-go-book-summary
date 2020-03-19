package problem04

import (
	"reflect"
	"sort"
	"testing"
)

func TestTasksSort(t *testing.T) {
	tasks := TaskList([]Task {
		{
			Title:"title1",
			Priority:4,
		},
		{
			Title:"title2",
			Priority:2,
		},
		{
			Title:"title3",
			Priority:3,
		},
		{
			Title:"title4",
			Priority:1,
		},
	})

	sort.Sort(tasks)
	want := TaskList([]Task {
		{
			Title:"title4",
			Priority:1,
		},
		{
			Title:"title2",
			Priority:2,
		},
		{
			Title:"title3",
			Priority:3,
		},
		{
			Title:"title1",
			Priority:4,
		},
	})
	if !reflect.DeepEqual(tasks, want) {
		t.Errorf("got %v, want %v", tasks, want)
	}
}