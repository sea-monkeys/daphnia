package daphnia

import (
	"math"
	"sort"

)

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

func GetTopNVectorRecords(records []VectorRecord, max int) []VectorRecord {
	// Sort the records slice in descending order based on CosineDistance
	sort.Slice(records, func(i, j int) bool {
		return records[i].CosineDistance > records[j].CosineDistance
	})

	// Return the first max records or all if less than three
	if len(records) < max {
		return records
	}
	return records[:max]
}
