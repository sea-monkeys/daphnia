package daphnia

import (
	"math"
)

/*
Cosine Similarity:

- Measures the cosine of the angle between two vectors, regardless of their magnitude.
- Values range from -1 to 1 (1 indicating identical vectors, 0 indicating orthogonal vectors, -1 indicating opposite vectors).
- Particularly useful for comparing documents or text embeddings.

CosineDistance calculates the dot product and magnitudes of the vectors to determine their similarity.
*/

func dotProduct(v1 []float64, v2 []float64) float64 {
	// Calculate the dot product of two vectors
	sum := 0.0
	for i := range v1 {
		sum += v1[i] * v2[i]
	}
	return sum
}

// CosineDistance calculates the cosine distance between two vectors
func CosineDistance(v1, v2 []float64) float64 {
	// Calculate the cosine distance between two vectors
	product := dotProduct(v1, v2)

	norm1 := math.Sqrt(dotProduct(v1, v1))
	norm2 := math.Sqrt(dotProduct(v2, v2))
	if norm1 <= 0.0 || norm2 <= 0.0 {
		// Handle potential division by zero
		return 0.0
	}
	return product / (norm1 * norm2)
}
