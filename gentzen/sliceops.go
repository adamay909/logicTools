package gentzen

import "slices"

// delete s from the slice set once. Duplicates are kept.
func deleteOnceFunc[T comparable](set []T, v T) func(T) bool {

	deletedOnce := false

	f := func(e T) bool {
		if deletedOnce {
			return false
		}
		if e == v {
			deletedOnce = true
			return true
		}
		return false
	}

	return f
}

// removes the subslice once. (dealing with slices, not sets)
func removeSubslice[T comparable](parentslice []T, subslice []T) []T {

	if !isSubslice(subslice, parentslice) {
		return parentslice
	}

	resp := make([]T, 0, len(parentslice))
	resp = append(resp, parentslice...)

	for i := range subslice {
		resp = slices.DeleteFunc(resp, deleteOnceFunc(resp, subslice[i]))
	}

	return resp
}

func isSubslice[T comparable](par, sub []T) bool {
	var d []T
	d = append(d, par...)
	for j := range sub {
		if !slices.Contains(d, sub[j]) {
			return false
		}
		d = slices.DeleteFunc(d, deleteOnceFunc(d, sub[j]))
	}
	return true
}
