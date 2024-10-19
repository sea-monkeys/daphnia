package daphnia

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestVectorStore(t *testing.T) {
	// Create a temporary directory for testing
	tmpdir, err := os.MkdirTemp("", "vectorstore_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	// Use a file path in the temporary directory
	testFilePath := tmpdir + "/test_vectors.gob"

	vs := &VectorStore{}
	err = vs.Initialize(testFilePath)
	if err != nil {
		t.Fatalf("Failed to initialize VectorStore: %v", err)
	}

	t.Run("Initialize", func(t *testing.T) {
		// Test initializing with an existing file
		err := vs.Initialize(testFilePath)
		if err != nil {
			t.Errorf("Failed to initialize with existing file: %v", err)
		}

		// Test initializing with a new file
		newFilePath := tmpdir + "/new_test_vectors.gob"
		err = vs.Initialize(newFilePath)
		if err != nil {
			t.Errorf("Failed to initialize with new file: %v", err)
		}
	})

	t.Run("Save", func(t *testing.T) {
		record := VectorRecord{
			Id:        "test1",
			Prompt:    "Test prompt",
			Embedding: []float64{0.1, 0.2, 0.3},
			Metadata:  map[string]interface{}{"key": "value"},
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(time.Hour),
		}

		savedRecord, err := vs.Save(record)
		if err != nil {
			t.Fatalf("Failed to save record: %v", err)
		}
		if savedRecord.Id != record.Id {
			t.Errorf("Expected saved record ID %s, got %s", record.Id, savedRecord.Id)
		}
	})

	t.Run("Get", func(t *testing.T) {
		record := VectorRecord{
			Id:        "test2",
			Prompt:    "Another test prompt",
			Embedding: []float64{0.4, 0.5, 0.6},
		}

		_, err := vs.Save(record)
		if err != nil {
			t.Fatalf("Failed to save record for Get test: %v", err)
		}

		retrievedRecord, err := vs.Get("test2")
		if err != nil {
			t.Fatalf("Failed to get record: %v", err)
		}
		if retrievedRecord.Prompt != record.Prompt {
			t.Errorf("Expected prompt %s, got %s", record.Prompt, retrievedRecord.Prompt)
		}
		if !reflect.DeepEqual(retrievedRecord.Embedding, record.Embedding) {
			t.Errorf("Expected embedding %v, got %v", record.Embedding, retrievedRecord.Embedding)
		}
	})

	t.Run("Get Non-existent Record", func(t *testing.T) {
		_, err := vs.Get("non-existent")
		if err == nil {
			t.Error("Expected error when getting non-existent record, got nil")
		}
		if err != nil && err.Error() != "value not found" {
			t.Errorf("Expected error 'value not found', got '%v'", err)
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		records := []VectorRecord{
			{Id: "test3", Prompt: "Test 3", Embedding: []float64{0.7, 0.8, 0.9}},
			{Id: "test4", Prompt: "Test 4", Embedding: []float64{1.0, 1.1, 1.2}},
		}

		for _, record := range records {
			_, err := vs.Save(record)
			if err != nil {
				t.Fatalf("Failed to save record: %v", err)
			}
		}

		allRecords, err := vs.GetAll()
		if err != nil {
			t.Fatalf("Failed to get all records: %v", err)
		}
		if len(allRecords) < 4 {
			t.Errorf("Expected at least 4 records, got %d", len(allRecords))
		}
	})

	t.Run("SearchSimilarities", func(t *testing.T) {
		queryRecord := VectorRecord{
			Prompt:    "Query",
			Embedding: []float64{0.15, 0.25, 0.35},
		}

		results, err := vs.SearchSimilarities(queryRecord, 0.5)
		if err != nil {
			t.Fatalf("Failed to search similarities: %v", err)
		}
		if len(results) == 0 {
			t.Error("Expected non-empty results, got empty")
		}

		for _, result := range results {
			if result.CosineDistance < 0.5 {
				t.Errorf("Expected cosine distance >= 0.5, got %f", result.CosineDistance)
			}
		}
	})

	t.Run("SearchTopNSimilarities", func(t *testing.T) {
		queryRecord := VectorRecord{
			Prompt:    "Query",
			Embedding: []float64{0.15, 0.25, 0.35},
		}

		results, err := vs.SearchTopNSimilarities(queryRecord, 0.5, 2)
		if err != nil {
			t.Fatalf("Failed to search top N similarities: %v", err)
		}
		if len(results) > 2 {
			t.Errorf("Expected at most 2 results, got %d", len(results))
		}

		for _, result := range results {
			if result.CosineDistance < 0.5 {
				t.Errorf("Expected cosine distance >= 0.5, got %f", result.CosineDistance)
			}
		}

		if len(results) > 1 && results[0].CosineDistance < results[1].CosineDistance {
			t.Errorf("Expected results to be sorted in descending order of cosine distance")
		}
	})
}

func TestCosineDistance(t *testing.T) {
	testCases := []struct {
		name     string
		v1       []float64
		v2       []float64
		expected float64
	}{
		{
			name:     "Identical vectors",
			v1:       []float64{1, 0, 0},
			v2:       []float64{1, 0, 0},
			expected: 1,
		},
		{
			name:     "Orthogonal vectors",
			v1:       []float64{1, 0, 0},
			v2:       []float64{0, 1, 0},
			expected: 0,
		},
		{
			name:     "Opposite vectors",
			v1:       []float64{1, 0, 0},
			v2:       []float64{-1, 0, 0},
			expected: -1,
		},
		{
			name:     "Similar vectors",
			v1:       []float64{1, 1, 1},
			v2:       []float64{2, 2, 2},
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CosineDistance(tc.v1, tc.v2)
			if !almostEqual(result, tc.expected, 1e-6) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

// Helper function to compare float64 values with a tolerance
func almostEqual(a, b, tolerance float64) bool {
	return abs(a-b) <= tolerance
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}