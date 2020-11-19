package resolver

import "math/rand"

// subset 乱序后的数组子集
func subset(set []string, sub int) []string {
	rand.Shuffle(len(set), func(i, j int) {
		set[i], set[j] = set[j], set[i]
	})
	if len(set) <= sub {
		return set
	} else {
		return set[:sub]
	}
}
