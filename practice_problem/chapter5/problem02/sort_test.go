package problem02

import (
	"reflect"
	"sort"
	"testing"
)

func TestCaseInsensitiveSort(t *testing.T) {
	cases := []struct {
		in CaseInsensitive
		out CaseInsensitive
	} {
		{
			CaseInsensitive([]string{"iPhone", "iPad", "MacBook", "AppStore"}),
			CaseInsensitive([]string{"AppStore", "iPad",  "iPhone", "MacBook"}),
		},
		{
			CaseInsensitive([]string{"orange", "apple", "banana", "tomato"}),
			CaseInsensitive([]string{"apple", "banana",  "orange", "tomato"}),
		},
	}

	for i, c := range cases {
		sort.Sort(c.in)
		if !reflect.DeepEqual(c.in, c.out) {
			t.Errorf("case index : %d\tgot %v, want %v", i, c.in, c.out)
		}
	}
}