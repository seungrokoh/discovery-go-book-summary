# Discovery_Go_Summary

# 문자열 및 자료구조

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
