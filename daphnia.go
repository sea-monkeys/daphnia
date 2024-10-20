package daphnia

import (
	"encoding/gob"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/sea-monkeys/artemia"
)

type VectorRecord struct {
	Id             string    `json:"id"`
	Prompt         string    `json:"prompt"` // chunk of text that the vector represents
	Embedding      []float64 `json:"embedding"`
	CosineDistance float64
	Metadata       map[string]interface{} `json:"metadata"` // additional metadata

	// If you want to manage the life time of the vector record
	CreatedAt      time.Time              `json:"created_at"`
	ExpiresAt      time.Time              `json:"expires_at"`
}

type VectorStore struct {
	pl *artemia.PrevalenceLayer
}

// Initialize initializes the VectorStore with the given path.
func (vs *VectorStore) Initialize(path string) error {
	// Register types for gob encoding
	gob.Register(VectorRecord{})

	pl, err := artemia.NewPrevalenceLayer(path)
	if err != nil {
		// If the file doesn't exist, create a new one
		if os.IsNotExist(err) {
			pl, err = artemia.NewPrevalenceLayer(path)
			if err != nil {
				return fmt.Errorf("failed to create new prevalence layer: %w", err)
			}
		} else {
			return fmt.Errorf("failed to open prevalence layer: %w", err)
		}
	}
	vs.pl = pl
	return nil
}

func (vs *VectorStore) Get(id string) (VectorRecord, error) {

	if val, ok := vs.pl.Get(id); ok {
		if vector, ok := val.(VectorRecord); ok {
			return vector, nil
		} else {
			return VectorRecord{}, fmt.Errorf("value is not a VectorRecord")
		}
	}
	return VectorRecord{}, fmt.Errorf("value not found")
}

func (vs *VectorStore) GetAll() ([]VectorRecord, error) {
	var vectorRecords []VectorRecord

	vectors := vs.pl.Query(func(item interface{}) bool {
		_, ok := item.(VectorRecord)
		return ok
	})

	for _, u := range vectors {
		vector := u.(VectorRecord)
		vectorRecords = append(vectorRecords, vector)
	}

	if len(vectorRecords) == 0 {
		return nil, fmt.Errorf("no vector records found")
	}

	return vectorRecords, nil
}

func (vs *VectorStore) Save(vectorRecord VectorRecord) (VectorRecord, error) {
	if vectorRecord.Id == "" {
		vectorRecord.Id = uuid.New().String()
	}
	if err := vs.pl.Set(vectorRecord.Id, vectorRecord); err != nil {
		return VectorRecord{}, err
	}
	return vectorRecord, nil
}

// SearchSimilarities searches for vector records in the VectorStore that have a cosine distance similarity greater than or equal to the given limit.
//
// Parameters:
//   - embeddingFromQuestion: the vector record to compare similarities with.
//   - limit: the minimum cosine distance similarity threshold.
//
// Returns:
//   - []VectorRecord: a slice of vector records that have a cosine distance similarity greater than or equal to the limit.
//   - error: an error if any occurred during the search.
func (vs *VectorStore) SearchSimilarities(embeddingFromQuestion VectorRecord, limit float64) ([]VectorRecord, error) {

	var records []VectorRecord

	allRecords, err := vs.GetAll()
	if err != nil {
		return nil, err
	}

	for _, v := range allRecords {
		distance := CosineDistance(embeddingFromQuestion.Embedding, v.Embedding)
		if distance >= limit {
			v.CosineDistance = distance
			records = append(records, v)
		}
	}
	return records, nil
}

// SearchTopNSimilarities searches for the top N similar vector records based on the given embedding from a question.
// It returns a slice of vector records and an error if any.
// The limit parameter specifies the minimum similarity score for a record to be considered similar.
// The max parameter specifies the maximum number of vector records to return.
func (vs *VectorStore) SearchTopNSimilarities(embeddingFromQuestion VectorRecord, limit float64, max int) ([]VectorRecord, error) {
	records, err := vs.SearchSimilarities(embeddingFromQuestion, limit)
	if err != nil {
		return nil, err
	}
	return getTopNVectorRecords(records, max), nil
}

func getTopNVectorRecords(records []VectorRecord, max int) []VectorRecord {
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
