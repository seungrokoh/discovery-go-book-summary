package problem03

func IsContain(slice []string, target string) int {
	return IterativeBinarySearch(slice, target)
}

func RecursiveBinarySearch(slice []string, target string, start, end int) bool {
	if start > end {
		return false
	}
	middle := (start + end) / 2
	if slice[middle] == target {
		return true
	}

	if slice[middle] > target {
		return RecursiveBinarySearch(slice, target, start, middle - 1)
	}
	return RecursiveBinarySearch(slice, target, middle + 1, end)
}

func IterativeBinarySearch(slice []string, target string) int {
	start := 0
	end := len(slice) - 1

	for start <= end {
		middle := (start + end) / 2

		if slice[middle] == target {
			return middle
		}

		if slice[middle] > target {
			end = middle - 1
		} else if slice[middle] < target {
			start = middle + 1
		}
	}
	return -1
}