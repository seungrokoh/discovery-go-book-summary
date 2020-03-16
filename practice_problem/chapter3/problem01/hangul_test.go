package problem01

import (
	"reflect"
	"testing"
)

func TestPrintFruitDelicious(t *testing.T) {

	t.Run("only hangul", func(t *testing.T) {
		fruits := []string{"사과", "바나나", "토마토", "귤"}

		got := PrintFruitDelicious(fruits)
		want := []string{"사과는 맛있다", "바나나는 맛있다", "토마토는 맛있다", "귤은 맛있다"}

		assertError(t, got, want)
	})

	t.Run("contain not hangul", func(t *testing.T) {
		fruits := []string{"사과", "바나나", "tomato", "귤"}

		got := PrintFruitDelicious(fruits)
		want := []string{"사과는 맛있다", "바나나는 맛있다", "귤은 맛있다"}

		assertError(t, got, want)
	})
}

func assertError(t *testing.T, got, want []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}