package helpers

import "math/rand"

// ShuffleStringSlice randomizes a string slice
func ShuffleStringSlice(s []string) []string {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
	return s
}
