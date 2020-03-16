package problem05

func NewMultiSet() map[string]int {
	return make(map[string]int)
}

func Count(m map[string]int, val string) int {
	return m[val]
}

func Insert(m map[string]int, val string) {
	m[val]++
}

func Erase(m map[string]int, val string) {
	if m[val] <= 1 {
		delete(m, val)
		return
	}
	m[val]--
}

func String(m map[string]int) string {
	result := "{ "
	for key, value := range m {
		for i := 0; i < value; i++ {
			result += key
			result += " "
		}
	}
	result += "}"
	return result
}