package utils

import "math/rand"

func MostFrequent[K comparable](arr []K) (int, K) {
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

	return maxCount, maxKey
}

// GetRandomSubArray Get random array from source array
// using Fisher–Yates shuffle algorithm (https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle)
func GetRandomSubArray[K any](arr []K, size int) []K {
	if size > len(arr) {
		return arr
	}

	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})

	return arr[:size]
}
