package problem01

import "errors"

var (
	un = "은 맛있다"
	nun = "는 맛있다"

	ErrNotHangul = errors.New("it is not hangul")

	hangulStart = rune(44032)
	hangulEnd = rune(55204)
	numEnds = 28
)

func PrintFruitDelicious(fruits []string) []string {
	index := 0
	var result []string

	for _, fruit := range fruits {
		// 한글인지 아닌지 error 체크
		if hasConsonantSuffix, err := HasConsonantSuffix(fruit); err == nil {
			if hasConsonantSuffix {
				// 받침이 있는 경우
				result = append(result, fruit + un)
			} else {
				// 받침이 없는 경우
				result = append(result, fruit + nun)
			}
			index++
		}
	}
	return result
}

// 한글의 받침이 있는지 없는지 확인하는 함수
// 한글이 아닐 경우 ErrNotHangul을 리턴합니다.
func HasConsonantSuffix(s string) (bool, error) {
	hangul := []rune(s)
	lastConsonantSuffix := hangul[len(hangul) - 1]

	if hangulStart > lastConsonantSuffix || lastConsonantSuffix >= hangulEnd {
		// 한글이 아닐 경우
		return false, ErrNotHangul
	}
	// 마지막 글자만을 확인
	index := int(hangul[len(hangul) - 1] - hangulStart)
	return index % numEnds != 0, nil
}