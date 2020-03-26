# 8장 실무패턴

실무에서 맞닥뜨리는 설계 문제들을 풀 수 있는 방법에 대해서 알아보자. 중요한 것은 **어떻게 go로 구현하는지가 아닌 어떤 문제를 풀려고 하는지를 생각하는 것이다.**



## 8.1 오버로딩

오버로딩은 같은 함수 및 메서드를 여러개 두는 것이지만 Go에서는 지원하지 않는다. 

Go에서 **오버로딩을 어떻게 흉내내는지가 중요한 것이 아니라 어떤 문제를 풀기 위해서 오버로딩이 필요한지 생각하는 것이 중요하다.**

* 자료형에 따라 다른 이름 붙이기
    * 오버로딩을 반드시 하지 않아도 되는 경우가 많다. 이 경우에는 자료형에 따라 다른 함수의 이름을 붙이자.
* 동일한 자료형의 자료 개수에 따른 오버로딩
    * 이 경우에는 가변 인자를 사용해서 문제를 해결하자.
* 자료형 스위치 활용하기
    * 오버로딩을 반드시 해야하는 경우에는 **인터페이스로 인자를 받고, 메서드 내에서 자료형 스위치로 다른 자료형에 맞추어 다른 코드가 수행되게** 할 수 있다.
* 다양한 인자 넘기기
    * 서로 다른 이름을 붙여도 상관은 없지만 **기본 값을 포함한 여러 설정을 넘기는 경우에는 이들을 모두 묶은 구조체를 넘기는 것을 고려하자**
* 인터페이스 활용하기
    * 인터페이스를 활용하는 것이 더 나은 경우가 있다.



:heavy_check_mark: 오버로딩을 하기 보다 다른 이름을 붙이는 것이 더 나은 경우

```c++
// volume of a cube
int volume(int s) {
    return s*s*s;
}

// volume of a cylinder
double volume(double r, int h) {
    return 3.14*r*r*static_cast<double>(h);
}

// volume of a cuboid
long volume(long l, int b, int h) {
    return l*b*h;
}
```

* 이 경우 volumeGube, volumeCylinder, volumeCuboid로 이름 짓는 것이 훨씬 낫다.



:heavy_check_mark: 편의상 다양한 인자가 오버로딩 되는 경우

```java
Element getElement(int idx) {
    return getElement(idx, DEFAULT);
}

Element getElement(int idx, Language lang) {
    return getElement(idx, lang, false);
}

Element getElement(int idx, Language lang, bool excludeEmpty) {
    // .. implementation ..
}
```

* 이럴 경우 구조체를 넘기는 것이 더 나을 수 있다.

```go
type Option struct {
    Idx				int
    Lang			Language
    ExcludeEmpty	bool
}

func GetElement(opt Option) *Element {
    // ...
}

func main() {
    GetElement(Option{Idx : 3})
}
```

* 이름이 같고 인자의 자료형이 다를 경우 인터페이스를 활용하는 것이 나을 수 있다.

```go
type Stringer interface {
    String() string
}

type Int int
type Double float64

func (i Int) String() string { ... }
func (d Double) String() string { ... }

func ExampleString() {
    fmt.Println(Int(5).String(), Double(3.7).String())
    // Output:
    // 5 3.7
}
```

* 메서드가 아닌 일반 함수로도 쉽게 변환이 가능하다.

```go
func ToString(s Stringer) string {
    return s.String()
}
```



### 8.1.1 연산자 오버로딩

* Go는 연산자 오버로딩을 지원하지 않는다.
* 연산자 오버로딩은 새로운 문제를 풀기보다는 편의성을 위한 기능이다.
* 이를 인터페이스를 이용하여 쉽게 해결할 수 있다.
* 예를들어 sort.Interface의 Less는 부등호 연산자인 < 를 오버로딩하기 위한 것이다.



## 8.2 템플릿 및 제네릭 프로그래밍

* 제네릭은 자료형을 배제할 수 있는 프로그래밍 패러다임이다.
* 제네릭을 사용하면 자료형에 상관 없이 동일한 알고리즘이라면 하나의 코드로 작성이 가능하다.
* 하지만 Go는 제네릭을 지원하지 않는다.



### 8.2.1 유닛 테스트

다른 언어에서 사용하는 함수인 assertEqual 함수와 같이 두 값이 서로 같은지 확인하는 방법에 대해서 알아본다.



* if를 이용하여 직접 비교하는 방법

```go
if expected != actual {
    t.Error("Not equal!")
}
```

* assertEqual은 한 줄로 간결하게 표현되지만 위 코드는 3줄로 늘어나는 단점이 있다.
* 이에 자료형을 한정시킨 assertStringEqual과 같은 함수를 만들수도 있다.

```go
func assertEqualString(t *testing.T, expected, actual string) {
    if expected != actual {
        t.Errorf("%s != %s", expected, actual)
    }
}

func assertEqualString(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("%d != %d" expected, actual)
    }
}
```

* reflect.DeepEqual을 이용하여 범용적인 assertEqual을 작성하여 비교한다.

```go
func assertEqual(t *testing.T, expected, actual interface{}) {
    if !reflect.DeepEqual(expected, actual) {
        t.Errorf("%v != %v", expected, actual)
    }
}
```

* 테이블 기반 테스트를 진행한다.

```go
func Test(t *testing.T) {
    examples := []struct {
        desc, expected, input string
    } {
        {
            desc: "empty case", 
            expected: "",
            input: "",
        },
        {
            ...
        }, 
    }
    
    for _, ex := range examples {
        actual := f(ex.input)
        if ex.expected != actual {
            t.Errorf("%s: %s != %s", ex.desc, ex.expected, actual)
        }
    }
}
```

* Example 함수를 만들어 테스트를 진행한다.

```go
func Example() {
    fmt.Println(Reverse("abc"))
    // Output:
    // cba
}
```



### 8.2.2 컨테이너 알고리즘

* 컨테이너에 어떤 알고리즘을 적용하고자 할 때, 그 컨테이너가 담고 있는 자료형은 큰 상관이 없이 구현할 수 있어야 하고 그것을 위하여, 혹은 더 효율적으로 구현하기 위하여 제네릭을 사용한다.
* sort, heap 인터페이스에서 대소를 비교하는 부분을 인덱스를 이용한다.
* 두 자료를 직접 비교하기 보단 인텍스를 주고 자료를 비교해 제네릭을 활용하지 않고 원하는 일을 할 수 있다.
* 특정 자료형에 국한되지 않고 출력하거나, 네트워크로 보낼 수 있는 알고리즘은 인터페이스를 활용하면 가능하다.
* 인터페이스를 이용하여 새로운 자료형에 대해서 해당 인터페이스에 맞게 구현해 사용한다.



### 8.2.3 자료형 메타 데이터

* 어떤 자료형이 넘어왔는지에 따라서 다른 코드가 동작하게 하려면 자료형 스위치를 이용하면 된다.
* 어던 자료형인지 뿐만 아니라 해당 자료형에 대한 메타 데이터를 처리하고 싶을 때 reflect 패키지를 이용하면 된다.
* 자료형을 넘겨받아 해당 자료형으로 무언가를 하는 함수를 구현할 수 있다.
* 자료에 대한 자료를 메타데이터라고 한다.

자료형을 넘겨받아 그 값을 자료형을 알아내고 각각의 자료형을 키와 값으로 가지는 맵을 만드는 함수를 구현해보자.

```go
func NewMap(key, value interface{}) interface{} {
    // key와 value의 타입을 알아낸다.
	keyType := reflect.TypeOf(key)
	valueType := reflect.TypeOf(value)
    // key, value 타입의 맵 타입을 만든다.
	mapType := reflect.MapOf(keyType, valueType)
    // map 자료형으로 메타데이터를 만든다.
	mapValue := reflect.MakeMap(mapType)
    // Interface()를 이용해 메타데이터에 대한 값을 가져온다.
	return mapValue.Interface()
}
```

* 이렇게 만든 함수는 interface{}를 반환한다.
* map[string]int로 사용하려면 형 단언을 해서 사용해야 한다.
* 실제 자주 사용하기 보단 주로 라이브러리를 작성할 때 어쩔 수 없이 이용하는게 대부분이다.

```go
m := NewMap("", 0).(map[string]int)
m["abc"] = 3
fmt.Println(m)
```

#### 구조체 필드 확인

* 구조체가 어떤 필드를 가지고 있는지 확인할 수 있다.

구조체의 필드를 확인하고 필드의 이름을 출력하는 예제

```go
func FieldNames(s interface{}) ([]string, error) {
    t := reflect.TypeOf(s)
    if t.Kind() != reflect.Struct {
        return nil, errors.New("FieldNames: s is not a struct")
    }
    var names []string
    n := t.NumField()
    for i := 0; i < n; i++ {
        names = append(names, t.Field(i).Name)
    }
    return names, nil
}

func main() {
	s := struct {
		id int
		Name string
		Age int
	} {}
	fmt.Println(FieldNames(s))
}
// Output:
// [id Name Age] <nil>
```

* reflect.Value를 이용하면 특정 필드의 값을 얻어낼 수도 있다.
* 함수를 받아 다른 함수형으로 변경하여 반환하는 것도 가능하다.

반환값이 없는 함수를 받아 error를 반환하는 함수로 변경하는 예제.

```go
func AppendNilError(f interface{}, err error) (interface{}, error) {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return nil, errors.New("AppendNilError: f is not a function")
	}
	var in []reflect.Type
	var out []reflect.Type

	for i := 0; i < t.NumIn(); i++ {
		in = append(in, t.In(i))
	}
	for i := 0; i < t.NumOut(); i++ {
		out = append(out, t.Out(i))
	}

	out = append(out, reflect.TypeOf((*error)(nil)).Elem())
	funcType := reflect.FuncOf(in, out, t.IsVariadic())
	v := reflect.ValueOf(f)
	funcValue := reflect.MakeFunc(funcType, func(args []reflect.Value) []reflect.Value {
		results := v.Call(args)
		results = append(results, reflect.ValueOf(&err).Elem())
		return results
	})
	return funcValue.Interface(), nil
}

func main() {
	f := func() {
		fmt.Println("called")
	}
	f2, err := AppendNilError(f, errors.New("test error"))
	fmt.Println("AppendNilError.err: ", err)
	fmt.Println(f2.(func() error)())
}

// Output:
// AppendNilError.err:  <nil>
// called
// test error
```

* 동적 자료형 언어에서 할 수 있는 많은 것을 할 수 있게 된다.
* 그러나 reflect를 사용하면 정적인 자료형 검사를 할 수 없으므로 꼭 필요한 경우에만 이용해야한다.



### 8.2.4 go generate

* go generate를 이용하면 임의의 명령을 수행하여 프로그램 코드를 생성할 수 있다.
* 문자열뿐만 아니라 다양한 자료형에 대해 MultiSet을 동작하게 만들고 싶을 때
* 템플릿을 이용해 소스코드를 생성할 수 있다.

```go
// Binary multisetgen provides an example code to generate generic
// code for type T.
package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"text/template"
)

var (
	packageName = flag.String(
		"package_name",
		"main",
		"package name",
	)
	multisetTypename = flag.String(
		"multiset_typename",
		"MultiSet",
		"container type",
	)
	elementTypename = flag.String(
		"element_typename",
		"string",
		"element type",
	)
	output = flag.String(
		"output",
		"",
		"output filename",
	)
)

// tmpl is a code template for multi set implementation for type T.
var tmpl = template.Must(template.New("multiset").Parse(`// Generated by multisetgen. DO NOT EDIT!
package {{.PackageName}}

import "fmt"

type {{.MultisetTypename}} map[{{.ElementTypename}}]int

func New{{.MultisetTypename}}() {{.MultisetTypename}} {
	return {{.MultisetTypename}}{}
}

func (m {{.MultisetTypename}}) Insert(val {{.ElementTypename}}) {
	m[val]++
}

func (m {{.MultisetTypename}}) Erase(val {{.ElementTypename}}) {
	if _, exists := m[val]; !exists {
		return
	}
	m[val]--
	if m[val] <= 0 {
		delete(m, val)
	}
}

func (m {{.MultisetTypename}}) Count(val {{.ElementTypename}}) int {
	return m[val]
}

func (m {{.MultisetTypename}}) String() string {
	vals := ""
	for val, count := range m {
		for i := 0; i < count; i++ {
			vals += fmt.Sprint(val) + " "
		}
	}
	return "{ " + vals + "}"
}
`))

// outputFilename returns a output or returns lower cased
// multisetTypename.go.
func outputFilename(output, multisetTypename string) string {
	if output != "" {
		return output
	}
	return strings.ToLower(multisetTypename + ".go")
}

func main() {
	flag.Parse()
	out, err := os.Create(outputFilename(*output, *multisetTypename))
	if err != nil {
		log.Println(err)
		return
	}
	if err := tmpl.Execute(out, struct {
		PackageName      string
		MultisetTypename string
		ElementTypename  string
	}{*packageName, *multisetTypename, *elementTypename}); err != nil {
		log.Println(err)
		return
	}
	log.Println("File written:", out.Name())
}
```

* tmpl은 MultiSet의 소스 코드를 템플릿으로 만든 것
* 각각의 패키지 이름, MultiSet의 자료형 이름, 각 원소의 자료형 이름 부분을 템플릿으로 대체해서 템플릿이 수행될 때 알맞은 값이 들어가도록 한다.
* output을 지정하지 않으면 MultiSet 자료형 이름을 소문자로 변환한 뒤 .go 확장자를 붙인다.



## 8.3 객체지향

Go는 객체지향을 완전히 지원하지는 않지만 객체지향과 같이 다루어지는 경우가 많다.



### 8.3.1 다형성

객체지향의 꽃이라 불리는 다형성으로 객체가 메서드에 대한 다양한 구현을 가지고 있을 수 있다. 

* 다형성은 메서드가 호출되었을 때, 어떤 자료형이냐에 따라서 다른 구현을 할 수 있도록 하는 문제를 풀기 위함이다.
* Go에서는 인터페이스 쉽게 구현이 가능하며 다른 언어보다 더 깔끔하게 구현이 가능하다.

```go
type Shape interface {
	Area() float32
}

type Square struct {
	Size float32
}

func (s Square) Area() float32 {
	return s.Size * s.Size
}

type Rectangle struct {
	Width, Height float32
}

func (r Rectangle) Area() float32 {
	return r.Width * r.Height
}

type Triangle struct {
	Width, Height float32
}

func (t Triangle) Area() float32 {
	return 0.5 * t.Width * t.Height
}

func TotalArea(shapes []Shape) float32 {
	var total float32
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}

func main() {
    fmt.Println(TotalArea([]Shape{
        Square{3}, 
        Rectangle{4, 5}, 
        Triangle{6, 7}, 
    }))
}
// Output:
// 50
```

* TotalArea 함수에 넘기는 슬라이스는 Area만 구현하고 있으면 어떤 도형들도 담아서 넘겨줄 수 있다.
* 다형성을 주로 이용하는 [커맨드 패턴](https://gmlwjd9405.github.io/2018/07/07/command-pattern.html)과 같은 것들도 Go에서는 쉽게 구현할 수 있다.



### 8.3.2 인터페이스

* 다른 언어의 인터페이스와 다른 점은 Go는 인터페이스 내의 메서드들을 구현하기만 하면 그 인터페이스를 구현하는 것이 된다. **매우 중요한 특성이다.**
* 그 이유는 다른 패키지에 있는 구조체를 확장하여 다형성을 구현하려 할 때 해당 구조체가 만든 인터페이스를 구현하기만 하면 확장하고 재사용할 수 있다.

두 번째 의미에 대해서 자세하게 살펴보자. 위 코드를 아직 만들지 않은 상태에서 만약 다른 패키지에 삼각형에 대한 구조체가 정의되어 있다고 가정하자. 아래와 같은 코드로 Triangle이 구현되어 있다.

```go
type Triangle struct {
    Width, Height float32
}

func (t Triangle) Area() float32 {
    return 0.5 * t.Width * t.Height
}
```

만약 여기서 Triangle 뿐만 아니라 여러 도형을 확장시키고 싶을 때, Triangle 구조체가 가지고 있는 Area() 메서드를 가진 인터페이스를 만들고 추가로 구현하고 싶은 여러 도형을 구현하면 된다.

```go
type Shape interface {
    Area() float32
}

type Square struct {
    Size float32
}

func (s Square) Area() float32 {
    return s.Size * s.Size
}

type Rectangle struct {
    Width, Height float32
}

func (r Rectangle) Area() float32 {
    return r.Width * r.Height
}
```

위와 같이 추가로 Square, Rectangle을 만들고 모두 인터페이스를 구현하게 하면 이 두가지는 Shape 인터페이스가 된다. 또한 다른 패키지의 Triangle 구조체 또한 Shape 인터페이스가 된다.



* 다른 언어와 달리 암시적(implict)으로 인터페이스를 구현함으로써 Triangle의 코드를 고치거나 상속받은 클래스를 만들어 다형성을 지원하게 하지 않아도 된다.
* 이 점은 유연하게 프로그램을 확장하고 재사용할 수 있는 큰 장점이다.



### 8.3.3 상속

* 상속은 어떤 클래스의 구현들을 재사용하기 위하여 사용된다.
* Is-A 관계와 Has-A 관계가 성립한다.
* Has-A 관계의 경우 전통적인 객체지향에서도 상속보다는 객체 구성(object composition)이다.
* Go에서는 재사용하고자 하는 구현의 자료형 변수를 struct에 내장하면 된다. Has-A 관계는 말 그대로 구현을 필드로 가지고 있으면 된다.
* Is-A 관계의 상속에서는 많은 경우에 추상 클래스를 상속한다.
* Go에서는 인터페이스를 이용하면 Is-A관계를 나타낼 수 있따.



:bulb: 추상 클래스를 사용하는 이유??

> 구현에 구애받지 않고 이용하고 필요한 경우에 서로 다른 구현으로 다형성을 활용하기 위함



* 추상클래스가 아닌 클래스를 상속받는 경우는 인터페이스와 구현을 함께 상속한다.
* 인터페이스를 상속함으로써 상위 클래스의 일종이라는 것을 이용한 다형성 코드가 가능해지고, 구현을 상속함으로써 반복을 피할 수 있다. 이는 인터페이스와 구조체 내장을 동시에 이용하면 가능하다.



#### 메서드 추가

* 기존에 있던 코드를 재사용하면서 기능을 추가하고 싶을 경우에 상속할 수 있다. 이때 메서드를 추가하여 기능을 추가한다.
* Go에서는 구조체 내장을 이용해 문제를 해결할 수 있다.

다형성 예제에서 도형에 Area() 외에 둘레를 구하는 기능을 추가하고자 할 때 **구조체 내장을 이용해서 기존 코드를 재사용하고 새로운 메서드를 만들면 된다.**

```go
type Rectangle struct {
    Width, Height float32
}

func (r Rectangle) Area() float32 {
    return r.Width * r.Height
}
```

위 코드를 재사용하면서 기능을 추가할 때 구조체 내장을 사용한 뒤 원하는 기능을 메서드로 추가하면 된다.

```go
type RectangleCircum struct {
    Rectangle
}

func (r RectangleCircum) Circum() float32 {
    return 2 * (r.Width + r.Height)
}

func main() {
    r := RetangleCircum{Rectangle{3, 4}}
    fmt.Println(r.Area())
    fmt.Println(r.Circum())
}
// Output: 
// 12
// 14
```

* 필요하다면 상속과 함께 생성자도 만들어줄 수 있다.

```go
func NewRectangleCircum(width, height float32) *RectangleCircum {
    return &RectangleCircum(Rectangle(width, height))
}
```

* 여러 단계로 내장이 이루어지는 경우에도 다시 RectangleCircum을 내장하여 이 구현에 다른 구현을 추가할 수 있다.

#### 오버라이딩

* 기존에 있던 구현을 다른 구현으로 대체하고자 하는 경우에도 상속을 쓸 수 있다. 
* 구조체 내장을 이용해 오버라이딩을 해결할 수 있다.

```go
type WrongRectangle struct { Rectangle }

func (r WrongRectangle) Area() float32 {
    return r.Rectangle.Area() * 2
}

func main() {
    r := WrongRectangle{Rectangle{3, 4}}
    fmt.Println(r.Area())
}
// Output:
// 24
```

* WrongRectangle이 Area()를 구현했기 때문에 바로 Rectangle의 Area()가 호출되지 않는다.

#### 서브타입

* 기존 객체가 쓰이던 곳에 상속받은 객체를 쓰고자 상속하기도 한다.
* Go에서는 인터페이스와 구조체 내장을 모두 사용하면 해결할 수 있다.



:bulb: 위에서 구현한 WrongRectangle도 Shape로 취급이 될까?

> Area 메서드를 구현하고 있기 때문에 당연히 Shape로 취급된다.

:bulb: RectangleCircum은 Shape로 취급이 될까?

> Area()를 구현하고 있지 않지만 내장하고 있는 구조체가 Area()를 구현하고 있기 때문에 Shape로 취급된다.



* reflect.Type.Implements 메서드를 이용하면 인터페이스를 구현하고 있는지 여부를 알 수 있다.
* 내장된 구조체가 있는지 확인하기 위해 내장된 구조체의 이름으로 필드를 찾은 다음 Anonymous 필드를 찾아볼 수 있다.

```go
impl := reflect.TypeOf(RectangleCircum()).Implements(reflect.TypeOf((*Shape)(nil)))

field, ok := reflect.TypeOf(RectangleCircum{}).FieldByName("Rectangle")
emb := ok && field.Anonymous && field.Type == reflect.TypeOf(Rectangle{})
```

* impl에는 RectangleCircum이 Shape 인터페이스를 구현하는지 여부가 기록된다.
* emb에는 Rectangle이 RectangleCircum에 내장되어 있는지 여부가 기록된다.

:bulb: 인터페이스 타입을 얻어낼 때의 주의점

> Shape(nil)은 nil 인터페이스이기 때문에 인터페이스의 포인터로 nil을 만들고 자료형을 얻어내면 인터페이스 포인터 자료형이 되고 여기에 Elem()을 호출하면 인터페이스 자료형을 얻어 낼 수 있다.

 

* RectangleCircum{}과 같이 빈객체를 만드는것이 싫으면 아래와 같이 작성할 수도 있다.

```go
impl := reflect.TypeOf((*RectangleCircum)(nil)).Elem().Implements(
    reflect.TypeOf((*Shape)(nil)).Elem()
)
```

* reflect로 서브 타입 검사를 할 수 있는 것은 좋지만 자주 쓰게 되지는 않는다.



### 8.3.4 캡슐화

* 객체 안에 있는 정보를 바깥에 숨기고자 하는 것이 캡슐화이다.
* 소문자로 이름을 시작하게 만들어 다른 패키지에서 참조가 불가능하게 만들어 캡슐화를 할 수 있다.
* 내가 만든 다른 패키지에서 접근이 가능하게 하고 남은 접근하지 못하게 하고 싶은 경우 패키지 경로에 internal을 넣으면 internal이 있는 경로에 있는 패키지를 포함한 범위에서만 참조가 가능하다.
* getter, setter를 만들어 외부에서도 접근할 수 있게 만들 수 있다.
* Go관례상 getter는 get을 prefix를 제거하고 setter는 prefix로 Set을 써서 이름을 짓는다. (ex. Len(), SetLen(x))
* 생성자 또한 NewX()의 형태로 다른 패키지에서 호출하게 할 수 있다.

:bulb: **NewX()의 반환값으로는 해당 자료형의 오브젝트를 그대로 돌려줄 수도 있지만, 그 자료형이 구현하는 인터페이스를 돌려줌으로써 마치 추상자료형을 돌려주는 것 같은 형태를 만들어 유연한 프로그래밍이 가능하도록 한다.**



## 8.4 디자인 패턴

몇 가지 유명한 디자인 패턴들이 풀고자 하는 문제를 Go에서 어떻게 풀 수 있는지 확인해보자.



### 8.4.1 반복자 패턴

* 클로져를 이용하여 호출하는 반복자, 콜백을 넘겨주어 함수가 모든 원소에 대하여 호출되게 하는 반복자, 인터페이스를 이용한 반복자, 채널을 이용한 반복자 모두가 해당된다.
* 채널을 이용한 반복자를 이용할 때, **중간에 중단할 수 있어야 하는 경우 반드시 done 채널 또는 context.Context를 받아서 처리할 수 있도록 작성해야 한다.**



### 8.4.2 추상 팩토리 패턴

* 여럿 묶어 놓은 팩토리를 추상화하는 패턴이다.
* [추상 팩토리란](https://gmlwjd9405.github.io/2018/08/08/abstract-factory-pattern.html) 구체적인 클래스에 의존하지 않고 서로 연관되거나 의존적인 객체들의 조합을 만드는 인터페이스를 제공하는 패턴
* 어떤 팩토리를 받았는지 상관없이 동일한 코드로 팩토리에서 인스턴스들을 찍어낼 수 있다.

```go
type Button interface {
	Paint()
	OnClick()
}

type Label interface {
	Paint()
}

// WinButton is a Button implementation for Windows.
type WinButton struct{}

func (WinButton) Paint()   { fmt.Println("win button paint") }
func (WinButton) OnClick() { fmt.Println("win button click") }

// WinLabel is a Label implementation for Windows.
type WinLabel struct{}

func (WinLabel) Paint() { fmt.Println("win label paint") }

// WinButton is a Button implementation for Mac.
type MacButton struct{}

func (MacButton) Paint()   { fmt.Println("mac button paint") }
func (MacButton) OnClick() { fmt.Println("mac button click") }

// WinLabel is a Label implementation for Mac.
type MacLabel struct{}

func (MacLabel) Paint() { fmt.Println("mac label paint") }

// UI factory can create buttons and labels.
type UIFactory interface {
	CreateButton() Button
	CreateLabel() Label
}

// WinFactory is a UI factory that can create Windows UI elements.
type WinFactory struct{}

func (WinFactory) CreateButton() Button {
	return WinButton{}
}

func (WinFactory) CreateLabel() Label {
	return WinLabel{}
}

// MacFactory is a UI factory that can create Mac UI elements.
type MacFactory struct{}

func (MacFactory) CreateButton() Button {
	return MacButton{}
}

func (MacFactory) CreateLabel() Label {
	return MacLabel{}
}

// CreateFactory returns a UIFactory of the given os.
func CreateFactory(os string) UIFactory {
	if os == "win" {
		return WinFactory{}
	} else {
		return MacFactory{}
	}
}

func Run(f UIFactory) {
	button := f.CreateButton()
	button.Paint()
	button.OnClick()
	label := f.CreateLabel()
	label.Paint()
}

func main() {
    f1 := CreateFactory("win")
    Run(f1)
    f2 := CreateFactory("mac")
    Run(f2)
}
```



### 8.4.3 비지터 패턴

* 비지터 패턴은 알고리즘을 객체 구조에서 분리시키기 위한 디자인 패턴이다.
* 구조를 수정하지 않고도 새로운 동작을 추가할 수 있다.
* 비지터에서 수행하는 기능과 관련된 코드를 한 곳에 집중시켜 놓을 수 있다.



```go
type CarElementVisitor interface {
	VisitWheel(wheel Wheel)
	VisitEngine(engine Engine)
	VisitBody(body Body)
	VisitCar(car Car)
}

type Acceptor interface {
	Accept(visitor CarElementVisitor)
}

type Wheel string

func (w Wheel) Name() string {
	return string(w)
}

func (w Wheel) Accept(visitor CarElementVisitor) {
	visitor.VisitWheel(w)
}

type Engine struct{}

func (e Engine) Accept(visitor CarElementVisitor) {
	visitor.VisitEngine(e)
}

type Body struct{}

func (b Body) Accept(visitor CarElementVisitor) {
	visitor.VisitBody(b)
}

type Car []Acceptor

func (c Car) Accept(visitor CarElementVisitor) {
	for _, e := range c {
		e.Accept(visitor)
	}
	visitor.VisitCar(c)
}

type CarElementPrintVisitor struct{}

func (CarElementPrintVisitor) VisitWheel(wheel Wheel) {
	fmt.Println("Visiting " + wheel.Name() + " wheel.")
}

func (CarElementPrintVisitor) VisitEngine(engine Engine) {
	fmt.Println("Visiting engine")
}

func (CarElementPrintVisitor) VisitBody(body Body) {
	fmt.Println("Visiting body")
}

func (CarElementPrintVisitor) VisitCar(car Car) {
	fmt.Println("Visiting car")
}

type CarElementDoVisitor struct{}

func (CarElementDoVisitor) VisitWheel(wheel Wheel) {
	fmt.Println("Kicking my " + wheel.Name() + " wheel.")
}

func (CarElementDoVisitor) VisitEngine(engine Engine) {
	fmt.Println("Starting my engine")
}

func (CarElementDoVisitor) VisitBody(body Body) {
	fmt.Println("Moving my body")
}

func (CarElementDoVisitor) VisitCar(car Car) {
	fmt.Println("Starting my car")
}

func main() {
	car := Car{
		Wheel("front left"),
		Wheel("front right"),
		Wheel("back left"),
		Wheel("back right"),
		Body{},
		Engine{},
	}
	car.Accept(CarElementPrintVisitor{})
	car.Accept(CarElementDoVisitor{})
}
// Output:
// Visiting front left wheel.
// Visiting front right wheel.
// Visiting back left wheel.
// Visiting back right wheel.
// Visiting body
// Visiting engine
// Visiting car
// Kicking my front left wheel.
// Kicking my front right wheel.
// Kicking my back left wheel.
// Kicking my back right wheel.
// Moving my body
// Starting my engine
// Starting my car
```