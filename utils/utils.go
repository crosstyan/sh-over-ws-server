package utils

func DeleteAt[S ~[]E, E any](s S, i int) S {
	if i < 0 || i >= len(s) {
		return s
	}
	return append(s[:i], s[i+1:]...)
}

// https://scribe.rip/how-to-write-generic-helper-functions-with-go-1-18-part-2-19e3d2ab45f5
// Only delete the first `el` in `s`.
// won't complain if `el` is not in `s`.
func DeleteIfOnce[S ~[]E, E any](s S, el E, cmp func(E, E) bool) S {
	for i, v := range s {
		if cmp(v, el) {
			s = DeleteAt(s, i)
			return s
		}
	}
	return s
}

// will delete all `el` in `s`.
func DeleteIfAll[S ~[]E, E any](s S, el E, cmp func(E, E) bool) S {
	var indexes []int
	for i, v := range s {
		if cmp(v, el) {
			indexes = append(indexes, i)
		}
	}
	var res []E
	var last int = 0
	for _, idx := range indexes {
		res = append(res,
			s[last:idx]...)
		last = idx + 1
	}
	// don't miss the last part
	if last < len(s) {
		res = append(res, s[last:]...)
	}
	return res
}
