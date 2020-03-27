# 7장 동시성

## 7.1 고루틴

* 고루틴은 가벼운 스레드와 같은 것으로 **현재 수행 흐름과 별개의 흐름을 만들어준다.**
* 함수 앞에 `go` 키워드를 붙여 고루틴을 실행할 수 있다.
* go 키워드가 붙은 함수는 메모리를 공유하는 논리적으로 별개의 흐름이 된다.



### 7.1.1 병렬성과 병행성

* 병렬성은 정말 동시에 각각의 흐름이 수행되는 경우를 뜻한다.
* 병행성은 실제로는 동시에 흐르는 것이 아닌 동시에 실행되는 것처럼 보이는 것이다.



go 키워드가 붙은 함수 호출은 별개의 흐름으로 동작한다. 만약 여러개의 goroutine 함수가 있다면 concurrency로 동작하게된다.



```go
func main() {
    go func() {
        fmt.Println("In goroutine")
    }()
    fmt.Println("In main routine")
}
```

* 메인 함수가 끝나버리면 goroutine이 모두 수행되지 않을 수도 있다.
* 위 코드는 goroutine이 항상 실행되지 않는 코드이다.
* main 함수가 끝날 때 goroutine이 준비가 되지 않아 출력을 하지 못할 수 있기 때문이다.



### 7.1.2 고루틴 기다리기

* 고루틴을 제어하기 위하여 sync library를 제공한다.
* sync.WaitGroup을 이용해 선행작업이 있을 때 goroutine들을 제어할 수 있다.

p.240 ~ p.242 의 사진 파일을 다운로드하고 zipping하는 코드를 살펴보면 이해할 수 있다. 해당 코드는 제대로 동작하지 않는데 이유는 다음과 같다.

* download 함수는 goroutine으로 동작하는데 이 함수가 언제 끝날지 알 수 없다.
* zipping을 하기 위해선 다운로드 된 사진이 필요하지만 goroutine은 별개로 동작하기 때문에 zipping하는 시점에 다운로드된 사진이 없을 수 있다.
* goroutine으로 동작 할 경우 여러 흐름으로 나눠서 하기 때문에 빠를 수 있지만 선행작업이 있을 경우 따로 제어를 해줘야한다.

위 예제는 두 작업이 연관되어 있다. 선행작업 (파일 다운로드) 이후 zipping을 해야하는 경우다. 이때 sync.WaitGroup을 이용해 만들어진 고루틴들이 모두 작업을 마칠 때까지 기다릴 수 있다. 

```go
var wg sync.WaitGroup
// 생성될 고루틴의 갯수만큼 add
wg.Add(len(urls))
for _, url := range urls {
    go func(ulr string) {
        // 고루틴이 종료될 때 Done을 호출
        defer wg.Done()
        if _, err := download(url); err != nli {
            log.Fatal(err)
        }
    }(url)
}
// 고루틴이 모두 끝날 때까지 기다림
wg.Wait()
```

* wg 에는 기본값이 0으로 맞춰져 있는 카운터가 들어 있다.
* Wait() 함수는 카운터가 0일 될 때까지 기다린다.
* Add() 함수는 호출될 때 넘긴 인자 만큼 카운터를 증가시킨다.
* Done() 함수는 Add(-1)과 동일하며 동작을 마쳤다는 의미로 표현한다.
* 즉 동작할 고루틴들의 갯수를 지정해놓고 각 고루틴들이 종료될 때 Done을 호출함으로써 모든 고루틴들이 끝났음을 확인할 수 있다.

```go
var wg sync.WaitGroup
for _, url := range urls {
    wg.Add(1)
    go func(url string) {
        defer wg.Done()
        if _, err := download(url); err != nil {
            log.Fatal(err)
        }
    }(url)
}
wg.Wait()
```

* 처음부터 생성 될 고루틴의 개수를 알지 못할 때 고루틴을 생성할 때 Add() 함수를 호출해 해결할 수 있다.
* 주의 할 점은 Add() 함수를 고루틴 내부에 포함시키면 안된다. 그 이유는 제대로 동작할 수도 있지만 고루틴 안에 선언해 놓을 경우 Add()하기 전에 wg.Wait()을 통과 할 가능성이 있기 때문이다. (race condition) 



### 공유 메모리와 병렬 최소값 찾기

* 고루틴들은 메모리도 서로 공유한다. 
* 변수의 포인터를 받아서 해당 변수에 원하는 값을 넣어줄 수 있다.

아래 예제는 배열에서 최소값을 찾는 문제이다. 먼저 배열에서 최소값을 찾는 함수를 만들고 배열 전체를 n개로 나눠 n개의 goroutine으로 최소값을 찾는 문제이다. 



```go
// 배열에서 최소값을 찾는 함수
func Min(a []int) int {
	if len(a) == 0 {
		return 0
	}
	min := a[0]
	for _, e := range a[1:] {
		if min > e {
			min = e
		}
	}
	return min
}

// 배열을 n조각으로 나눠 각 고루틴에서 최소값을 찾는 함수
// 최소값을 찾을 배열과 몇 개의 goroutine으로 나눌지 인자로 넘겨준다.
func ParallelMin(a []int, n int) int {
	if len(a) < n {
		return Min(a)
	}
    // 각 고루틴에서 찾은 최소값들을 저장할 mins 배열
	mins := make([]int, n)
	bucketSize := (len(a) + n - 1) / n
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
        // 고루틴 생성
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			begin, end := i*bucketSize, (i+1)*bucketSize
			if end > len(a) {
				end = len(a)
			}
            // 고루틴에서 최소값을 찾고 배열에 저장
			mins[i] = Min(a[begin:end])
		}(i)
	}
    // 모든 고루틴에서 최소값을 찾을 때까지 대기
	wg.Wait()
    // 각 고루틴에서 찾은 최소값들 중 최소값을 최종 반환
	return Min(mins)
}
```

* 실제 1억개의 임의의 데이터를 만들고 최소값을 찾는데 걸리는 시간을 측정해 볼 수 있다.

```go
func main() {
	var nums []int

	// 1억개의 임의 값들을 생성
	for i := 0; i < 100000000; i++ {
		nums = append(nums, rand.Int())
	}
	start := time.Now()
	fmt.Println("linear Min 최소 값 : ", Min(nums))
	fmt.Println("linear Min 수행 시간 : ", time.Since(start))
	start = time.Now()
	fmt.Println("parallel Min 최소 값 : ", ParallelMin(nums, 100))
	fmt.Println("parallel Min 수행 시간 : ", time.Since(start))
}
```

* 각 수행 시간은 실행 환경마다 다르며 ParallelMin이 더 빠른 수행결과를 보인다는 것을 볼 수 있다.
* 위 예제는 100개의 고루틴들이 1억개의 데이터를 100등분 하여 각각 계산하는 예제이다.



## 7.2 채널

* 서로 다른 고루틴끼리 통신할 수 있는 수단이다.
* 채널은 넣은 데이터를 뽑아낼 수 있는 마치 파이프와 같은 형태의 자료구조이다.
* 채널은 일급시민이다.
* 양방향 채널과 단반향 채널이 있으며 양방향은 단방향으로 변환할 수 있다.
* 맵처럼 생성해서 사용하며 채널끼리 서로 복사한 경우 동일한 채널을 가리킨다. (포인터와 비슷한 레퍼런스형)

```go
c1 := make(chan int)
var c2 chan int = c1
var c3 <-chan int = c1
var c4 chan<- int = c1
```

* c1과 c2는 동일한 채널을 보고 있으며 c3는 (receive only), c4는 (send only) 채널이 된다.
* receive only, send only는 화살표의 위치에 따라 정해진다. (<-chan or chan<-)

```go
// 채널에 자료 보내기
c <- 100

// 채널에서 받은 자료 버리기
<-c

// 채널에서 자료 받기
data := <-c
```



### 7.2.1 일대일 단방향 채널 소통

* 채널 이용의 가장 단순한 형태로 고루틴 하나에서는 보내고 다른 고루틴에서 받는 형태이다.

```go
func main() {
	c := make(chan int)
	go func() {
		c <- 1
		c <- 2
		c <- 3
	}()
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(<-c)
}

// Output:
// 1
// 2
// 3
```

* go func()로 호출된 고루틴 내에서 1,2,3 데이터를 보내고, 밖에서는 데이터를 받는다.
* 메인 함수도 고루틴으로 동작하기 때문에 두 고루틴 사이에서 서로 데이터를 주고 받는 모습이다.
* 보내는 고루틴과 받는 고루틴이 주고 받는 데이터의 개수가 다르면 고루틴은 멈추게 된다.

```go
func main() {
	c := make(chan int)
	go func() {
		c <- 1
		c <- 2
	}()
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(<-c)
}
```

* 위 코드는 보내는 개수가 받는 개수보다 적을 때 나타나는 오류를 알려준다
* 세 번째 fmt.Println(<-c) 부분에서 채널에서 데이터가 오기만을 기다리다 deadlock이 발생해 프로그램이 죽는다.

```go
func main() {
	c := make(chan int)
	go func() {
		c <- 1
		c <- 2
		c <- 3
	}()
	fmt.Println(<-c)
	fmt.Println(<-c)
}
// Output:
// 1
// 2
```

* 위 코드는 받는 개수가 보내는 개수보다 적을 때 나타나는 현상이다.
* 2개의 코드만 받았으므로 출력이 되고 main 함수가 종료되어 이상이 없어 보이지만 문제가 있다.
* go func()로 생성된 고루틴은 c <- 3 에서 멈추게 된다. 
* 하지만 main 함수는 이미 종료가 됐고 만들어진 고루틴은 허공에 뜨게 된다. 이게 반복되면 memory leak이 발생한다.
* 자세히 알아보기 위해 runtime.NumGoroutine() 을 이용해 현재 생성된 고루틴의 수를 세어보면 된다.

```go
func main() {
	c := make(chan int)
	go func() {
		c <- 1
		c <- 2
		c <- 3
	}()
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println("num of goroutine : ", runtime.NumGoroutine())
}
// Output:
// 1
// 2
// num of goroutine : 2
```

* 메인과 생성된 고루틴 2개가 남아있게 된다.
* 보내고 받는 데이터의 개수를 서로 알지 못해도 동작하도록 구현할 수 있다.

```go
func main() {
	c := make(chan int)
	go func() {
		c <- 1
		c <- 2
		c <- 3
		close(c)
	}()
	for num := range c {
		fmt.Println(num)
	}
}
// Output:
// 1
// 2
// 3
```

* close(c)를 이용해 채널을 닫을 수 있다.
* range 를 이용해 채널에서 값을 받아올 수 있다.
* range로 채널을 받아오는건 채널이 오픈되어 있을 때 받아온다. 즉, 채널이 닫히면 for문도 벗어나게 된다.



### :star::star::star::star::star:

#### 너무나 중요한 사용 패턴

별 5개인 이유는 이 코드를 이해 하면 다른 코드들도 이해하기 쉽기 때문이다. 그 말은 이 코드 패턴을 많은 곳에서 사용한다는 의미이다. 

* 채널 하나를 만들어서 넘겨주고 받는 것은 깔끔해 보이지 않는다.
* 따라서 함수가 채널을 반환하게 하는 패턴을 사용한다.
* 채널을 닫을 땐 defer close(c)를 이용하자.
* 보내는 쪽에서 단반향 채널을 반환하여 이 채널을 이용하는 고루틴이 받아가기만 할 수 있게 구현한다.



:star::star::star::star::star:

```go
func main() {
	c := func() <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			out <- 1
			out <- 2
			out <- 3
		}()
		return out
	}()
	for num := range c {
		fmt.Println(num)
	}
}
```

이 패턴은 반드시 외워두도록 하자. 패턴을 만드는 방법은 다음과 같다.

* 안에서 채널(이용하는 고루틴이 받아가는)을 생성 (out)
* out 채널을 사용하는 고루틴 생성
* 고루틴 안에서 채널에 데이터 넣기
* 생성한 채널을 반환

자세히 보면 클로져(closure) 형태로 구성되어 있다. 클로져를 모른다면 [4장 함수](https://github.com/seungrokoh/discovery-go-book-summary/tree/master/summary/chapter4)를 다시 읽어보고 이해 해야한다.



### 7.2.2 생성기 패턴

* 채널을 이용하여 생성기를 만들 수 있다.

```go
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
```

클로저를 이용해 만들 수 있지만 **채널을 이용하는 방법에는 몇 가지 장점이 있다.**

* 생성하는 쪽에서 상태 저장 방법을 복잡하게 고민할 필요가 없다.
* 받는 쪽에서 for range를 사용할 수 있다.
* 채널 버퍼를 이용하면 멀티 코어를 활용하거나 입출력 성능상의 장점을 이용할 수 있다.

클로져를 이용한 코드와 비교는 책 또는 [여기](https://github.com/seungrokoh/discovery-go-book-summary/tree/master/summary/chapter7/examples/fibonacci/fibonacci.go)서 확인할 수 있다.



### 7.2.3 버퍼 있는 채널

```go
c := make(chan int, 10)
```

* 받는 쪽이 준비가 되어 있지 않아도 보내는 쪽에서 미리 보낼 수 있다.
* 버퍼가 가득차기 전까지는 받는 쪽에서 준비가 되어 있지 않아도 보낼 수 있다.
* 보내는 쪽과 받는 쪽의 어느 정도 격차가 생겨도 계속 동작할 수 있기 때문에 성능 향상이 일어날 수 있다.
* 하지만 동시성은 강력하지만 복잡할 수 있으므로 알려진 패턴을 따르는게 더 좋다.
* **버퍼 없는 채널로 동작하는 코드를 만들고 필요에 따라 성능 향상을 위해 버퍼를 조절해주는게 좋다.**



### 7.2.4 닫힌 채널

* 채널이 닫히면 for range를 이용할 때 반복이 종료된다.
* 채널에서 값을 받아올 때 2개의 변수로 받아올 수 있는데 첫 번째는 값, 두 번째는 채널의 open 여부(bool)이다.
* 첫 번째 값은 오픈 여부에 상관 없이 값이 들어온다. 열려있으면 받은 값, 닫혀있으면 기본값이 들어온다.
* 두 번째 값은 열려 있으면 true, 닫혀 있으면 false가 들어온다.
* 채널이 닫혀 있으면 채널에 값이 보내질 리 없으므로 전혀 기다리지 않고 무작정 기본값을 받아올 수 있는 채널이 된다.

```go
func Example_closedChannel() {
    c := make(chan int)
    close(c)
    fmt.Println(c)
    fmt.Println(c)
    fmt.Println(c)
    // Output:
    // 0
    // 0
    // 0
}
```

* 닫혀있는 채널을 또 닫으려고 시도하면 패닉(panic)이 발생한다.



### :star::star::star::star::star:

## 7.3 동시성 패턴

이 파트에서 나오는 패턴들은 되도록 외워두는 것이 좋다. 위에서 설명한 **너무나 중요한 패턴**을 익히고 있다면 이해하기 수월하다.



### 7.3.1 파이프라인 패턴

* 파이프라인 패턴은 한 단계의 출력이 다음 단계의 입력으로 이어지는 구조이다.
* 분업구조로 컨베이어 벨트를 생각하면 쉽다. (상세 예시는 책에 나와있다.)
* 들어오는 데이터와 나가는 데이터에 집중하여 문제를 해결할 때 좋다.
* 파이프라인 패턴은 생성기 패턴의 일종이다.
* 받기 전용 채널을 입력으로 활용한다는 점이 생성기 패턴과의 차이이다.
* 반환된 받기 전용 채널을 다른 파이프라인의 입력으로 넘겨줄 수 있기 때문에 매우 자연스럽게 출력을 입력으로 연결하여 일직선으로 사슬처럼 연결된 파이프라인을 구성할 수 있다.

```go
// 채널로 들어온 값들에 1을 더해주는 함수
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
```

* Input 채널에서 받은 값들에 1을 더해 out 채널에 넣고 이를 반환하는 함수이다.
* 입력과 출력 모두 채널을 받고 반환하기 때문에 다음과 같이 사용할 수 있다.

```go
func main() {
    c := make(chan int)
    go func {
        defer close(c)
        c <- 5
        c <- 3
        c <- 8
    }()
    
    for num := range PlusOne(PlusOne(c)) {
        fmt.Println(num)
    }
}
// Output:
// 7
// 5
// 10
```

* 안쪽 PlusOne 함수를 거져 1 증가된 수가 저장된 채널을 다시 두 번째 PlusOne 함수에서 입력으로 사용하는 구조이다.
* 형태만 같다면 서로 다른 함수들도 이렇게 이어 붙일 수 있다.
* 1을 더하고, 제곱을 한 뒤, 2를 곱해주는 함수를 파이프 라인으로 구현할 수 있다.

```go
// 채널로 들어온 값에 1을 더해주는 함수
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

// 채널로 들어온 값을 제곱하는 함수
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

// 채널로 들어온 값에 2를 곱해주는 함수
func MultiplyDouble(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			out <- num * 2
		}
	}()
	return out
}

// 메인 함수
func main() {
	c := make(chan int)
	go func() {
		defer close(c)
		c <- 1
		c <- 2
		c <- 3
	}()

	for num := range MultiplyDouble(Square(PlusOne(c))) {
		fmt.Println(num)
	}
}
// Output:
// 4
// 18
// 32
```

* 같은 함수 뿐만 아니라 형태만 같다면 서로 다른 함수들도 이어 붙일 수 있다.

```go
type IntPipe func(<-chan int) <-chan int
```

* 생성기 패턴과 마찬가지로 무조건 **데이터를 보내는 쪽에서 채널을 닫아야 한다.**
* 사슬처럼 이어진 파이프라인을 하나로 보이게 만들어야 할 때 다음과 같은 패턴을 이용한다.

```go
func Chain(ps ...IntPipe) IntPipe {
    return func(in <-chan int) <-chan int {
        c := in
        for _, p := range ps {
            c = p(c)
        }
        return c
    }
}
```

* Chain(A, B)(c)는 B(A(c))와 같다.

* for문으로 파이프라인 함수들(ps)를 순차적으로 실행시킨다.
* 각 함수 p를 이용해 c를 입력으로 사용해 반환된 채널을 다시 c에 입력시킨다.
* 이를 반복적으로 반환되어 나온 c를 다음 파이프라인 함수 p(c)에 입력으로 넣는다.
* 바로 이용하지 않고 다른 곳에 넘겨야 하는 경우에 Chain 고계 함수를 이용하면 편리하다.



### 7.3.2 채널 공유로 팬아웃하기

* 한 곳에서의 출력을 여러곳에 나누어주어야 할 때 사용할 수 있다.
* 채널 하나를 여러 고루틴에게 공유하면 된다.
* 여러 고루틴이 한 채널에서 자료를 받아가려고 할 때, 하나의 고루틴에만 자료가 전달된다.
* 보내는 쪽에서 모든 자료를 보내면 채널을 닫는다.
* 채널이 닫히면 모든 해당 채널을 이용하는 모든 고루틴이 닫혔음을 인지한다. (broadcast)



:heavy_check_mark: FanOut 예제

```go
func main() {
    c := make(chan int)
    defer close(c)
    // 총 3개의 고루틴을 생성하고 채널에서 값을 받아 고루틴 번호와 값을 출력한다.
    for i := 0; i < 3; i++ {
        go func(i int) {
            for n := range c {
                // 한 고루틴이 독식하는 것을 막기위한 임시방편
                time.Sleep(1)
                fmt.Println(i, n)
            }
        }(i)
    }
    // 채널에 0 ~ 9 값을 넣는다.
    for i := 0; i < 10; i++ {
        c <- i
    }
}
```

* 채널을 닫아주지 않으면 프로그램이 종료되지 않는 경우에는 숫자들을 기다리는 3개의 고루틴들이 메모리에 남아있게 된다.
* 메인 함수가 아닌 반복적으로 호출되는 함수 안에서 위와 같은 코드에 채널을 닫는 코드가 없다면 고루틴들은 계속해서 메모리에 쌓이게 되고 이는 memory leak으로 이어진다.
* **채널을 닫는것은 신호를 모두에게 전달하기 위한 매우 강력하고 깔끔한 방법이다.**
* for문 안에 있는 제어 변수를 사용할 때는 **고루틴에 넘겨서 복사 및 고정**시켜 사용해야한다.



### 7.3.3 팬인하기

* 여러 곳에서 수행되던 작업들을 한 곳에 모아 처리하는 것을 FanIn이라 한다.
* 채널을 닫는 것에 주의해야 한다. 보내는 고루틴이 여럿이므로 보내는 곳에서 채널을 닫으면 여러번 닫히기 때문에 패닉이 발생한다.
* 따라서 채널을 닫기 위한 고루틴을 하나 만들고 해당 고루틴에서 모든 고루틴이 데이터를 보내는 것을 기다렸다 채널을 닫는다.

```go
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
```

:bulb: 왜 채널을 닫기 위해 굳이 다른 고루틴을 사용할까?

> 고루틴을 사용하지 않으면 FanIn 함수가 out채널을 반환할 때 함수 안의 리소스가 해제되므로, 클로져 형태로 고루틴을 만들어 FanIn 함수가 종료되더라도 리소스를 가져갈 수 있게 만든다.



### 7.3.4 분산처리

* 팬아웃해서 파이프라인을 통과시키고 다시 팬인시키면 분산처리가 된다.

```go
func Distribute(p IntPipe, n int) IntPipe {
    // IntPipe형을 반환
    return func(in <-chan int) <-chan int {
        // n개로 분산 시킬 채널을 생성
        cs := make([]<-chan int, n)
        // n개로 나누어진 채널에 작업들을 할당해서 진행
        for i := 0; i < n; i++ {
            cs[i] = p(in)
        }
        // 여러개로 나뉘어진 작업들을 다시 FanIn을 이용해서 합침
        return FanIn(cs...)
    }
}
```

* Chain과 함께 사용하면 다양한 파이프라인을 구성할 수 있다.

```go
out := Chain(Cut, Distribute(Chain(Draw, Paint, Decorate), 10), Box)(in)

or 

out := Chain(Cut, Distribute(Draw, 6), Distribute(Paint, 10), Distribute(Decorate, 3), Box)(in)
```

* Go에서는 고루틴마다 스레드를 모두 할당하지 않으므로 고루틴 개수가 많은 것은 크게 걱정할 필요가 없다.
* 동시에 수행될 필요가 없는 고루틴들은 모두 하나의 스레드에서 순차적으로 수행된다.



### 7.3.5 select

switch문과 비슷하지만 동시성 프로그래밍에서 사용되며 다음과 같은 특징이 있다.

* select를 이용하면 **동시에 여러 채널과 통신할 수 있다.**
* 모든 case가 계산된다. 거기에 함수 호출 등이 있으면 select를 수행할 때 모두 호출된다.
* 각 case는 채널에 입출력하는 형태가 되며 막히지 않고 입출력이 가능한 case가 있으면 그중에 하나가 선택되어 입출력이 수행되고 해당 case의 코드만 수행된다.
* default가 있으면 모든 case에 입출력이 불가능할 때 코드가 수행된다. default가 없고 모든 case에 입출력이 불가능하면 어느 하나라도 가능해질 때까지 기다린다.



select는 모든 case를 한번에 검사한다. 또한 select는 입출력이 발생하지 않으면 무한정 대기하게 된다. 만약 10개의 채널에서 값이 들어오는 것을 기다리고 있다가 1개의 채널에서 입/출력이 발생하게 된다면 10개의 case를 모두 검사하고 들어온 채널의 body만 실행하게 된다.



:bulb: 만약 입/출력을 기다리고 있는 select에서 **10개의 채널에서 입/출력이 동시에 발생하게 된다면??**

select는 여러개가 발생했을 때 무작위로 1개의 입/출력만 처리하게 된다. 즉 9개는 입/출력을 처리하지 않는다. 따라서 select를 for문으로 감싸게 된다면 나머지 입출력에 대해서도 계속해서 처리하게 되는 것이다.



### select로 팬인하기

* select를 이용하면 고루틴을 여러 개 이용하지 않고도 팬인을 할 수 있다.

```go
select {
case n := <-c1: c <- n
case n := <-c2: c <- n
case n := <-c3: c <- n
}
```

* c1, c2, c3중 어느 채널이라도 자료가 준비되어 있으면 그것을 c로 보내는 코드
* 셋 중 어떤 채널이 닫혀 있는 경우, 닫혀 있는 채널은 막히지 않고 기본값을 계속해서 받아갈 수 있기 때문에 닫힌 채널로부터 기본값이 받아질 가능성이 있다.

```go
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
```



:bulb: 만약 입/출력을 기다리지 않고 스킵하려면?

select에 default를 추가하게 되면 입/출력이 들어오지 않았을 때 스킵할 수 있게 된다. 

```go
select {
case n := <-c:
    fmt.Println(n)
default:
    fmt.Println("Data is not ready. Skipping...")
}
```



### 시간제한

* 채널과 통신을 기다리되 일정 시간 동안만 기다리겠다면 time.After 함수를 이용하면 된다.
* time.After 함수는 채널을 반환한다.
* 해당 채널로 값이 들어온다면 시간이 넘었다는 뜻이므로 return으로 함수 전체를 빠져나가게 구현하면 된다.

```go
select{
case n := <- recv:
    fmt.Println(n)
case send <- 1:
    fmt.Println("sent 1")
case <-time.After(5 * time.Second):
    fmt.Println("No send and receive communication for 5 seconds")
    return
}
```

* for 문으로 둘러싸여 있다면 매번 타이머가 새로 생성되므로 전체 시간 제한을 두고 싶다면 타이머 채널을 따로 보관하면 된다.

```go
timeout := time.After(5 * time.Second)
select{
case n := <- recv:
    fmt.Println(n)
case send <- 1:
    fmt.Println("sent 1")
case <-timeout:
    fmt.Println("No send and receive communication for 5 seconds")
    return
}
```

* recv와 send에 빈번하게 자료가 반복적으로 오고 가더라도 5초 동안만 처리하게 된다.



### 7.3.6 파이프라인 중단하기

* 파이프라인을 구성할 때 데이터를 끝까지 받지 않고 중간에 끊고 싶을 때도 있다.
* 가장 간단한 방법은 특정 조건이 수행됐을 때 데이터를 받는 for range를 break 하고 남은 값들은 의미 없이 소진시키는 것이다.

```go
func main() {
    c := make(chan int)
    go func() {
        defer close(c)
        for i := 3; i < 103; i++ {
            c <- i
        }
    }()
    
    nums := PlusOne(PlusOne(PlusOne(PlusOne(PlusOne(c)))))
    for num := range nums {
        fmt.Println(num)
        if num == 18 {
            break;
        }
    }
    for _ = range nums {
        // Consume all nums
    }
}
```

* 하지만 이런 방법은 별로 좋지 않은 방법이다.
* 만약 해제되지 않은 고루틴이 존재한다면 memory leak이 발생할 수 있다.
* 또한 많은 네트워크 트래픽을 유발하거나 배터리를 소모한다면 의미없이 데이터를 소모할 때 그만큼의 비용이 발생한다.

:bulb: 그럼 중간에 데이터를 그만 받고 보내는 채널을 닫는 방법은?

>  중간에 데이터를 그만받고 보내는 채널을 닫기 위해서 done채널(보내는 채널을 닫기 위한 채널)을 하나 더 둔다. 보내는 채널에서 done채널을 감시하고 있다가 데이터가 발생되면 보내는 채널을 종료시키는 방식이다. 신호는 close(done)으로 발생시킨다.

```go
func PlusOne(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			select {
			case out <- num + 1:
			case <-done:
				return
			}
		}
	}()
	return out
}

func main() {
    c := make(chan int)
    go func() {
        defer close(c)
        for i := 3; i < 103; i += 10 {
            c <- i
        }
    }()
    done := make(chan struct{})
    nums := PlusOne(done, PlusOne(done, PlusOne(done, PlusOne(done, PlusOne(done, c)))))
    for num := range nums {
        fmt.Println(num)
        if num == 18 {
            break
        }
    }
    close(done)
}
```

* select를 이용해 done채널을 관찰하고 있다 신호가 들어오면 고루틴을 빠져나가도록 만든다.
* close(done) 한 번으로 이 채널로부터 값을 기다리고 있는 모든 고루틴에 일이 끝났다고 방송한다.



### 7.3.7 컨텍스트(context.Context) 활용하기

* 더 복잡한 상황이 발생하면 context 패턴을 사용하는 것이 좋다.
* 종료 신호뿐만 아니라 다른 정보들도 공유되어야 할 때 사용된다.
* 대표적으로 사용자의 인증정보나 요청 마감등이 있다.

```go
func PlusOne(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			select {
			case out <- num + 1:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func main() {
    c := make(chan int)
	go func() {
		defer close(c)
		for i := 3; i < 103; i += 10 {
			c <- i
		}
	}()
    // 취소 기능을 가진 Context를 최상위 Context인 Background에 붙인다.
	ctx, cancel := context.WithCancel(context.Background())
	nums := PlusOne(ctx, PlusOne(ctx, PlusOne(ctx, PlusOne(ctx, PlusOne(ctx, c)))))
	for num := range nums {
		fmt.Println(num)
		if num == 18 {
			cancel()
			break
		}
	}
	// Output:
	// 8
	// 18
}
```

* context.Context는 계층 구조로 되어 있다.
* 최상위 계층은 context.Background() 이며 프로그램이 끝날 때까지 절대로 취소되지 않고 계속 살아있다.
* 하위 구조를 계속해서 트리 구조로 붙일 수 있다.
* 상위 구조가 취소되면 그 하위에 있는 모든 컨텍스트도 취소된다.
* 위 예제에서 ctx는 새로 생성된 컨텍스트가 들어가고, cancel은 해당 컨텍스트를 취소하는데 호출할 수 있는 함수가 들어간다.
* WithDeadline, WithTimeout, WithValue등 다양한 방식의 컨텍스트를 생성할 수 있다.
* 컨텍스트는 관례상 다른 구조체 안에 넣지 않고 함수의 맨 첫 번째 인자로 넘겨주고 받는다.



### 7.3.8 요청과 응답 짝짓기

* 분산처리를 하게 된다면 응답을 받은 고루틴의 순서는 무작위이기 때문에 어떤 요청에 대한 응답인지 알 길이 없다.
* 이를 해결하기 위해 요청과 응답을 짝지어 어떤 요청에서 들어온 응답인지 알 수 있는 방법이다.



간단한 방법은 채널로 자료를 넘겨줄 때 ID번호를 같이 넘겨서 ID 번호를 확인해보는 방법이 있다. 하지만 보내는 쪽에서는 ID를 보관하고 있지만 이 요청에 대한 응답을 다른 고루틴이 받아 갈 수 있다면 복잡해진다.



:bulb: 그렇다면 해결책은?

> 요청을 보낼 때 결과를 받고 싶은 채널을 함께 실어서 보내는 방법을 사용한다. 즉, 요청을 보내는 메시지에 응답을 받을 채널도 함께 넣어서 보낸다.



```go
type Request {
    Num int
    Resp chan Response
}

type Response {
    Num int
    WorkerID int
}

func PlusOneService(reqs <-chan Request, workerID int) {
	for req := range reqs {
        // 요청이 들어왔을 때 각 요청을 처리하는 고루틴을 만든다.
		go func(req Request) {
            // 요청에 대한 처리가 끝나면 response 채널을 닫는다.
			defer close(req.Resp)
            // 요청에 대한 결과를 처리하고 매칭된 응답 채널에 보낸다.
			req.Resp <- Response{req.Num + 1, workerID}
		}(req)
	}
}

func main() {
    reqs := make(chan Request)
	defer close(reqs)

    // 요청을 처리하는 3개의 고루틴을 만듦
	for i := 0; i < 3; i++ {
		go PlusOneService(reqs, i)
	}

	var wg sync.WaitGroup
	for i := 3; i < 53; i += 10 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			resps := make(chan Response)
            // 처리할 값과 응답받을 response를 함께 넘김
			reqs <- Request{i, resps}
            // 응답 채널에서 온 값을 소진한다. 이후 Response 고루틴은 닫힌다.
			log.Println(i, "=>", <-resps)
		}(i)
	}
	wg.Wait()
}
```

* 요청과 응답이 1:1이 아닌 요청에 대한 응답을 0개 이상이 되도록 처리할 수 있다.

```go
for resp := range resps {
    fmt.Println(i, "=>", resp)
}
```

* 이는 검색 요청의 결과처럼 여러 개의 결과를 받아야 하는 경우에 유용하다.



### 7.3.9 동적으로 고루틴 이어붙이기



:heavy_check_mark: prime의 배수를 걸러내는 고루틴을 계속해서 이어붙여 prime number를 출력하는 예제이다. 

* 동적으로 채널을 통하여 고루틴들을 이어붙일 수 있다.

* 2부터 숫자를 하나씩 증가 시켜가며 채널에 숫자를 보낸다. 
* 고루틴에서는 이 채널에서 숫자를 받을 때마다 출력하고, 출력된 숫자의 배수가 되는 숫자들을 걸러내는 필터 고루틴을 붙인다.
* 그러면 이미 출력된 숫자의 배수들은 다시는 출력되지 않게 된다. 따라서 소수만 출력되게 된다.



:bulb: 파이프라인을 일직선으로 이어붙일 때 생성기 패턴으로 필터를 만들고 컨텍스트 등을 이용해 소수 생성기를 만든다.



```go
// start부터 정수를 무한정 생산하는 정수 생성기
// step 만큼 더한 수를 계속해서 생성한다.
// 외부에서 해당 고루틴을 멈출 수 있게 ctx.Done을 검사한다.
func Range(ctx context.Context, start, step int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := start; ; i += step {
			select {
			case out <- i:
			case <-ctx.Done():
				return
			}
		}
	}()
    // 정수를 무한정 생산하는 채널을 반환한다.
	return out
}

// 들어온 n의 배수를 걸러내는 파이프라인을 반환한다.
func FilterMultiple(n int) IntPipe {
    // 클로저를 이용해 다른 파이프라인에 연결해서 쓸 수 있도록 한다.
	return func(ctx context.Context, in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for x := range in {
				if x%n == 0 {
					continue
				}
                // 배수가 아닐경우 채널에 x를 계속해서 넣는다.
				select {
				case out <- x:
				case <-ctx.Done():
					return
				}
			}
		}()
		return out
	}
}

// 무한 소수 생성기
func Primes(ctx context.Context) <-chan int {
    // 소수를 담는 채널을 반환한다.
	out := make(chan int)
	go func() {
		defer close(out)
        // 2부터 정수를 무한정 생성하는 채널
		c := Range(ctx, 2, 1)
		for {
            // <-c가 막혀있을 때 ctx가 취소될 수 있고, out <- i 에서 막혀있다가 
            // ctx가 취소될 수 있기 때문에 select를 다중으로 만들어야 한다.
			select {
			case i := <-c:
				c = FilterMultiple(i)(ctx, c)
				select {
				case out <- i:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func PrintPrimes(max int) {
    // 취소 가능한 컨텍스트 생성
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for prime := range Primes(ctx) {
		if prime > max {
			break
		}
		fmt.Print(prime, " ")
	}
	fmt.Println()
}
```

* Primes에서 소수를 순서대로 꺼내어 이용하다가 범위가 넘어버리면 반복문을 빠져나간다.
* 함수가 끝날 때 derfer cancel()로 인해 Primes에 넘어간 ctx가 취소되고, 생성되었던 고루틴들이 모두 소멸된다.



### 7.3.10 주의점

* 무조건 보내는 쪽에서 채널을 닫아야한다.
* 채널을 닫지 않거나 위험하게 사용하는 실수는 범하지 않도록 정형화된 패턴을 잘 사용해야 한다.



:heavy_check_mark: 잘못 사용하는 예 (생성자-소비자 패턴)

```go
func main() {
    c := make(chan int)
    done := make(chan int)
    // 생성자
    go func() {
        for i := 0; i < 10; i++ {
            c <- i
        }
        done <- true
    }()
	// 소비자
    go func() {
        for {
            fmt.Println(<-c)
        }
    }()
    <-done
}
```

:bulb: 어떤 문제점들이 있을까?

* 두 번째 고루틴이 끝나지 않는다.
    * 첫 번째 고루틴에서 채널을 닫지 않아서 무한정 기다리는 고루틴(좀비)이 되어버린다.
    * 만약 메인함수가 아닌 계속해서 호출되는 함수라면 고루틴이 무한정 쌓일것이므로 memory leak이 발생한다.
* 메인 고루틴에서 <-done으로 인해 소비가 끝나기 전에 메인 고루틴이 끝나버릴 가능성이 있다.
    * fmt.Println(<-c)는 한 줄이더라도 내부적으로는 여러 줄이 실행되기 때문에
    * 이럴 경우 버그가 발생했을 때 이유를 찾기가 어려워진다.



:bulb: 어떻게 해결할까?

* 자료를 보내는 채널에서는 무조건 채널을 닫는다.
* 받는쪽에서 중간에 return 할 수 있으므로 닫을 때는 defer를 이용한다. 그렇지 않으면 닫지 않고 종료될 수 있다.
* 받는 쪽이 끝날 때까지 기다리는 것이 모든 자료의 처리가 끝나는 시점까지 기다리는 방법으로 더 안정적이다.
    * 위 예제는 소비하는 쪽에서 done<-true를 했어야 한다.
    * 위 예제는 소비자쪽에서 생산이 끝났는지 알 수 없으므로 이는 생산자쪽에서 채널을 닫는 것으로 신호를 줬어야했다.
* 특별한 이유가 없으면 소비자쪽은 range를 이용해라.
* WaitGroup을 사용해 두 고루틴이 모두 끝났을 때까지 기다리면 문제가 되질 않는다.
* 끝났음을 알리는 done 채널은 자료를 보내는 쪽에서 결정할 사항이 아니다.
    * 물론 자료를 보내는 쪽에서 채널을 닫아서 자료가 끝났음을 알리는 것이 더 낫다.
    * 소비하는 쪽에서 done 채널에 값을 보내 보내는 쪽에 더 이상 자료를 보내지 말라는 cancel 요청으로 보는게 좋다.
* done 채널에 자료를 보내 신호를 주는 것보다 close(done)으로 채널을 닫아 알리는것이 좋다.



## 7.4 경쟁 상태

공유된 자원에 둘 이상의 프로세스가 동시에 접근하여 잘못된 결과가 나올 수 있는 상태를 경쟁 상태 (race condition) 이라 한다. 이는 타이밍에 따라서 결과가 달라질 수 있기 때문에 고치기 번거로운 버그를 만들어 낸다.



* 복잡한 상황에서 모든 고루틴들이 막혀서 교착 상태(deadlock)가 발생하는 경우가 있다.
    * deadlock은 프로그램이 오류를 출력해주기 때문에 오류를 알 수 있다.
* 하지만 쉽게 발견하지 못하는 버그들이 있다. 그 중 하나가 경쟁 상태 (race condition)이다.
* 채널을 잘 활용하면 경쟁 상태 문제를 많이 해결할 수 있다.
* 몇가지 경우에는 채널보다 sync 라이브러리를 활용하는게 더 간단하다.
* Atomic 라이브러리를 활용하면 해당 연산이 반드시 원자성(atomicity)를 띄게 만들어줄 수 있다.



### 7.4.1 동시성 디버그

* 경쟁 상태 탐지 기능으로 -race 옵션을 줄 수 있다.

```powershell
$ go test -race mypkg		// to test the package
$ go run -race mysrc.go		// to run the source file
$ go build -race mycmd		// to build the command
$ go install -race mypkg	// to install the package
```

* runtime.NumGoroutine()을 호출하여 현재 동작하는 고루틴의 개수를 알 수 있다.
* runtime.NumCpu()로 현재 사용 가능한 CPU수를 알 수 있다.
* runtime.GOMAXPROCS() 얼마나 많은 CPU를 사용할 것인지 통제할 수 있다.



### 7.4.3 sync.Once

* 한 번만 어떤 코드를 수행하고자 할 때 sync.Once를 사용할 수 있다.
* 주로 분산처리를 할 때 initialize용으로 사용할 수 있다.
* 채널을 이용해 같은 효과를 낼 수 있지만 분명한 의미전달을 위해 sync.Once를 사용하자.

```go
func main() {
    var once sync.Once
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            once.Do(func() {
                fmt.Println("Initialized")
            })
            fmt.Println("Goroutine: ", i)
        }(i)
    }
    wg.Wait()
}
// Output:
// initialized
// Goroutine 0
// Goroutine 2
// Goroutine 1
```



### 7.4.4 Mutex와 RWMutex

* 뮤텍스(Mutex)는 상호 배타 잠금 기능이 있다. 
* 동시에 둘 이상의 고루틴에서 코드의 흐름을 제어할 수 있다.
* 뮤텍스를 가장 잘 활용하는 방법은 접근하고자 하는 자원 포인터와 뮤텍스 포인터를 하나의 구조체에 넣어두고 사용하는 방식이다.

```go
type Accessor struct {
    R *Resource
    L *sync.Mutex
}

// 생성
Accessor{&resource, &sync.Mutex{}}

// 사용
func (acc *Accessor) Use() {
    // do something
    acc.L.Lock()
    // Use acc.R
    acc.L.Unlock()
    // Do something else
}
```



* sync.RWMutex는 동시에 읽는것은 허용하지만 한 군데서라도 쓰기를 시도하면 접근할 수 없게 한다.
* 읽는 것은 문제가 되지 않기 때문에 RWMutex를 이용하면 성능상의 장점을 조금 가져갈 수 있다.
* Go의 map에서 RWMutex를 이용하기 적합한 성질을 가지고 있다.



```go
type ConcurrentMap struct {
    M map[string]string
    L *sync.RWMutex
}

func (m ConcurrentMap) Get(key string) string {
    m.L.RLock()
    defer m.L.RUnlock()
    return m.M[key]
}

func (m ConcurrentMap) Set(key, value string) {
    m.L.Lock()
    m.M[key] = value
    m.L.Unlock()
}

func main() {
    m := ConcurrentMap(map[string]string{}, &sync.RWMutex{})
}
```

* RLock과 RUnlock을 사용하지 않고 Lock, Unlock만 사용한다면 Mutex와 동일하다.
* Mutex, RWMutex 모두 sync.Locker 인터페이스를 구현하고 있다.



## 7.5 문맥 전환

문맥 전환(context switching)이란 **프로그램이 여러 프로세스 혹은 스레드에서 동작할 때 기존에 하던 작업들을 메모리에 보관해두고 다른 작업ㅇ르 시작하는 것**을 말한다. 문맥 전환은 비용이 발생하게 된다.



* 하던 일을 중단하고 보관해두어야 하기 때문에 CPU 레지스터에 들어 있던 것들을 메모리에 보관한다.
* 이때 CPU 파이프라인에서 다음 순서에 수행할 완료되지 못한 작업들은 버려지게 된다.



Go 컴파일러는 주로 다음의 경우에 문맥 전환(context switching)을 하는 코드를 생성한다.

* 파일이나 네트워크 연산처럼 시간이 오래 걸리는 입출력 연산이 있을 때
* 채널에 보내거나 받을 때
* go로 고루틴이 생성될 때
* 가비지 컬렉션이 사이클을 지난 뒤



Go에서 고루틴 간의 문맥 전환이 일어나는 것은 자연스럽다. **채널에 자료를 보내는 쪽에서 고루틴이 수행되다가 보내는 코드가 수행될 때, 같은 채널에서 데이터를 받는 다른 고루틴의 해당 부분으로 문맥 전환하면 자연스럽게 수행된다.** 

이것을 컴파일 시간에 예측하면 변수들을 레지스터에 할당하는 전략을 더 잘 세울 수 있어서 성능 향상에 도움이 된다.



* 고루틴을 새로 생성할 때에도 새로 생성된 고루틴으로 건너뛸 수 있다. 그리고 가비지 컬렉션 사이클이 지난 뒤에 문맥 전환이 가능하다.
* time.Sleep(0)을 이용해 문맥 전환을 강제로 시킬 수 있다.