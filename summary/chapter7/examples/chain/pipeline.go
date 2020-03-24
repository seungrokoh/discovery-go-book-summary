package concurrency

func PlusOne(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			out <- num + 1
		}
	}()
	return out
}

func Square(in <-chan int) <- chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			out <- num * num
		}
	}()
	return out
}

type IntPipe func(<-chan int) <-chan int

func Chain(ps ...IntPipe) IntPipe {
	return func(in <-chan int) <- chan int {
		c := in
		for _, p := range ps {
			c = p(c)
		}
		return c
	}
}

func FanIn(ins ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(ins))

	for _, in := range ins {
		go func(in <-chan int) {
			defer wg.Done()
			for num := range in {
				out <- num
			}
		}(in)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func FanIn3(in1, in2, in3 <-chan int) <-chan int {
	out := make(chan int)
	openCnt := 3
	closeChan := func(c *<-chan int) bool {
		*c = nil
		openCnt--
		return openCnt == 0
	}

	go func() {
		defer close(out)
		for {
			select {
			case n, ok := <-in1:
				if ok {
					out <- n
				} else if closeChan(&in1) {
					return
				}
			case n, ok := <-in2:
				if ok {
					out <- n
				} else if closeChan(&in2) {
					return
				}
			case n, ok := <-in3:
				if ok {
					out <- n
				} else if closeChan(&in3) {
					return
				}
			}
		}
	}()
	return out
}
