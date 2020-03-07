package utils

func Min(first int, numbers ...int) int {
	min := first
	for _, num := range numbers {
		if num < min {
			min = num
		}
	}

	return min
}
