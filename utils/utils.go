package utils

import (
	"strconv"
)

func Min(first int, numbers ...int) int {
	min := first
	for _, num := range numbers {
		if num < min {
			min = num
		}
	}

	return min
}

func AtoUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}
