package concurrency

import (
	"log"
	"sync"
	"testing"
)

func TestPlusOneService(t *testing.T) {
	reqs := make(chan Request)
	defer close(reqs)

	for i := 0; i < 3; i++ {
		go PlusOneService(reqs, i)
	}

	var wg sync.WaitGroup
	for i := 3; i < 53; i += 10 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			resps := make(chan Response)
			reqs <- Request{i, resps}
			log.Println(i, "=>", <-resps)
		}(i)
	}
	wg.Wait()
}
