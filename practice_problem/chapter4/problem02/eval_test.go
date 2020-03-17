package problem02

import (
	"reflect"
	"testing"
)

func TestEval(t *testing.T) {
	eval := NewEvaluator(map[string]BinOp {
		"**": func(a, b int) int {
			if a == 1 {
				return 1
			}
			if b < 0 {
				return 0
			}
			r := 1
			for i := 0; i < b; i++ {
				r *= a
			}
			return r
		},
		"*": 	func(a, b int) int { return a * b},
		"/": 	func(a, b int) int { return a / b},
		"mod": 	func(a, b int) int { return a % b},
		"+": 	func(a, b int) int { return a + b},
		"-": 	func(a, b int) int { return a - b},
	}, PrecMap {
		"**":	NewStrSet(),
		"*":	NewStrSet("**", "*", "/", "mod"),
		"/":	NewStrSet("**", "*", "/", "mod"),
		"mod":	NewStrSet("**", "*", "/", "mod"),
		"+":	NewStrSet("**", "*", "/", "mod", "+", "-"),
		"-":	NewStrSet("**", "*", "/", "mod", "+", "-"),
	})

	exprs := []string{"5",
		"1 +      2",
		"1 - 2 - 4",
		"( 3 - 2 ** 3 ) * ( -2 )",
		"( 3 - 2 ** 3 ) * ( -2 ))",
		"3 * ( ( 3 + 1 ) * 3 ) / 2",
		"2 ** 3 mod 3",
		"2 ** 2 ** 3",
		"2 * 2 *",
		"2 ! 3"}
	var got []int
	for _, expr := range exprs {
		val, err := eval(expr)
		if err != nil {
			t.Errorf("%s", err)
		} else {
			got = append(got, val)
		}
	}

	want := []int{5, 3, -5, 10, 18, 2, 256}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
