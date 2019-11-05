package util

import "testing"

var (
	v0 = []float64{}
	v1 = []float64{5}
	v2 = []float64{1, 2}
	v3 = []float64{1, 3, 2, 4}
	v4 = []float64{1, 5, 6, 7, 2, 3, 4, 8, 9}
	v5 = []float64{1, 5, 6, 7, 8, 2, 3, 4, 9, 10}
)

func TestAvg(t *testing.T) {
	assert(t, "Avg", v3, Avg(v3), 2.5)
}

func TestAvg_whenEmpty(t *testing.T) {
	assert(t, "Avg", v0, Avg(v0), 0.0)
}

func TestAvg_whenSingle(t *testing.T) {
	assert(t, "Avg", v1, Avg(v1), 5.0)
}

func TestMed_whenEven(t *testing.T) {
	assert(t, "Med", v3, Med(v3), 2.5)
	assert(t, "Med", v5, Med(v5), 5.5)
}

func TestMed_whenOdd(t *testing.T) {
	assert(t, "Med", v4, Med(v4), 5.0)
}

func TestMed_whenEmpty(t *testing.T) {
	assert(t, "Med", v0, Med(v0), 0.0)
}

func TestMed_whenOne(t *testing.T) {
	assert(t, "Med", v1, Med(v1), 5.0)
}

func TestMed_whenTwo(t *testing.T) {
	assert(t, "Med", v2, Med(v2), 1.5)
}

func TestMin(t *testing.T) {
	assert(t, "Min", v3, Min(v3), 1.0)
}

func TestMin_whenEmpty(t *testing.T) {
	assert(t, "Min", v0, Min(v0), 0.0)
}

func TestMin_whenSingle(t *testing.T) {
	assert(t, "Min", v1, Min(v1), 5.0)
}

func TestMax(t *testing.T) {
	assert(t, "Max", v3, Max(v3), 4.0)
}

func TestMax_whenEmpty(t *testing.T) {
	assert(t, "Max", v0, Min(v0), 0.0)
}

func TestMax_whenSingle(t *testing.T) {
	assert(t, "Max", v1, Min(v1), 5.0)
}

func assert(t *testing.T, f string, v []float64, a, e float64) {
	if a != e {
		t.Errorf("%s(%#v) = %f; expected %f", f, v, a, e)
	}
}
