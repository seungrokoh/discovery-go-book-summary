# 3장 문자열 및 자료구조

# 3.1 문자열
* 문자열은 바이트들이 연속적으로 나열되어 있는 것
* string과 []byte로 나타낼 수 있다.
* string은 읽기 전용, []byte는 읽기/쓰기 모두 가능하다.

## 3.1.1 유니코드 처리
* Go 소스 코드는 UTF-8로 되어 있다.
* 파일 내에 모든 문자들이 UTF-8로 자연스럽게 인코딩 된다.

```go
for i, r := range "가나다" {
    fmt.Println(i, r)
}
fmt.Println(len("가나다"))
// Output:
// 0 44032
// 3 45208
// 6 45796
// 9
```

* 각 문자를 그대로 출력하면 유니코드가 출력된다.
* 해당 문자를 찍어보고 싶다면 string(r)을 이용하면 된다.
* 어떤 문자들이 들어있는지를 중시한다면 string, 실제 바이트 표현이 어떤지를 중시한다면 []byte를 쓰는 습관을 붙여보자.

## 3.1.5 문자열 잇기
`문자열 (string)`은 읽기 전용이다. 문자열을 이어 붙이는건 **문자열을 수정하는 것이 아닌 이어붙인 새로운 문자열을 만드는 것이다.**

:heavy_check_mark: 문자열 이어붙이기
> '+' 연산을 이용해 두 문자열을 이어 붙일 수 있다.

:heavy_check_mark: 문자열과 포인터
> Go의 문자열은 문자열에 대한 포인터와 비슷하다.

:heavy_check_mark: 문자열 이어붙이는 다양한 방법
* 간단하게 `'+'`
* `fmt` 패키지의 S로 시작하는 함수들(ex. Sprint, Sprintf)로 이어 붙일 수 있다.
    * 문자열이 아닌 다른 것들도 이어붙일 수 있다.
* `strings` 패키지의 join 함수를 이용
    * 구분자를 이용해 이어붙일 때 용이하다.
* 이 외에 `bytes`, `strconv` 패키지를 살펴보자.

## 3.1.6 문자열을 숫자로
:heavy_check_mark: 문자열에서 숫자로
* strconv.Atoi()
* strconv.ParseXXX()

:heavy_check_mark: 숫자에서 문자열로
* strconv.Itoa()
* strconv.FormatInt()

:heavy_check_mark: fmt 패키지 사용하기
* fmt.Sscanf() : `문자열`로 부터 `숫자` 또는 `다른 형식`을 읽을 수 있다.
* fmt.Sprint() : `숫자`를 `문자열`로 바꿀 수 있다.

***
# 3.2 배열과 슬라이스
`배열`과 `슬라이스` 모두 연속된 메모리 공간을 순차적으로 이용하는 자료구조이다. 주로 `슬라이스`를 이용해 간접적으로 배열을 이용한다.


## 3.2.1 배열
* 연속된 메모리 공간을 순차적으로 사용하는 것
* ...를 이용해 컴파일러가 크기를 알아서 지정하게 할 수 있다.

```go
fruits := [3]string{"사과", "바나나","토마토"}
or
fruits := [...]string{"사과", "바나나", "토마토"}
```

## 3.2.2 슬라이스
* 배열은 자주쓰이지 않는다. 더 유연한 구조의 슬라이스를 사용함
* 슬라이스는 `길이`와 `용량`을 갖고 길이가 변할 수 있는 유연한 구조이다.

```go
var fruits []string
or
var fruits make([]string, n)
```

* make로 만든 슬라이스는 해당 자료형의 `기본값`이 들어감
* 슬라이스를 잘라서 사용할 수도 있음. 이를 `슬라이싱` 이라함
* 범위가 넘어가지 않도록 조심해야함. `슬라이싱`은 세심한 주의가 필요함

:heavy_check_mark: 슬라이싱 예제
```go
func Example_slicing() {
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println(nums)
	fmt.Println(nums[1:3])
	fmt.Println(nums[2:])
	fmt.Println(nums[:3])
	// Output:
	// [1 2 3 4 5]
	// [2 3]
	// [3 4 5]
	// [1 2 3]
}
```

## 3.2.3 슬라이스 덧붙이기
* append 함수를 사용해 붙일 수 있다.
* 가변인자로 여러개를 붙일 수 있다.
* 두 슬라이스를 이어붙이기 위해 ...을 사용해 서로 이어 붙일 수 있다.
* 슬라이싱된 결과를 이용할 수도 있다.

:heavy_check_mark: 슬라이스 덧붙이기 예제
```go
func Example_append() {
	f1 := []string{"사과", "바나나", "토마토"}
	f2 := []string{"포도", "딸기"}
	f3 := append(f1, f2...)     // 이어붙이기
	f4 := append(f1[:2], f2...) // 토마토를 제외하고 이어붙이기

	f4 := append([]string(nil), f2...)
	fmt.Println(f1)
	fmt.Println(f2)
	fmt.Println(f3)
	fmt.Println(f4)
	// Output:
	// [사과 바나나 토마토]
	// [포도 딸기]
	// [사과 바나나 토마토 포도 딸기]
	// [사과 바나나 포도 딸기]
}
```

## 3.2.4 슬라이스 용량
* 슬라이스는 `연속된 메모리 공간`을 활용하는 것이라 `용량에 제한`이 있다.
* 용량이 꽉 찻을 때 덮붙이고자 한다면 `더 넓은 메모리 공간`으로 이사를 가야하며 전에 있던 내용들은 `복사된다.`
* `make([]int, 5)`를 이용해 초기화 했다면 길이뿐만 아니라 `용량도 5로 초기화된다.` 여기서 다시 `append`하게 된다면 새로운 곳으로 복사가 이뤄지게 된다. `(용량 부족)`
* 슬라이스 용량 확인 : `cap(slice)`, 길이 확인 : `len(slice)`

:heavy_check_mark: 슬라이스 용량 예제
```go
func Example_sliceCap() {
	nums := []int{1, 2, 3, 4, 5}

	fmt.Println(nums)
	fmt.Println("len:", len(nums))
	fmt.Println("cap:", cap(nums))
	fmt.Println()

	sliced1 := nums[:3]
	fmt.Println(sliced1)
	fmt.Println("len:", len(sliced1))
	fmt.Println("cap:", cap(sliced1))
	fmt.Println()

	sliced2 := nums[2:]
	fmt.Println(sliced2)
	fmt.Println("len:", len(sliced2))
	fmt.Println("cap:", cap(sliced2))
	fmt.Println()

	sliced3 := sliced1[:4]
	fmt.Println(sliced3)
	fmt.Println("len:", len(sliced3))
	fmt.Println("cap:", cap(sliced3))
	fmt.Println()

	nums[2] = 100
	fmt.Println(nums, sliced1, sliced2, sliced3)
	// Output:
	// [1 2 3 4 5]
	// len: 5
	// cap: 5
	//
	// [1 2 3]
	// len: 3
	// cap: 5
	//
	// [3 4 5]
	// len: 3
	// cap: 3
	//
	// [1 2 3 4]
	// len: 4
	// cap: 5
	//
	// [1 2 100 4 5] [1 2 100] [100 4 5] [1 2 100 4]
}
```
* 슬라이스를 `뒤에서 자를 경우` 용량은 그대로다.
* 슬라이스를 `앞에서 자를 경우` 용량은 줄어든다.
* 잘라냈더라도 `뒤에 공간이 있으면 그 공간을 살릴 수도 있다.`
* 슬라이싱을 했을 때 `메모리 재 할당이 이루어지지 않는 이상 동일 메모리를 바라보고 있다.`
* 초기화 때 용량도 미리 할당할 수 있다.

:bulb: append시 원소 개수만큼 capacity가 증가하는것이 아니다.
>  append 할 때 cap은 1씩 증가하는 것이 아닌 go version 마다 다르게 증가한다.

:heavy_check_mark: 용량 할당
```go
nums := make([]int, 3, 5)
```
길이가 3, 용량이 5인 슬라이스가 생성된 것이다. 이는 아래 코드와 동일하다.
```go
nums := make([]int, 5)
nums = nums[:3]
```

* 빈 슬라이스를 만들지만 `공간을 미리 할당받을 수 있다.`
* 공간을 미리 할당 받으면 `원소가 해당 용량을 초과하지 않는 이상 복사가 이루어지지 않기 때문에` 효율적이다.

:heavy_check_mark: 공간 미리 할당받기
```go
nums := make([]int, 0, 15)
```

## 3.2.5 슬라이스의 내부 구현
* 슬라이스는 `배열을 가리키고 있는 구조체`이다.
* 슬라이스는 `시작 주소`, `길이`, `용량` 3개의 필드로 구성되어 있다.
* 복사가 일어나서 이동이 일어난다면, `새로운 배열을 보고` 있게 된다. `배열은 크기가 변경될 수 없기 때문에` 크기가 다른 배열을 하나 더 만들어야 하기 때문이다.
* append로 원소 추가시 일단 nums의 `늘어난 길이가 용량을 초과할 것인지 아닌지를 먼저 조사`한다.
* 용량이 초과하지 않을 경우 시작위치에서 길이만큼 오른쪽으로 이동 후 길이가 증가한 슬라이스를 반환한다.
* 용량을 초과 할 경우 `더 큰 크기의 배열을 새로 하나 더 만들고` 슬라이스도 거기에 맞게 고쳐서 반환한다.

## 3.2.6 슬라이스 복사
* copy(dest, src) 함수를 이용하여 쉽게 복사할 수 있다.
```go
copy(dest, src)
```

* `src 길이 만큼 복사`하게 된다.
* dest의 길이가 src보다 작다면 `일부분만 복사될 수 있다.`
* len(src), len(dest) 중 작은 값만큼 복사가 이루어 진다.
* copy 함수의 `반환값은 복사된 원소의 개수`이다.
```go
if n := copy(dest, src); n != len(src) {
    fmt.Println("복사가 덜 됐습니다.")
}
```

* append 를 이용해서도 복사가 가능하다.
```go
src := []int{30, 20, 50, 10, 40}
dest := append([]int(nil), src...)
```

## 3.2.7 슬라이스 삽입 및 삭제
배열에서 삽입과 삭제는 메서드를 이용하여 호출하면 단순해보이지만 실제로는 연속된 공간에 있어야 하기 때문에 굉장히 비효율적인 과정을 거칠 수 밖에 없다.

:heavy_check_mark: 슬라이스에 원소 1개 삽입
```go
if i < len(a) {
    a = append(a[:i + 1], a[i:]...)
    a[i] = x
} else {
    a = append(a, x)
}

or

a = append(a, x)
copy(a[i + 1:], a[i:])
a[i] = x
```

:heavy_check_mark: 슬라이스에 원소 여러개 삽입
```go
// x := []int{7, 8, 9}
a = apeend(a, x...)
copy(a[i + len(x):], a[i:])
copy(a[i:], x)
```

:heavy_check_mark: 슬라이스에서 원소 1개 삭제
```go
a = append(a[:i], a[i + 1:]...)
```

:heavy_check_mark: 슬라이스에서 원소 여러개 삭제
```go
a = append(a[:i], a[i + k:]...)
```

* 위 방법으로 원소를 삭제하면 O(n)의 시간이 걸린다.
* 순서가 보장되지 않는 방법으로 O(1)의 시간으로 원소를 삭제할 수 있다.

:heavy_check_mark: O(1)로 원소 삭제
```go
a[i] = a[len(a) - 1]
a = a[:len(a) - 1]
```

* 슬라이스에서 삭제시 삭제되는 슬라이스 내부에 포인터가 있다면 메모리 릭이 일어난다.

***
# 3.3 맵
* Go에서 map은 해시테이블로 구현된다.
* 키를 이용하여 값을 상수 시간에 가져올 수 있다.
* 해시맵은 순서가 보장되지 않는다.
* 맵을 생성해야지만 사용할 수 있다.

```go
m := make(map[keyType]valueType)

or

m := make(map[keyType]valueType{})
```

* m[key]를 이용하여 맵의 값을 읽을 수 잇다.
* 해당 키가 없으면 값의 자료형의 기본값을 반환한다.
* 맵을 읽을 때 두개의 변수로 받게 되면, 두 번째 변수에 키의 존재 여부를 받을 수 있다.

```go
value, exist = m[key]
```

* 맵에 키가 이미 존재한다면 기존에 있던 값이 변경되고 없는 경우에는 새로 만든다.
```go
m[key] = value
```

* 맵은 순서가 보장되지 않기 때문에 순서에 의존적인 코딩을 하면 안된다.
* 해시 함수가 바뀔 수 있기 때문이다.
* 맵 테스트를 할 때 비교하는 방법은 reflect.DeepEqual 을 이용하면 된다.

***
# 3.4 입출력
Go의 입출력에 대한 표준 라이브러리는 io package에 들어있다. fmt package에 형식을 이용한 입출력도 구현되어 있다.

## 3.4.1 io.Reader와 io.Writer
* 입출력은 io.Reader와 io.Writer 인터페이스와 파생된 다른 인터페이스들을 이용한다.
* fmt package에서 F로 시작하는 함수들이 io.Reader, io.Writer를 인자로 받는다.
* 파일뿐만 아니라 버퍼, 소켓 등을 이용해서 읽고 쓸 수 있다
* 따라서 기본 입출력 역시 파일을 읽고 쓰는 것과 거의 동일한 방법으로 읽고 쓸 수 있다.
* 함수를 작성할 때 io.Reader, io.Writer 등을 받아서 처리하게 만드는 습관을 들이자.

> io.Reader, io.Writer를 이용하면 표준 입출력, 파일, 네트워크 등 모두 적용할 수 있으며 테스트 등을 할 때도 좋다.

## 3.4.2 파일 읽기
```go
f, err := os.Open(filename)
if err != nil {
    return err // 혹은 다른 에러 처리
}
defer f.Close()

var num int
if _, err := fmt.Fscanf(f, "%d\n", &num); err == nil {
    // 읽은 num 값 사용
}
```

* is.Open()은 파일 오브젝트, 에러 두 가지 반환값이 있다.
* defer는 해당 함수를 벗어날 때 호출할 함수를 등록하는 역할을 한다.
* open 이후 defer로 Close() 하는 습관을 들이자.

## 3.4.3 파일 쓰기
```go
f, err := os.Create(filename)
if err != nil {
    return err // 혹은 다른 에러 처리
}
defer f.Close()
if _, err := fmt.Fprintf(f, "%d\n", num); err != nil {
    return err // 혹은 다른 에러 처리
}
```

## 3.4.4 텍스트 리스트 읽고 쓰기
:heavy_check_mark: 문자열 슬라이스 라인별로 출력
```go
func WriteTo(w io.Writer, lines []string) error {
    for _, line := range lines {
        if _, err := fmt.Fprintln(w, line); err != nil {
            return err
        }
    }
    return nil
}
```

:heavy_check_mark: 문자열 읽기
```go
func ReadFrom(r io.Reader, lines *[]string) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		*lines = append(*lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
```

:heavy_check_mark: 읽기 및 쓰기 사용
```go
func ExampleWriteTo() {
	lines := []string{
		"bill@mail.com",
		"tom@mail.com",
		"jane@mail.com",
	}
	if err := WriteTo(os.Stdout, lines); err != nil {
		fmt.Println(err)
	}
	// Output:
	// bill@mail.com
	// tom@mail.com
	// jane@mail.com
}

func ExampleReadFrom() {
	r := strings.NewReader("bill\ntom\njane\n")
	var lines []string
	if err := ReadFrom(r, &lines); err != nil {
		fmt.Println(err)
	}
	fmt.Println(lines)
	// Output:
	// [bill tom jane]
}
```

## 3.4.5 그래프의 인접 리스트 읽고 쓰기

:heavy_check_mark: 출력함수
```go
func WriteTo(w io.Writer, adjList [][]int) error {
	size := len(adjList)
	if _, err := fmt.Fprintf(w, "%d", size); err != nil {
		return err
	}
	for i := 0; i < size; i++ {
		lsize := len(adjList[i])
		if _, err := fmt.Fprintf(w, "\n%d", lsize); err != nil {
			return err
		}
		for j := 0; j < lsize; j++ {
			if _, err := fmt.Fprintf(w, " %d", adjList[i][j]); err != nil {
				return err
			}
		}
	}
	if _, err := fmt.Fprintf(w, "\n"); err != nil {
		return err
	}
	return nil
}
```

:heavy_check_mark: 입력함수
```go
func ReadFrom(r io.Reader, adjList *[][]int) error {
	var size int
	if _, err := fmt.Fscanf(r, "%d", &size); err != nil {
		return err
	}
	*adjList = make([][]int, size)
	for i := 0; i < size; i++ {
		var lsize int
		if _, err := fmt.Fscanf(r, "\n%d", &lsize); err != nil {
			return err
		}
		(*adjList)[i] = make([]int, lsize)
		for j := 0; j < lsize; j++ {
			if _, err := fmt.Fscanf(r, " %d", &(*adjList)[i][j]); err != nil {
				return err
			}
		}
	}
	if _, err := fmt.Fscanf(r, "\n"); err != nil {
		return err
	}
	return nil
}
```

:heavy_check_mark: 테스트 함수
```go
package graph

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestWriteTo(t *testing.T) {
	adjList := [][]int{
		{3, 4},
		{0, 2},
		{3},
		{2, 4},
		{0},
	}
	w := bytes.NewBuffer(nil)
	if err := WriteTo(w, adjList); err != nil {
		t.Error(err)
	}
	expected := "5\n2 3 4\n2 0 2\n1 3\n2 2 4\n1 0\n"
	if expected != w.String() {
		t.Logf("expected: %s\n", expected)
		t.Errorf("found: %s\n", w.String())
	}
}

func ExampleReadFrom() {
	r := strings.NewReader("5\n2 3 4\n2 0 2\n1 3\n2 2 4\n1 0\n")
	var adjList [][]int
	if err := ReadFrom(r, &adjList); err != nil {
		fmt.Println(err)
	}
	fmt.Println(adjList)
	// Output:
	// [[3 4] [0 2] [3] [2 4] [0]]
}
```
