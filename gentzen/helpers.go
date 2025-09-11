package gentzen

// Permutations generates all permutations of a slice using Heap's algorithm.
func permutations[T any](s []T) [][]T {
	var resp [][]T

	var generate func([]T, int)

	generate = func(a []T, n int) {
		if n == 1 {
			e := make([]T, len(a))
			copy(e, a)
			resp = append(resp, e)
			return
		}
		for i := 0; i < n; i++ {
			generate(a, n-1)
			if n%2 == 1 {
				a[0], a[n-1] = a[n-1], a[0]
			} else {
				a[i], a[n-1] = a[n-1], a[i]
			}
		}
	}

	if len(s) == 0 {
		return nil
	}

	generate(s, len(s))
	return resp
}
