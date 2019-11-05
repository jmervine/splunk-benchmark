package math

import (
	"sort"
)

func Avg(v []float64) float64 {
	switch l := len(v); l {
	case 0:
		return 0.0
	case 1:
		return v[0]
	default:
		t := float64(0)

		for _, n := range v {
			t = t + n
		}

		return t / float64(len(v))
	}
}

func Med(v []float64) float64 {
	l := len(v)
	switch l {
	case 0:
		return 0.0
	case 1:
		return v[0]
	case 2:
		return Avg(v)
	default:
		if !sort.Float64sAreSorted(v) {
			sort.Float64s(v)
		}

		half := l / 2
		mn := v[half]

		if l%2 == 0 {
			mn = (mn + v[half-1]) / 2
		}

		return mn
	}
}

func Min(v []float64) float64 {
	switch l := len(v); l {
	case 0:
		return 0.0
	case 1:
		return v[0]
	default:
		if !sort.Float64sAreSorted(v) {
			sort.Float64s(v)
		}

		return v[0]
	}
}

func Max(v []float64) float64 {
	switch l := len(v); l {
	case 0:
		return 0.0
	case 1:
		return v[0]
	default:
		if !sort.Float64sAreSorted(v) {
			sort.Float64s(v)
		}

		return v[len(v)-1]
	}
}
