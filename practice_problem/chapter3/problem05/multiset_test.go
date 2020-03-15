package problem05

import (
	"strconv"
	"testing"
)

func TestMultiset(t *testing.T) {
	t.Run("Test Multiset", func(t *testing.T) {
		multiset := NewMultiSet()
		var got string

		got += String(multiset)
		got += strconv.Itoa(Count(multiset, "3"))

		Insert(multiset, "3")
		Insert(multiset, "3")
		Insert(multiset, "3")
		Insert(multiset, "3")

		got += String(multiset)
		got += strconv.Itoa(Count(multiset, "3"))

		Insert(multiset, "1")
		Insert(multiset, "2")
		Insert(multiset, "5")
		Insert(multiset, "7")
		Erase(multiset, "3")
		Erase(multiset, "5")

		got += strconv.Itoa(Count(multiset, "3"))
		got += strconv.Itoa(Count(multiset, "1"))
		got += strconv.Itoa(Count(multiset, "2"))
		got += strconv.Itoa(Count(multiset, "5"))

		want := "{ }0{ 3 3 3 3 }43110"

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}

