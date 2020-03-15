package problem03

func IsContain(slice []string, target string) bool {
	return BinarySearch(slice, target, 0, len(slice) - 1)
}

func BinarySearch(slice []string, target string, start, end int) bool {
	if start > end {
		return false
	}
	middle := (start + end) / 2
	if slice[middle] == target {
		return true
	}

	if slice[middle] > target {
		return BinarySearch(slice, target, start, middle - 1)
	}
	return BinarySearch(slice, target, middle + 1, end)
}