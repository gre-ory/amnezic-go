package util

import "math/rand"

// //////////////////////////////////////////////////
// convert

func Shuffle[T any](items []T) {
	nb := len(items)
	for i := 0; i < nb; i++ {
		index := rand.Intn(nb)
		if index != i {
			items[i], items[index] = items[index], items[i]
		}
	}
}
