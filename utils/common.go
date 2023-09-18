package utils

import (
	"math/rand"

	"golang.org/x/exp/constraints"
)

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

func Max[T constraints.Ordered](arr ...T) T {
	if len(arr) == 0 {
		var zero T
		return zero
	}

	m := arr[0]
	for _, v := range arr {
		if m < v {
			m = v
		}
	}

	return m
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func RandFloat(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
