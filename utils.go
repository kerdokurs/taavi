package main

func ForEach[T any](ts []T, fn func(t T)) {
	for _, t := range ts {
		fn(t)
	}
}

func Map[T, U any](ts []T, mapper func(t T) U) []U {
	us := make([]U, len(ts))

	for i, t := range ts {
		us[i] = mapper(t)
	}

	return us
}

func SliceEq[T comparable](as, bs []T) bool {
	if len(as) != len(bs) {
		return false
	}

	for i, a := range as {
		if bs[i] != a {
			return false
		}
	}

	return true
}
