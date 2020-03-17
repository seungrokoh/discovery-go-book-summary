package problem03

import (
	"regexp"
	"strconv"
	"strings"
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

	in := strings.Join([]string{"다들 그 동안 고생이 많았다.",
		"첫째는 분당에 있는 { 2 ** 4 * 3 }평 아파트를 갖거라.",
		"둘째는 임야 { 10 ** 5 mod 7777 }평을 가져라."}, "\n")

	rx := regexp.MustCompile(`{[^}]+}`)
	got := rx.ReplaceAllStringFunc(in, func(expr string) string {
		expr = strings.Trim(expr, "{ }")
		val, _ := eval(expr)

		return strconv.Itoa(val)
	})

	want := strings.Join([]string{"다들 그 동안 고생이 많았다.",
		"첫째는 분당에 있는 48평 아파트를 갖거라.",
		"둘째는 임야 6676평을 가져라."}, "\n")

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}

}
