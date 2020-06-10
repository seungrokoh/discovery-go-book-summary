# 4장 함수
서브루틴은 코드 덩어리를 만든 다음 호출하고 귀환할 수 있는 구조이다. 코드를 재사용하여 중복된 코드를 줄일 수 있고, 내부와 외부를 분리하여 생각할 수 있어 코드를 추상화하고 단순화할 수 있다.

* Go 에서 서브루틴을 함수라고 부른다.
* 내부적으로 서브루틴은 주로 스택으로 구현되어있다.
* 호출이 이루엉지면 현재 프로그램 카운터(PC)와 넘겨줄 인자들을 넣은 뒤 PC 값을 변경하여 호출될 서브루틴으로 건너뛴다.
* Go 언어는 값에 의한 호출 (Call by value) 만을 지원한다.
* 주소값을 넘겨받아 그 주소에 있는 값을 변경하여 참조에 의한 호출 (Call by reference)과 비슷한 효과를 낼 수 있다.

## 4.1 값 넘겨주고 넘겨받기
함수를 호출할 때 값을 넘겨주거나 받을 수 있다.

## 4.1.1 값 넘겨주기
```go
func ReadFrom(r io.Reader, lines *[]string) error {}
```
* []string이 아닌 * []string으로 받는 이유는 lines 변수의 값을 변경하고자 하기 때문
* []string을 이용하여 넘겼다면 포인터, 길이, 용량 세 값이 넘어간다.
* 세 값이 넘어가면 세 값을 변경한다 해도 바깥 세상과는 무관한 일이 된다.
* 세 값을 변경하여도 바깥 세상에는 영향을 주지 않지만 포인터가 가리키고 있는 배열을 변경하면 영향을 받는다.

:heavy_check_mark: 포인터 인자 예제
```go
package main

import (
	"fmt"
)

func main() {
	n := []int {1, 2, 3, 4}
	fmt.Println(n)
	addOneCallByValue(n)
	fmt.Println(n)
	addOneCallByReference(&n)
	fmt.Println(n)
}

func addOneCallByValue(nums []int) {
	for i := range nums {
		nums[i]++
	}
	nums = append(nums, 0)
}

func addOneCallByReference(nums *[]int) {
	for i := range *nums {
		(*nums)[i]++
	}
	*nums = append(*nums, 100)
}

// Output:
// [1, 2, 3, 4]
// [2, 3, 4, 5]
// [3, 4, 5, 6, 100]
```

* 슬라이스의 길이나 용량의 변화가 필요 없다면 포인터를 넘겨주지 않아도 값은 변경이 된다.
* 슬라이스의 길이나 용량의 변화가 필요하다면 포인터를 넘겨주어야 변경이 일어날 수 있다.
* 포인터로 넘어온 값은 * 을 앞에 붙여서 값을 참조할 수 있다.
* 포인터 자료형으로 받으면 주소값이 넘어와서 받는 포인터 변수가 담게 되는 것이다.

## 4.1.2 둘 이상의 반환값
* Go의 함수는 둘 이상의 값을 반환할 수 있다.
* 반환값이 둘 이상인 경우 괄호로 둘러싸야 한다.
* 반환할 때에는 쉼표로 구분하여 반환한다.
* 반환값을 받을 때 불필요한 값은 _ 를 이용하여 무시할 수 있다.

## 4.1.3 에러값 주고 받기
* 음수를 반환해 에러를 의미하는 것으로 약속할 수 있음.
* strings.Index()는 어떤 입력에 대해서도 성공하는 함수로 정의됨.
* 변수의 포인터나 레퍼런스를 넘겨 에러값을 받음

Go에서는 에러를 반환함으로써 쉽게 처리할 수 있다.

* Go의 관례상 에러는 가장 마지막 값으로 반환함.
* 패닉(panic)이 있지만 일반적인 에러보다는 심각한 에러 상황에서 쓰임
* 에러값을 돌려주는 방식에 익숙해지자.

```go
if err := MyFunc(); err != nil {}
```

* 조건문 안에서 변수를 새로 만들 수 있어 에러 처리할 때 유용함
* 위의 err 변수는 조건문을 벗어나면 소멸됨
* Go의 에러처리 방식은 현재 문맥에서 처리할 수 없을 때에는 해당 에러를 그대로 반환

:heavy_check_mark: 새로운 에러 생성
```go
return errors.New("stringlist.ReadFrom: line is too long")

or

return fmt.Errorf("stringlist: too long line at %d", count)
```
* 에러를 호출한 곳으로 반환할 때 가장 단순한 방법은 문자열 메세지를 주고 받는 방법
* 부가 정보를 추가하여 돌려줄 때에는 fmt.Errorf()를 사용
* 하지만 넘겨준 숫자는 사람이 보기에는 좋지만 프로그램이 처리하기에는 편리하지 않음. 이건 error 처리에서 배우기

## 4.1.4 명명된 결과 인자
* Go에서 return 값도 변수명을 지정할 수 있음.
* 명명된 return arguments는 기본값으로 초기화됨
* return 뒤에 쉼표로 구분하여 나열할 수도 있고, 생략하고 return만 쓸 수 있다.

```go
func Sum(a, b int) (sum int) {
    sum = a + b
    return
}

func Sum(a, b int) (sum int, err error) {
    sum = a + b
    return
}
```

## 4.1.5 가변 인자
넘겨받는 인자의 개수를 여러개 받을 수 있다.

```go
func WriteTo(w io.Writer, lines... string) (n int64, err error)
```
* lines를 가변인자로 변경하여도 lines는 슬라이스가 된다.
* 원래는 슬라이스를 넘겨줘야 했지만 이제는 나열하는 것으로 호출할 수 있다.

```go
WriteTo(w, "hello", "world", "Go language")
```

* 슬라이스를 넘기기 위해서 뒤에 점 세개를 붙이면 넘길 수 있다.

```go
lines := []string{"hello", "world", "Go language"}
WriteTo(w, lines...)
```
***

## 4.2 값으로 취급되는 함수
Go 언어에서 함수는 일급 시민(First-class citizen)으로 분류된다. 이 뜻은 함수 역시 값으로 변수에 담길 수 있고 다른 함수로 넘기거나 돌려받을 수 있다는 뜻이다.

## 4.2.1 함수 리터럴
* 마치 함수의 이름은 함수의 값을 담는 변수와 같이 보인다.
* add라는 함수는 함수의 이름으로 그 함수를 담고 있는 변수와 같이 보자.

```go
func add(a, b int) int {
    return a + b
}
```
* 여기서 순수하게 함수의 값만 표현하려면? 이름을 빼보자.

```go
func (a, b int) int {
    return a + b
}
```
* 이것을 함수 리터럴(Function literal) 이라고 부르고, 익명 함수라고 부를 수도 있다.
* 함수형 언어에서 람다 함수와 동일한 방법으로 사용할 수 있다.

```go
func printHello() {
    fmt.Println("Hello!")
}
```
```go
func Example_funcLiteral() {
    func() {
        fmt.Println("Hello!")
    }()
    // Output:
    // Hello!
}
```
```go
func Example_funcLiteralVar() {
    printHello := func() {
        fmt.Println("Hello!")
    }
    printHello()
    // Output:
    // Hello!
}
```
* 위의 세 가지 예는 모두 동일하다.
* 함수 리터럴을 호출하기 위해 ()를 붙여 해당 함수를 호출한다.

## 4.2.2 고차 함수
고차 함수(higher-order function)은 함수를 넘기고 받을 수 있는 함수이다.

* 특정 함수가 여러곳에서 쓰일 때 각 사용처마다 처리를 다르게 하고 싶을 때 사용할 수 있다.
* ReadFrom 함수를 여러곳에서 쓰일 수 있도록 고차함수를 이용해서 변경하기

```go
func ReadFrom(r io.Reader, f func(line string)) error {
    scanner := bufio.NewScanner(r)
    for scanner.Scan() {
        f(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        return err
    }
    return nil
}
```
* 이제 r에서 한 줄씩 읽어서 매 줄마다 f 함수를 호출한다.
* 이로써 호출자가 자유롭게 원하는 동작을 f 함수를 통해 지시할 수 있게 된다.

```go
func ExampleReadFrom_Print() {
    r := strings.NewReader("bill\ntom\njane\n")
    err := ReadFrom(r, func (line string) {
        fmt.Println("(", line, ")")
    })
    if err != nil {
        fmt.Println(err)
    }

    // Output:
    // (bill)
    // (tom)
    // (jane)
}
```
:heavy_check_mark: 고차함수 예제
```go
package main

import (
	"fmt"
	"strconv"
)

func main() {
	higherOrderFunction(printItem)
}

func higherOrderFunction(f func (line string)) {
	for _, item := range []int{1, 2, 3} {
		f(strconv.Itoa(item))
	}
}

func printItem(item string) {
	fmt.Println(item)
}

// Output:
// 1
// 2
// 3
```
* 함수 리터럴을 이용해서 함수를 호출할 수 있다.
* ReadFrom에서 읽은 값들을 슬라이스에 추가해 넣고 싶은 경우 클로저를 이용하면 된다.

## 4.2.3 클로저
* 클로저는 외부에서 선언한 변수를 함수 리터럴 내에서 마음대로 접근할 수 있는 코드를 의미한다.

```go
func ExampleReadFrom_append() {
    r := strings.NewReader("bill\ntom\njane\n")
    var lines []string
    err := ReadFrom(r, func(line string) {
        lines = append(lines, line)
    })

    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(lines)

    // Output:
    // [bill, tom, jane]
}
```
* lines라는 문자열 슬라이스를 ReadFrom에 넘겨주는 함수 리터럴 안에서 사용할 수 있다.
* ReadFrom에 넘기는 함수는 그 함수가 이용하는 변수들도 함께 가지고 넘어간다.



클로저에 대해 더 자세히 알아보자. 



클로저는 함수와 그 함수가 선언됐을 때의 렉시컬 환경(Lexical environment)과의 조합이다. 내용을 이해하기 위해 예제를 살펴보자.



```go
func OuterFunc() {
	x := 10
	innerFunc := func() { fmt.Println(x) }
	innerFunc()
}
func main() {
	OuterFunc()
}

// Output:
// 10
```

위 예제에서 OuterFunc 내부에 innerFunc가 선언되고 호출되었다. 이때 innerFunc는 OuterFunc 내부에 선언되었기 때문에 OuterFunc 내부 변수에 접근할 수 있게 된다. 따라서 x의 값인 10이 호출된다.

이제 다른 예제를 보자.



```go
func OuterFunc() func() {
    x := 10
	innerFunc := func() {
		fmt.Println(x)
		x++
	}
	return innerFunc
}

func main() {
    inner := OuterFunc()
    inner()
    inner()
    inner()
}

// Output:
// 10
// 11
// 12
```

OuterFunc는 innerFunc를 반환하고 끝났다. 예상대로라면 함수가 종료됐기 때문에 지역변수인 x는 소멸되어야 한다. 하지만 마치 x가 살아있는 듯이 inner()를 호출하면 10이 출력되고 계속 출력했을 때 값이 증가하는걸 볼 수 있다.



이렇게 자신을 포함한 외부함수보다 내부함수가 더 오래 유지되는 경우, 외부 함수 밖에서 내부함수가 호출되더라도 외부함수의 지역 변수에 접근할 수 있는 함수를 클로저라고 한다.



> 즉, **클로저는 반환된 내부함수가 자신이 선언됐을 때의 환경(Lexical environment)인 스코프를 기억하여 자신이 선언됐을 때의 환경(스코프) 밖에서 호출되어도 그 환경(스코프)에 접근할 수 있는 함수**를 말한다. 이를 조금 더 간단히 말하면 **클로저는 자신이 생성될 때의 환경(Lexical environment)을 기억하는 함수다**라고 말할 수 있겠다.



## 클로저 사용

클로저가 많이 사용되는 유형에 대해서 알아본다.

[reference](https://poiemaweb.com/js-closure)



### 1. 상태 유지

현재 상태를 기억하고 변경된 최신 상태를 기억하는 것이다.

전등을 키고 끌 수 있는 toggle 버튼을 만든다고 가정했을 때 현재 전등의 상태를 나타내기 위하여 전역변수를 선언해야 한다. 하지만 전역변수로 선언하면 외부에서 접근할 수 있고 변경할 수 있기 때문에 오류를 유발시킬 수 있으므로 지양해야 한다.



이때 클로저를 사용하면 최신 상태를 유지할 수 있다.  동작하는 코드는 아니지만 간단하게 전등을 키고 끄는 함수를 클로저로 만든다면 다음과 같다.

```go
func Toggle() func() {
    isOn := false
    return func() {
        // change display ...
        isOn = !isOn
    }
}

func main() {
    toggle := Toggle()
    LightBulb.onClickListener {
        toggle()
    }
}

```



### 2. 전역변수 사용의 억제

버튼이 클릭될 때마다 클릭한 횟수가 누적되는 프로그램을 만든다고 가정해보자. 클릭된 횟수가 유지되어야 하는 상황이다. 이때 전역변수를 사용한다면 오류를 발생시킬 가능성이 있으므로 좋지 않은 코드다. 이를 클로저를 이용하여 작성해보자.



```go
func Increase() func() int {
	counter := 0
	return func() int {
		counter++
		return counter
	}
}
func main() {
	counter := Increase()
    button.onClickListener {
        display(counter())
    }
}
```



## 4.2.4 생성기
* 함수를 호출할 때마다 증가된 값을 받을 수 있는 생성기(generator)를 만들어보자.

```go
func NewIntGenerator() func() int {
    var next int
    return func() int {
        next++
        return next
    }
}

func ExampleNewIntGenerator() {
    gen := NewIntGenerator()

    fmt.Println(gen(), gen(), gen(), gen(), gen())
    fmt.Println(gen(), gen(), gen(), gen(), gen())
    // Output:
    // 1 2 3 4 5
    // 6 7 8 9 10
}
```

* NewIntGenerator는 클로저를 반환하는 고계함수이다.
* 반환하는 함수 리터럴이 속해 있는 스코프 안에 있는 next 변수에 접근하고 있다.
* 따라서 이 함수는 next 변수와 함께 세트로 묶인다.
* 만약 NewIntGenerator()를 여러번 호출하여 함수를 여러개 가지고 있다면 각 함수가 갖고 있는 next도 분리되어 있다.

```go
func ExampleNewIntGenerator_multiple() {
    gen1 := NewIntGenerator()
    gen2 := NewIntGenerator()
    fmt.Println(gen1(), gen1(), gen1())
    fmt.Println(gen2(), gen2(), gen2(), gen2(), gen2())
    fmt.Println(gen1(), gen1(), gen1(), gen1())

    // Output:
    // 1 2 3
    // 1 2 3 4 5
    // 4 5 6 7
}
```
* 같은 방식을 이용해 느긋한 계산법 (Lazy evaluation)을 구현하거나 무한한 크기의 자료구조도 만들 수 있다.

## 4.2.5 명명된 자료형
* 자료형에 새로 이름을 붙일 수 있다.
* type keyword를 사용하면 된다.
```go
type rune int32
```
* 이름만으로 자료형을 지칭하는 것을 명명된 자료형이라한다.
* 이름만으로 자료형을 지칭하지 않는 것을 명명되지 않은 자료형이라한다.
* 자료형을 검사함으로써 프로그램을 직접 수행해보기 전에 컴파일 시점에서 버그를 어느 정도 예방할 수 있다.

정점과 간선으로 이루어진 그래프를 다루는 코드를 예로 들었을 때 각 정점과 간선은 정수형으로 된 ID값을 사용한다고 하자.

:heavy_check_mark: 정점 ID 생성기
```go
func NewVertexIDGenerator() func() int {
    var next int
    return func() int {
        next++
        return next
    }
}
```
:heavy_check_mark: 간선 ID 생성기
```go
func NewEdgeIDGenerator() func() int {
    // ...
}
```
:heavy_check_mark: 간선 ID를 받아 새로운 간선을 생성하는 함수
```go
func NewEdge(eid int) {
    // ...
}
```
```go
func main() {
    gen := NewVertexIDGenerator()
    gen2 := NewEdgeIDGenerator()
    ...
    e := NewEdge(gen())
}
```
* NewEdge 함수에 간선 ID를 넘겨야 하지만 정점 ID를 넘겨도 오류가 발생하지 않는다.
* 찾기 힘든 버그가 발생 할 가능성이 있다.
* 명명된 타입을 이용해 컴파일단에서 실수를 잡아낼 수 있다.

:heavy_check_mark: 정정 및 간선 타입 선언
```go
type VertextID int
type EdgeID int
```

:heavy_check_mark: 명명된 자료형 사용
```go
func NewVertexIDGenerator() func() VertexID {
    var next int
    return func() VertexID {
        next++
        return VertexID(next)
    }
}
```
* int와 VertexID 둘 다 명명된 자료형이기 때문에 호환이 되지 않는다.
* 명명된 자료형끼리는 형변환을 해줘야한다.

```go
func NewEdge(eid EdgeID) {
    // ...
}

func main() {
    gen := NewVertexIDGenerator()
    // ...
    e := NewEdge(gen()) // 컴파일 에러
}
```
* 명명된 자료형을 통해 실수를 미연에 방지할 수 있다.
* 일괄적으로 해당 자료형의 표현을 변경할 수 있다.

```go
type runes []rune

func main() {
    var a []rune = runes{65, 66}
    fmt.Println(string(n))
}
```
* 명명되지 않은 자료형과 명명된 자료형 사이에는 표현이 같으면 호환된다.

## 4.2.6 명명된 함수형
* 함수도 사용자가 자료형을 정의할 수 있다.
* 명명된 함수형도 자료형 검사를 한다.


```go
type BinOp func(int, int) int

func OpThreeAndFour(f BinOp) {
    fmt.Println(f(3, 4))
}

func main() {
    OpThreeAndFour(func (a, b int) int {
        return a + b
    })
    // 컴파일 오류 없음.
}
```

* 함수 리터럴과 명명된 함수형 사이에 자동으로 형변환이 일어난다.
* 하지만 명명된 함수형 사이에는 자동으로 형변환이 일어나지 않는다.

## 4.2.7 인자 고정

```go
type MultiSet map[string]int
type SetOp func(m Multiset, val string)

// Insert 함수는 집합에 val을 추가한다.
func Insert(m Multiset, val string)
```

ReadFrom에서 읽은 각 줄을 지정된 MultiSet에 각 줄을 집어 넣을 수 있다.

```go
m := NewMultiSet()
ReadFrom(r, func(line string) {
    Insert(m, line)
})
```

여기서 함수의 형태를 변환하는 것 또한 추상화가 가능하다.

```go
func InsertFunc(m Multiset) func (val string) {
    return func(val string) {
        Insert(m, val)
    }
}
```

이렇게 하면 호출 부분을 단순화 시킬 수 있다.

```go
m := NewMultiSet()
ReadFrom(r, InsertFunc(m))
```

여기서 한 번 더 일반화 시킬 수 있다.

```go
func BindMap(f SetOp, m MultiSet) func (val string) {
    return func(val string) {
        f(m, val)
    }
}
```

이러면 다음과 같이 호출할 수 있다.

```go
m := NewMultiSet()
ReadFrom(r, BindMap(Insert, m))
```

마치 Insert 함수의 첫 인자인 m을 고정한 함수를 이용하는 것처럼 사용할 수 있다.



## 4.2.8 패턴의 추상화

* 변수로 어떤 값이나 연산의 결과에 대하여 이름을 붙이고 추상화할 수 있다.
* 함수로 코드를 값의 입출력으로 추상화할 수 있다.
* 고차함수를 이용하여 좀 더 높은 수준의 추상화를 할 수 있다.

> 반복되는 패턴이 있다는 것은 추상화할 수 있다는 얘기이다.

```go
func NewIntGenerator() func() int {
    var next int
    return func() int {
        next++
        return next
    }
}
```

```go
func NewVertexIDGenerator() func() VertexID {
    var next int
    return func() VertextID {
        next++
        return VertexID(next)
    }
}
```

위의 두 함수는 결과를 반환하기 전에 VertexID로 형 변환하는 것밖에 다른게 없다. 따라서 재사용이 가능하다.

```go
func NewVertexIDGenerator() func () VertexID {
    gen := NewIntGenerator()
    return func() VertexID {
        return VertexID(gen())
    }
}
```



## 4.2.9 자료구조에 담은 함수

Go에서 함수는 일급시민에 해당하므로 자료구조에 마음대로 담을 수 있다.



***

## 4.3 메서드

* 함수 앞에 리시버가 붙으면 메서드가 된다.
* 리시버 부분의 모양은 인자 목록과 같지만 함수이름 앞에 온다는 것이 다르다.

```go
func (recv T) MethodName(p1 T1, p2 T2) R1
```

* recv에 담겨 있는 자료에 대한 연산들을 메서드로 정의할 수 있다.



## 4.3.1 단순 자료형 메서드

* 모든 명명된 자료형에서 메서드를 정의할 수 있다.

```go
type VertextID int

func (id VertexId) String() string {
    return fmt.Sprintf("VertexId(%d)", id)
}

func ExampleVertexID_print() {
    i := VertexID(100)
    fmt.Println(i)
    // Output:
    // VertexID(100)
}
```



## 4.3.2 문자열 다중 집합

* Multiset을 메서드를 활용하여 자료 추상화를 해보자.

```go
type MultiSet map[string]int

func (m MultiSet) Insert(val string) {
    m[val]++
}

func (m MultiSet) Erase(val string) {
    if m[val] <= 1 {
        delete(m, val)
    } else {
        m[val]--
    }
}

func (m MultiSet) Count(val string) int {
    return m[val]
}

func (m MultiSet) String() string {
    s := "{ "
    for val, count := range m {
        s += strings.Repeat(val + " ", count)
    }
    return s + " }"
}
```

* 메서드를 활용하면 자료 추상화를 할 수 있다.
* 철저히 MultiSet에 집중하여 이 자료를 다루는 메서드들을 정의한 것이다.
* map[string]int로 표현되는 명명된 자료형이 여러 개일 수 있으며 이들은 모두 다른 메서드들을 가질 수 있다.



## 4.3.3 포인터 리시버

* 리시버도 다른 인자들과 같이 값으로 전달된다.
* 자료형이 포인터형인 리시버를 포인터 리시버라고 말한다.
* 포인터로 전달해야 할 경우에는 포인터 리시버를 사용해야 한다.



앞에서 다룬 인접 리스트를 파일에서 읽어오고 파일로 쓰는 부분을 메서드로 변경해보자.



:heavy_check_mark: 기존 코드

```go
func WriteTo(w io.Writer, adjList [][]int) error
func ReadFrom(r io.Reader, adjList *[][]int) error
```

:heavy_check_mark: 변경 코드

```go
type Graph [][]int

func (adjList Graph) WriteTo(w io.Writer) error
func (adjList Graph) ReadFrom(r io.Reader) error
```

* Go언어의 관습상 리시버의 이름을 길게 붙이지 않는다.
* 보통 자료형의 첫 글자를 따서 이름을 붙인다 (ex. func (g Graph) WriteTo)



## 4.3.4 공개 및 비공개

* 식별자 이름의 첫글자가 대문자면 공개, 소문자면 비공개이다.
* 공개일 경우 다른 모듈에서 호출할 수 있으며, 반대로는 호출할 수 없다.
* 자료형, 변수, 상수, 함수, 메서드, 명명된 자료형 등 모두 적용된다.
* 공개된 요소들은 반드시 주석을 달도록 하자. godoc이 문서를 만들어준다.



***

## 4.4 활용

몇 가지 라이브러리 활용법을 익혀보자.

## 4.4.1 타이머 활용하기

* 프로그램의 수행을 잠시 멈추고 싶을 때 time.Sleep 함수를 사용하면 된다.
* time.Sleep 함수는 블로킹(blocking) 타이머이다.
* Time 모듈안에 Timer를 이용하면 넌 블로킹(non-blocking) 타이머를 이용할 수 있다.



:heavy_check_mark: 블로킹 타이머 사용

```go
func CountDown(seconds int) {
    for seconds > 0 {
        fmt.Println(seconds)
        time.Sleep(time.Second)
        seconds--
    }
}

func main() {
    fmt.Println("Ladies and gentlemen!")
    CountDown(5)
}
```

:heavy_check_mark: 넌블로킹 타이머 사용

```go
time.AfterFunc(5 * time.Second, func() {
    // doSomething...
})
```

* 콜백을 구현할 수 있다.
* time.AfterFunc가 호출되면 반환 값으로 타이머를 돌려준다.
* 반환된 타이머를 이용하여 해당 타이머를 stop할 수 있다.



```go
timer := time.AfterFunc(...)
//...
timer.Stop()
```



## 4.2.2 path/filepath 패키지

* 파일 이름 경로를 다루는 패키지
* Walk 함수는 지정된 디렉터리 경로 아래에 있는 파일들에 대하여 어떤 일을 할 수 있는 함수
* 디렉터리 안에 디렉터리가 있으면 그것도 추적해 들어간다
* 음악 파일들을 폴더별로 구분해 놓았다면 이 함수를 이용해 음악들을 모두 재생 목록에 넣는 등의 일을 할 수 있다
* 각 파일에 대하여 어떤 일을 할지는 호출자가 결정할 수 있도록 고차함수로 되어있다.



:heavy_check_mark: Walk 함수

```go
func Walk(root string, walkFn WalkFunc) error
```

:heavy_check_mark: WalkFunc

```go
type WalkFunc func(path string, info os.FileInfo, err error) error
```

