package tools

import (
	"errors"
	"math"
)

// Could possibly use it for normalized vectors
// instead of cosine similarity for vector search
func DotProduct(a, b []float32) (float64, error) {
	if len(a) != len(b) {
		errMsg := "must have same len"
		return 0.0, errors.New(errMsg)
	}
	var dot float64
	for i := range a {
		va := float64(a[i])
		vb := float64(b[i])
		dot += va * vb
	}
	return dot, nil
}

func CosineSim(a, b []float32) (float64, error) {
	if len(a) != len(b) {
		errMsg := "must have same len"
		return 0.0, errors.New(errMsg)
	}
	var dot, na, nb float64
	for i := range len(a) {
		va := float64(a[i])
		vb := float64(b[i])
		dot += va * vb
		na += va * va
		nb += vb * vb
	}
	denom := math.Sqrt(na * nb)
	sim := dot / (denom + 1e-9)
	return sim, nil
}
