package utils

import (
	"errors"
	"math/rand"
)

// MostFrequent get item with most frequency in array
func MostFrequent[K comparable](arr []K) (count int, value K, err error) {
	if len(arr) == 0 {
		err = errors.New("empty array")
		return
	}
	m := map[K]int{}
	var maxCount int
	var maxKey K
	for _, a := range arr {
		m[a]++
		if m[a] > maxCount {
			maxCount = m[a]
			maxKey = a
		}
	}

	return maxCount, maxKey, nil
}

// GetRandomSubArray Get random array from source array
func GetRandomSubArray[K any](arr []K, size int) []K {
	if size > len(arr) {
		return arr
	}

	indexArr := rand.Perm(len(arr))
	subArray := make([]K, size)
	for i, index := range indexArr[:size] {
		subArray[i] = arr[index]
	}

	return subArray
}
