func Fibonacci(max int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		a, b := 0, 1
		for a <= max {
			out <- a
			a, b = b, a + b
		}
	}()
	return out
}

func FibonacciGenerator(max int) func() int {
    next, a, b := 0, 0, 1
    return func() int {
        next, a, b = a, b, a + b
        if next > max {
            return -1
        }
        return next
    }
}

func main() {
    // 채널을 이용한 fibonacci
    for fib := range Fibonacci(15) {
        fmt.Print(fib, ", ")
    }

    // 클로져를 이용한 fibonacci
    fib := FibonacciGenerator(15)
    for n := fib(); n >= 0; n = fib() {
        fmt.Print(n, ", ")
    }
}
