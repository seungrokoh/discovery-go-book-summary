package practice

func Fibonacci() func() int {
	x1, x2 := 1, 1
	return func() int {
		result := x1
		x1, x2 = x2, x1 + x2
		return result
	}
}
