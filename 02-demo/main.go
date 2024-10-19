package main

import (
    "fmt"
    "github.com/sea-monkeys/daphnia"
)

func main() {
    vs := &daphnia.VectorStore{}
    vs.Initialize("meme_database.gob")

    // Store some memes
    vs.Save(daphnia.VectorRecord{
        Prompt:    "Distracted boyfriend meme",
        Embedding: []float64{0.1, 0.8, 0.3}, // Represents "funny", "relatable", "popular"
    })

    vs.Save(daphnia.VectorRecord{
        Prompt:    "This is fine dog meme",
        Embedding: []float64{0.2, 0.9, 0.7}, // Represents "sarcastic", "relatable", "coping"
    })

    // Find a meme for a user feeling sarcastic and trying to cope
    userMood := daphnia.VectorRecord{
        Prompt:    "User feeling sarcastic and trying to cope",
        Embedding: []float64{0.25, 0.85, 0.8},
    }

    results, _ := vs.SearchTopNSimilarities(userMood, 0.9, 1)

    fmt.Println("Recommended meme:", results[0].Prompt)
    // Output: Recommended meme: This is fine dog meme
}