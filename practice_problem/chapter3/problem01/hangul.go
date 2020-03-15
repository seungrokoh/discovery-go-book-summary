package problem01

var (
	un = "은 맛있다"
	nun = "는 맛있다"
)

func PrintFruitDelicious(fruits []string) []string {
	result := make([]string, len(fruits))

	for i, fruit := range fruits {
		if HasConsonantSuffix(fruit) {
			// 받침이 있는 경우
			result[i] = fruit + un
		} else {
			// 받침이 없는 경우
			result[i] = fruit + nun
		}
	}
	return result
}

var (
	start = rune(44032)
	end = rune(55204)
)

// 한글의 받침이 있는지 없는지 확인하는 함수
func HasConsonantSuffix(s string) bool {
	numEnds := 28
	hangul := []rune(s)
	// 마지막 글자만을 확인
	index := int(hangul[len(hangul) - 1] - start)
	return index % numEnds != 0
}