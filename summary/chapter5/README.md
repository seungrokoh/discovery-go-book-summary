# 5장 구조체 및 인터페이스

* 구조체는 필드들을 묶어 놓은 것이다.
* 더 복잡한 자료형을 정의할 수 있다. 
* 자료를 네트워크를 통하여 전송하거나 파일에 저장하고 불러오는 경우 직렬화 및 역직렬화를 활용해야 한다.
* JSON은 이런 직렬화 및 역질렬화 형식 중 하나이다.
* 인터페이스는 메서드의 집합으로 구현은 없고 형태만 존재한다.
* 인터페이스의 메서드를 모두 정의하기만 하면 해당 인터페이스를 구현한 것으로 취급된다.
* 인터페이스는 프로그램을 유연하게 해주며, 외부 의존성을 줄일 수 있는 방법 중 하나이다.
* 빈 인터페이스는 와일드카드와 같은 존재로 형 스위치를 이용해 활용된다.



## 5.1 구조체

* 구조체는 필드을의 모음 혹은 묶음을 말한다.
* 명명된 구성 요소들을 필드라고 한다.
* 서로 다른 자료형의 자료들도 묶을 수 있다.



## 5.1.1 구조체 사용법

```go
var task = struct {
    title string
    done bool
    due *time.Time
} {"laundry", false, nil}
```

* 구조체도 이름을 붙일 수 있다.

```go
type Task struct {
    title string
    done bool
    due *time.Time
}
```

```go
var myTask = Task{"laundry", false, nil}
```

* 필드명과 함께 입력하면 원하는 값만 넣을 수도 있다.
* 넣지 않은 필드들은 기본값으로 설정된다.
* 여러 필드의 값을 정할 때는 쉼표로 구분한다.

```go
var myTask = Task{title: "laundry", done: true}
```

* 보통 한 줄보다는 여러줄로 선언한다.
* 가독성이 좋아지고 수정하기에 편하다.
* 마지막 필드에 쉼표를 붙여줘야 한다.

```go
var myTask = Task {
    title: "laundry", 
    done: true,
}
```



## 5.1.2 const와 iota

* 확장성을 위해 bool형을 쓸 곳에 enum 형을 쓰는 것을 고려해보자.
* 더 많은 종류의 상태가 추가될 때 쉽게 변경할 수 있다.

```go
type status int

type Task struct {
    title string
    status status
    due *time.Time
}
```

* Go에서는 enum이 따로 없고 상수로 정의해서 사용한다.
* 서로 연관된 상수는 묶어서 정의할 수 있다.

```go
const (
	UNKNOWN status = 0
    TODO	status = 1
    DONE	status = 2
)
```

* 0, 1, 2를 순서대로 붙일 때는 iota를 쓰면 편리하다.
* 첫 번째에만 iota를 써주면 아래 변수들은 같은 형이 되며 숫자가 순서대로 붙는다

```go
const (
	UNKNOWN status = iota
    TODO
    DONE
)
```

* 상수 이름을 밑줄(_)을 쓰게 되면 무시된다.
* iota의 값이 0부터 증가하는 것을 이용해 다양한 방법으로 iota를 활용할 수 있다.

```go
type ByteSize float64

const (
	_			= iota
    KB ByteSize = 1 << (10 * iota)
    MB
    GB
    TB
    PB
    EB
    ZB
    YB
)
```



## 5.1.3 테이블 기반 테스트

* Go언어에서는 다른 언어에서 지원하는 assertEqual과 같은 함수를 제공하지 않는다.
* if 문을 사용하여 각 자료형에 대해 비교한다.
* 여러 사례를 테스트 하고 싶을 때 구조체와 배열을 이용하여 테이블기반 테스트를 할 수 있다.
* 입력과 출력의 목록들만 따로 빼내 구조체 배열을 만들어 사용한다.
* 어느 테스트 케이스에서 문제가 발생하는지 눈에 잘 띄기 위하여 인덱스 번호를 같이 출력해주면 좋다.

```go
func TestFib(t *testing.T) {
    cases := []struct {
        in, want int
    } {
        // 0. case description...
        {0, 0},
        // 1. case description...
        {5, 5}, 
        // 2. case description...
        {6, 8},
    }
    
    for i, c := range cases {
        got := seq.Fib(c.in)
        if got != c.want {
            t.Errorf("Case %d, Fib(%d) : %d, wantß : %d", i, c.in, got, want)
        }
    } 
}
```

* assertion 함수를 직접 만들기 전에 단순히 if를 활용하면 어떤지, 테이블 기반 테스트를 이용하면 어떤지 한 번 생각해보자.
* Go 커뮤니티 관습 역시 if를 이용하거나 테이블 기반 테스트를 하는 것을 권장하고 있다.
* 실제 fmt 패키지의 Sprintf 함수도 테이블 기반 테스트를 활용한다. [link](https://golang.org/src/fmt/fmt_test.go)



## 5.1.4 구조체 내장

* 구조체는 여러 자료형의 필드들을 가질 수 있다는 점이 장점이다.
* 이 뜻은 구조체 안에 구조체 필드를 가질 수 있다는 뜻이 된다.
* 구조체를 내장하여 재사용할 수 있다.
* 구조체 내장을 활용하여 포함 관계(Has A)를 사용할 수 있다.



마감 시간을 표현하기 위해 Deadline 자료형을 선언해보고 마감시간을 체크하는 함수를 만들어보자.

데드라인이 없는 경우도 구현하기 위해 리시버 포인터를 사용한다.

```go
type Deadline time.Time

func (d *Deadline) OverDue() bool {
    return d != nil && time.Time(*d).Before(time.Now())
}
```

```go
type Task struct {
    Title string
    Status status
    Deadline *Deadline
}

// OverDue returns true if the deadline is before the current time.
func (t Task) OverDue() bool {
    return t.Deadline.OverDue()
}
```

* Task의 OverDue()는 그저 Deadline에 자기가 할 일을 위임 할 뿐이다.
* 메서드마다 모두 같은 이름의 메서드를 호출하는 코드를 작성해야 하는 일은 매우 비 효율적이다.
* 내장 기능을 이용하여 이 문제를 해결할 수 있다.



```go
type Task struct {
    Title string
    status status
    *Deadline
}
```

* 필드 이름을 생략하면 내장된다.
* 자료형의 이름과 같은 필드를 갖게 된다. (위에서는 Deadline 필드가 생김)
* 정의되어 있는 메서드도 바로 호출할 수 있는 상태가 된다.
* 따라서 기존에 위임하기만 하던 Task의 불필요한 메서드는 삭제 가능하다.
* 필드가 내장되어 있으면 내장된 필드가 구조체 전체의 직렬화 결과를 바꾸는 문제가 있다는 점을 명심하자.
* 내장 할 필드를 구조체로 Wrapping해서 문제를 해결해보자



```go
type Deadline struct {
    time.Time
}

func NewDeadline(t time.Time) *Deadline {
    return &Deadline{t}
}

type Task struct {
    Title string
    Status status
    Deadline *Deadline
}
```

* 구조체를 내장하게 되면 내장된 구조체에 들어있는 필드들도 바로 접근 가능하다.
* 따라서 여러 구조체에 있는 필드들이 모두 합쳐진 구조체 같은 것을 만들 수 있다.
* 상속과는 달리 실제로는 내부에 필드를 내장하고 있으면서 편의를 제공하는 것뿐이라는 걸 명심하자.



## 5.2 직렬화와 역직렬화

* 직렬화(Serialization)란 객체의 상태를 보관이나 전송 가능한 상태로 변환하는 것
* 역직렬화(Deserialization)란 보관되거나 전송받은 것을 다시 객체로 복원하는 것
* 보조 기억장치에 저장 및 불러오기, 네트워크를 통한 메세지 전송 등 다양한 곳에서 직렬화가 필요하다.



## 5.2.1 JSON

* 자료 교환 형식 중 하나로 XML에 비해 사람이 읽기 쉽고 간단하기 때문에 널리 사용되는 방식이다.
* 키와 값이 대응되어 있다.
* 값에는 불리언, 숫자, 문자열 뿐만 아니라 객체나 배열도 들어갈 수 있다.

```json
{
    "Title": "Laundry",
    "Status": 2,
    "Deadline": "2016-11-10T23:00:00Z"
}
```



### JSON 직렬화 및 역직렬화

* json.Marshal 함수를 이용하여 직렬화를 할 수 있다.
* JSON 패키지는 대문자로 시작하는 필드들만 JSON으로 직렬화한다.
* 직렬화를 하고 싶지 않은 필드들은 소문자로 시작하게 만들면 된다.

```go
func Example_marshalJSON() {
	t := Task {
		"Laundry",
		DONE,
		NewDeadline(time.Date(2015, time.August, 16, 15, 43, 0, 0, time.UTC)),
	}
	b, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(b))
}
```

* Json.Unmarshal 함수를 이용하여 역직렬화 할 수 있다.
* 포인터를 이용해야 수정된 것이 반영되므로 구조체 포인터를 넘겨준다.

```go
func Example_unmarshalJSON() {
	b := []byte(`{"Title":"Buy Milk","Status":2,"Deadline":"2015-08-16T15:43:00Z"}`)
	t := Task{}
	err := json.Unmarshal(b, &t)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(t.Title)
	fmt.Println(t.Status)
	fmt.Println(t.Deadline.UTC())
}
```



### JSON 태그

* 통신시 직렬화 필드들의 이름을 태그를 이용해 지정할 수 있다.
* 구조체의 필드에 json 태그를 붙일 수 있고 JSON 라이브러리가 읽고 처리해준다.

```go
type MyStruct struct {
    Title 		string 	`json:"title"`
    Internal 	string 	`json:"-"`
    Value 		int64	`json:",omitempty"`
    ID			int64	`json:",string"`
}
```

* json:title 로 태깅함으로써 "title"을 JSON 필드로 사용한다.
* 대시(-)를 이용해 필드를 무시할 수 있다.
* omitempty를 이용해 값이 비어있으면 JSON 결과를 나타내지 않을 수 있다.
* JSON에서 string으로 나타내게 할 수 있다.

마지막 ID는 Go 에서만 해당 json을 읽을 것 같으면 그냥 int64로 나둬도 된다. 하지만 웹 어플리케이션을 작성하는 경우 등 자바스크립트가 json을 읽게 되면ㅇ문제가 발생한다.



* 64비트 정수를 json으로 주고 받을 떄에는 \`json:",string"\`을 해주는 습관을 기르도록 하자.



### JSON 직렬화 사용자 정의

* 명명된 자료형은 marshal, unmashal시 원하는 값을 받을 수 있도록 사용자가 직접 정의할 수 있다.

```go
func (s status) MarshalJson() ([]byte, error) {
	switch s {
	case UNKNOWN:
		return []byte(`"UNKNOWN"`), nil
	case TODO:
		return []byte(`"TODO"`), nil
	case DONE:
		return []byte(`"DONE"`), nil
	default:
		return nil, errors.New("status.MarshalJson: unknown value")
	}
}

func (s *status) UnmarshalJson(data []byte) error {
	switch string(data) {
	case `"UNKNOWN"`:
		*s = UNKNOWN
	case `"TODO"`:
		*s = TODO
	case `"DONE"`:
		*s = DONE
	default:
		return errors.New("status.MarshalJson: unknown value")
	}
	return nil
}
```

* go generate 도구를 이용해 반복적인 코드를 자동으로 생성할 수 있다.
* 역 따옴표로 한 번 더 둘러싼 이유는 넘어오는 데이터가 따옴표까지 포함된 문자열이기 때문이다.
* json에서는 따옴표 없이 수를 표현하고 따옴표를 넣어서 문자열을 표현한다.
* ID 값과 같이 정확한 값의 전달이 필요하거나 시간 값이라고 해도 열화가 일어나지 않아야 하는 경우 \`json:,string\` 태그를 붙여서 문자열로 전달해야 한다.

```go
func (d Deadline) MarshalJSON() ([]byte, error) {
	return strconv.AppendInt(nil, d.Unix(), 10), nil
}

func (d *Deadline) UnmarshalJSON(data []byte) error {
	unix, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	d.Time = time.Unix(unix, 0)
	return nil
}
```



* JSON 관례상 키 이름은 소문자로 시작하는 것이 일반적이다.

```go
type Task struct {
    Title	string		`json:"title,omitempty"`
    Status	status 		`json:"status,omitempty"`
    Deadline *Deadline	`json:"deadline,omitempty"`
    Priority int 		`json:"priority,omitempty"`
}
```



### 구조체가 아닌 자료형 처리

* 반드시 구조체를 이용하여 JSON 라이브러리를 이용할 필요는 없다.
* 배열을 직렬화 및 역직렬화를 하는데 JSON 라이브러리를 사용할 수 있다.
* 그러나 여러 필드가 있는 자바스크립트의 오브젝트를 처리할 때 구조체가 매우 적절하다.
* 맵을 이용하여 자바스크립트의 오브젝트를 처리할 수 있다.

```go
func Example_mapMarshalJson() {
    b, _ := json.Marshal(map[string]string {
        "Name": "John", 
        "Age": "16",
    })
    fmt.Println(string(b))
    // Output:
    // {"Age": "16", "Name":"John"}
}
```

* 맵은 순서가 없기 때문에 Name과 Age 누가 먼저 나올지 모른다.
* 하지만 JSON 라이브러리의 구현에서는 키를 정렬하여 출력해준다.
* JSON에 이용하는 맵은 키가 문자열형이어야 한다.
* 아무 자료형을 담으려면 interface{} 자료형을 쓰면 된다.

```go
func Example_mapMarshalJson() {
    b, _ := json.Marshal(map[string]interface{} {
        "Name": "John", 
        "Age": 16,
    })
    fmt.Println(string(b))
    // Output:
    // {"Age": 16, "Name":"John"}
}
```

* interface{} 는 맵도 담을 수 있기 때문에 중첩 JSON 오브젝트를 나타낼 수 있다.
* interface{} 는 JSON으로 ㅛ현된 모든 데이터를 역직렬화할 수 있다.
* JSON 역직렬화 중 interface{} 자료형을 만나면 map[string]interface{} 형태로 역직렬화 하기 때문이다.



### JSON 필드 조작하기

* MarshalJSON 및 UnmarshalJSON을 구현해주는 것으로 쉽게 해결되지 않는 경우들이 많다.
* 입출력하는 JSON의 구조에 따라 구조체의 구조가 제한되어 버린다.



> 구조체 내장을 이용하여 여러 문제를 해결할 수 있다.



* 첫 번째로 구조체에서 특정 필드들을 빼고 직렬화하고 싶은 경우에 쓸 수 있다.
    * 해당 필드가 제거된 구조체 자료형을 만든 다음 그 자료형으로 모든 자료를 복사한 다음에 직렬화 하는방법
    * 이 방법은 번거롭고, 필드가 추가될 때마다 복사하는 코드도 변경해야 하는 단점이 존재한다.
    * 다른 방법은 빼고자 하는 필드들의 값을 지운 다음 직렬화하는 방법
    * 이 경우에는 JSON 태그를 변경할 수 없기 때문에 omitempty를 하고 싶지 않은 경우에는 제외시킬 방법이 없다.
* 구조체 내장을 이용하면 위에 나열된 단점들을 어느정도 해소할 수 있다.
    * 원래 구조체를 고치지 않고, 원하는 필드들만 제외하거나 추가하여 직렬화할 수 있다.



```go
type Fields struct {
    VisivleField	string
    InvisibleField	string
}

func ExampleOmitFields() {
    f := &Fields{"a", "b"}
    b, _ := json.Marshal(struct {
        *Fields
        InvisibleField 	string	`json:",omitempty"`
        Additional 		string
    } {
        Fields: f,
        Additional: "c"
    })
    fmt.Println(string(b))
    // Output: {"VisibleField":"a", "Additional":"c"}
}
```

* 직렬화시 struct에 Fields를 내장한다.
* 빼고 직렬화하고 싶은 필드 InvisibleField를 함께 넣어주고 omitempty 태그를 넣어준다.
* 원래 구조체와 빼고자 하는 구조체 모두 같은 필드 이름으로 해줘야한다.

```go
type Fields struct {
    VisibleField	string	`json:"visibleField"`
    InvisibleField 	string	`json:"invisibleField"`
}

InvisibleField	string	`json:"invisibleField,omitempty"`
Additional 		string 	`json:"additional,omitempty"`
```

* 두 구조체를 필드들을 합집합 할 수 있다.

```go
b, _ := json.Marshal(struct {
        *Fields
        *AnotherStruct
    } {
        // ...
    })
```



### 5.2.2 Gob

* Go 언어에서 기본으로 제공하는 또 다른 직렬화 방식이다.
* Go 언어에서만 읽고 쓸 수 있는 형태이다.
* 더 효율적인 변환이 가능하다.
* 주고받는 코드가 모두 Go로 되어 있을때 Gob 이용을 고려해볼 수 있다.
* gob.NewEncoder, gob.NewDecoder로 인코더와 디코더를 생성해 사용한다.
* 이때 각각 io.Writer, io.Reader를 넘긴다. byte.Buffer를 넘기면 []byte를 뽑아낼 수 있다.

맵을 인코딩 한 다음 한 줄에 16바이트씩 16진수로 출력하고 다시 맵으로 복원하는 예제

```go
package gob

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func Example_gob() {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	data := map[string]string{"N": "J"}
	if err := enc.Encode(data); err != nil {
		fmt.Println(err)
	}
	const width = 16
	for start := 0; start < len(b.Bytes()); start += width {
		end := start + width
		if end > len(b.Bytes()) {
			end = len(b.Bytes())
		}
		fmt.Printf("% x\n", b.Bytes()[start:end])
	}
	dec := gob.NewDecoder(&b)
	var restored map[string]string
	if err := dec.Decode(&restored); err != nil {
		fmt.Println(err)
	}
	fmt.Println(restored)
	// Output:
	// 0e ff 81 04 01 02 ff 82 00 01 0c 01 0c 00 00 08
	// ff 82 00 01 01 4e 01 4a
	// map[N:J]
}
```



## 5.3 인터페이스

* 인터페이스는 메서드들의 묶음이다.
* 인터페이스 네이밍은 주로 인터페이스의 메서드 이름에 er을 붙인다.
* 인터페이스에 있는 메서드들을 모두 구현한 자료형은 이 인터페이스로 사용할 수 있다.
* 예로 MarshalJSON을 구현하기만 해도 json.Marshaler가 될 수 있다.



## 5.3.2 인터페이스 정의

```go
interface {
    Method1()
    Method2(i int) error
}
```

* 인터페이스도 이름을 붙일 수 있다.

```go
type Loader interface {
    Load(filename string) error
}
```

* 구조체 내장과 비슷한 형식으로 여러 인터페이스를 합칠 수 있다.

```go
type ReadWriter {
    io.Reader
    io.Writer
}
```

io.Reader와 io.Writer의 모든 메서드를 구현하는 이름 붙인 자료형은 모두 ReadWriter가 된다.



## 5.3.2 커스텀 프린터

* 이름 붙인 자료형은 Print 계열의 함수들을 이용하여 출력할 때 나타나는 형식을 정의할 수 있다.
* fmt 패키지에 있는 Stringer 인터페이스의 String() 메서드를 구현함으로써 출력을 커스텀할 수 있다.
* fmt.Println 함수가 인터페이스를 검사해서 Stringer일 경우 String 메서드를 호출하여 출력하기 때문이다.
* 자료형 변환을 이용하면 다른 구현의 출력을 하게 만들 수 있다.

하위 작업 개념을 만들어 작업안에 여러개의 작업이 있는 구조를 만들고 출력할 수 있다.

```go
// Task is a struct to hold a single task.
type Task struct {
	Title    string    `json:"title,omitempty"`
	Status   status    `json:"status,omitempty"`
	Deadline *Deadline `json:"deadline,omitempty"`
	Priority int       `json:"priority,omitempty"`
	SubTasks []Task    `json:"subTasks,omitempty"`
}

// IncludeSubTasks is a Task but its String method returns the string
// including the sub tasks.
type IncludeSubTasks Task

func (t IncludeSubTasks) indentedString(prefix string) string {
	str := prefix + Task(t).String()
	for _, st := range t.SubTasks {
		str += "\n" + IncludeSubTasks(st).indentedString(prefix+"  ")
	}
	return str
}

// String returns the string representation of t including the sub
// tasks.
func (t IncludeSubTasks) String() string {
	return t.indentedString("")
}

func ExampleIncludeSubTasks_String() {
	fmt.Println(IncludeSubTasks(Task{
		Title:    "Laundry",
		Status:   TODO,
		Deadline: nil,
		Priority: 2,
		SubTasks: []Task{{
			Title:    "Wash",
			Status:   TODO,
			Deadline: nil,
			Priority: 2,
			SubTasks: []Task{
				{"Put", DONE, nil, 2, nil},
				{"Detergent", TODO, nil, 2, nil},
			},
		}, {
			Title:    "Dry",
			Status:   TODO,
			Deadline: nil,
			Priority: 2,
			SubTasks: nil,
		}, {
			Title:    "Fold",
			Status:   TODO,
			Deadline: nil,
			Priority: 2,
			SubTasks: nil,
		}},
	}))
	// Output:
	// [ ] Laundry <nil>
	//   [ ] Wash <nil>
	//     [v] Put <nil>
	//     [ ] Detergent <nil>
	//   [ ] Dry <nil>
	//   [ ] Fold <nil>
}
```



## 5.3.3 정렬과 힙



## 5.3.4 외부 의존성 줄이기

* 외부 리소스를 접근하는 것을 막고 싶은 경우 인터페이스를 사용할 수 있다.

```go
func Save(f *os.File) {
    // ...
}
```

* 위와 같이 작성하면 테스트를 할 때 실제 파일을 넘겨야 한다.
* 인터페이스를 이용하면 실제 구현에서는 실제 파일을 넘기고 테스트에서는 더미 데이터를 넘길 수 있다.

```go
func Save(w io.Writer) {
    // ...
}
```

* 가능하면 만들어져 있는 인터페이스를 받아서 동작하게 코드를 작성하면 유연하게 테스트를 할 수 있다.
* 닫는 연산이 필요하면 io.WriterCloser를 io.Writer 대신에 사용할 수 있다.
* 순차적으로 쓰는 경우가 아니라면 io.WriteSeeker 및 io.WriterAt을 이용할 수 있다.



파일 시스템을 이용하여 파일의 이름을 바꾸고 목록을 살펴보는 코드를 작성할 때 파일 시스템 인터페이스를 만들어서 이용하는 것이 유연성을 높이는데 도움이 된다.

파일 이름 변경과 삭제를 해야 하는 경우라면 다음과 같이 인터페이스를 만들고 내가 사용하는 파일 연산들만 포함하면 된다.

```go
type FileSystem interface {
    func Rename(oldpath, newpath string) error
    func Remove(name string) error
}
```

```go
type OSFileSystem struct {}

func (fs OSFileSystem) Rename(oldpath, newpath string) error {
    return os.Rename(oldpath, newpath)
}

func (fs OSFileSystem) Remove(name string) error {
    return os.Remove(name)
}
```

* 실제 구현 부분에서 만든 인터페이스를 사용한다.

```go
func ManageFile(fs FileSystem) {
    // ...
}
```

* OSFileSystem을 이용하여 호출하면 실제 파일시스템을 이용하고, 테스트 용도로 가짜 파일시스템을 이용할 수 있다.
* 표준 라이브러리를 살펴보고 인터페이스의 용도를 익혀보자.



## 5.3.5 빈 인터페이스와 형 단언

* 빈 인터페이스는 아무 자료형이나 취급할 수 있다.
* interface{} 타입을 원래의 자료형으로 변환하기 위해 형 단언을 사용한다.
* 인터페이스는 실제로 자료형과 값을 갖고 있는 구조체로 표현이된다. 따라서 형변환을 할 때 자료형이 맞는지 실행시간에 검사가 일어나야 한다.

```go
func ExampleCaseInsensitive_heapString() {
    // ...
    popped := heap.Pop(&apple)
    s := popped.(string)
}
```

* .(string)을 이용해 string 타입으로 형 단언을 할 수 잇다.
* 단언한 형이 아니라면 패닉을 발생시키게 된다.
* 인터페이스를 실제 자료형으로 받을 때 마찬가지로 형 단언을 사용한다.

```go
var r io.Reader = NewReader()
f := r.(os.File)
```

* 자료형이 맞지 않으면 패닉이 발생하기 때문에 단언한 형이 맞는지 검사해야한다.
* 형 단언시 두 번째 값으로 단언한 형이 맞는지에 대한 bool 값을 받을 수 있다.

```go
var r io.Reader = NewReader()
f, ok := r.(os.File)
```

* 형 단언은 구체적인 자료형뿐만 아니라 다른 인터페이스로도 가능하다.

```go
var r io.Reader = NewReader()
f, ok := r.(io.ReadCloser)
```



## 5.3.6 인터페이스 변환 스위치

* 인터페이스들이 지정하는 범위는 다양할 수 있다.
* interface{}는 모든 자료형을 포함하는 가장 넓은 인터페이스다.
* 포괄적으로 인터페이스를 받아서 특정 자료형일 때, 혹은 좀 더 좁은 범위의 인터페이스를 구현할 때 구현을 달리하고 싶을 때 자료형 스위치(type switch)를 이용하면 된다.

```go
func Join(sep string, a ...interface{}) string {
	if len(a) == 0 {
		return ""
	}
	t := make([]string, len(a))
	for i := range a {
		switch x := a[i].(type) {
		case string:
			t[i] = x
		case int:
			t[i] = strconv.Itoa(x)
		case fmt.Stringer:
			t[i] = x.String()
		}
		// The switch-case block above is equivalent to the
		// following if-else if.
		//
		// if x, ok := a[i].(string); ok {
		// 	t[i] = x
		// } else if x, ok := a[i].(int); ok {
		// 	t[i] = strconv.Itoa(x)
		// } else if x, ok := a[i].(fmt.Stringer); ok {
		// 	t[i] = x.String()
		// }
	}
	return strings.Join(t, sep)
}
```

* a[i]를 type으로 형 단언 해서 x에 할당한 모양새이다.
* case에서 자료형의 이름을 지정하여 각각의 경우에 대하여 구현을 다르게 해줄 수 있다.
* case 내부에서는 해당 자료형의 값으로 x가 지정된다.
* if를 이용할 수도 있지만 switch문을 이용하는 것이 보기도 좋고 편리하다.