package math

import "sort"

func Avg(v []float64) float64 {
	t := float64(0)

	for _, n := range v {
		t = t + n
	}

	return t / float64(len(v))
}

func Med(v []float64) float64 {
	if !sort.Float64sAreSorted(v) {
		sort.Float64s(v)
	}

	nn := len(v) / 2
	if nn%2 != 0 {
		return v[nn]
	}

	return (v[nn+1] + v[nn]) / 2
}

func Min(v []float64) float64 {
	if !sort.Float64sAreSorted(v) {
		sort.Float64s(v)
	}

	return v[0]
}

func Max(v []float64) float64 {
	if !sort.Float64sAreSorted(v) {
		sort.Float64s(v)
	}

	return v[len(v)-1]
}
