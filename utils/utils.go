package utils

import "golang.org/x/exp/slices"

// https://scribe.rip/how-to-write-generic-helper-functions-with-go-1-18-part-2-19e3d2ab45f5
// Only delete the first `el` in `s`.
// won't complain if `el` is not in `s`.
func DeleteIfOnce[S ~[]E, E any](s S, el E, cmp func(E, E) bool) S {
	for i, v := range s {
		if cmp(v, el) {
			s = slices.Delete(s, i, i)
			return s
		}
	}
	return s
}

// will delete all `el` in `s`.
func DeleteIfAll[S ~[]E, E any](s S, el E, cmp func(E, E) bool) S {
	for i, v := range s {
		if cmp(v, el) {
			s = slices.Delete(s, i, i)
		}
	}
	return s
}
