package problem05

func NewMultiSet() map[string]int {
	return make(map[string]int)
}

func Count(m map[string]int, val string) int {
	if val, exist := m[val]; exist {
		return val
	}
	return 0
}

func Insert(m map[string]int, val string) {
	m[val] += 1
}

func Erase(m map[string]int, val string) {
	if m[val] == 0 {
		return
	}
	m[val] = m[val] - 1
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