package main

import (
    "fmt"
    "github.com/sea-monkeys/daphnia"
)

func main() {
    // Initialize Daphnia
    vs := &daphnia.VectorStore{}
    vs.Initialize("terminator_facts.gob")

    // Store some Terminator facts
    vs.Save(daphnia.VectorRecord{
        Prompt:    "The T-800 is a cybernetic organism with living tissue over a metal endoskeleton.",
        Embedding: []float64{0.1, 0.2, 0.3}, // Pretend this is a real embedding
    })

    vs.Save(daphnia.VectorRecord{
        Prompt:    "The T-1000 is made of mimetic polyalloy, allowing it to shape-shift.",
        Embedding: []float64{0.4, 0.5, 0.6}, // Pretend this is a real embedding
    })

    // Query for similar facts
    question := daphnia.VectorRecord{
        Prompt:    "What is the T-800 made of?",
        Embedding: []float64{0.15, 0.25, 0.35}, // Pretend this is a real embedding
    }

    results, _ := vs.SearchTopNSimilarities(question, 0.8, 1)

    fmt.Println("Bot says:", results[0].Prompt)
    // Output: Bot says: The T-800 is a cybernetic organism with living tissue over a metal endoskeleton.
}