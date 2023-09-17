package utils

import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](arr ...T) T {
	if len(arr) == 0 {
		var zero T
		return zero
	}

	m := arr[0]
	for _, v := range arr {
		if m > v {
			m = v
		}
	}

	return m
}
