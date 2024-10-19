# 🦐 Daphnia: Your Friendly Neighborhood Vector Store

Welcome to Daphnia, the vector database that's small in size but big on personality! Just like its namesake water flea, Daphnia is here to help you filter through your data ocean with ease and efficiency.

## 🌊 What is Daphnia?

Daphnia is an embedded vector database library for Go that allows you to store, retrieve, and search vector embeddings with the speed of a startled water flea. Whether you're building the next big AI application or just trying to find the perfect meme, Daphnia has got your back!

## 🚀 Features

- 🧠 Embedded vector storage
- 🔍 Lightning-fast similarity search
- 🧮 Cosine distance calculations
- 🔐 Persistence with Artemia
- 🤖 Integration with Ollama for embeddings
- 🦸‍♂️ Superhero-level performance (okay, maybe more like water flea-level, but still impressive!)

## 📦 Installation

To get started with Daphnia, simply run:

```bash
go get github.com/sea-monkeys/daphnia
```

## 🎭 Usage Examples

### Example 1: The Terminator Trivia Bot

Imagine you're building a Terminator movie trivia bot. You want to store information about different Terminator models and quickly retrieve relevant facts. Here's how you might use Daphnia:

```go
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
```

### Example 2: The Meme Matchmaker

Let's say you're building a meme recommendation engine (because why not?). You want to find the perfect meme based on a user's mood. Daphnia can help!

```go
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
```

### Example 3: The Ollama-Powered Terminator Knowledge Base

In this example, we'll use Daphnia with Ollama to create a knowledge base about Terminators. We'll embed chunks of text about Terminators and then search for relevant information.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os"
    "strings"

    "github.com/ollama/ollama/api"
    "github.com/sea-monkeys/daphnia"
)

func main() {
    ctx := context.Background()

    // Initialize Daphnia
    vectorStore := daphnia.VectorStore{}
    vectorStore.Initialize("vectors.gob")

    // Set up Ollama client
    ollamaUrl, err := url.Parse("http://host.docker.internal:11434")
    if err != nil {
        log.Fatal("Error parsing Ollama URL:", err)
    }
    client := api.NewClient(ollamaUrl, http.DefaultClient)

    // Read Terminator data
    terminatorsDataFile, err := os.ReadFile("./terminators.data.md")
    if err != nil {
        log.Fatal("Error reading terminators data file:", err)
    }
    terminatorsData := string(terminatorsDataFile)

    // Split the data into chunks
    chunks := strings.Split(terminatorsData, "<!-- SPLIT -->")

    // Embed and store each chunk
    for idx, chunk := range chunks {
        req := &api.EmbeddingRequest{
            Model:  "mxbai-embed-large:latest",
            Prompt: chunk,
        }

        resp, err := client.Embeddings(ctx, req)
        if err != nil {
            fmt.Println("Error getting embedding:", err)
            continue
        }

        _, err = vectorStore.Save(daphnia.VectorRecord{
            Prompt:    chunk,
            Embedding: resp.Embedding,
            Id:        fmt.Sprintf("terminator-%d", idx),
        })

        if err != nil {
            fmt.Println("Error saving vector record:", err)
        }
    }

    // Search for information about T-800 Strengths
    questionEmbed, err := client.Embeddings(ctx, &api.EmbeddingRequest{
        Model:  "mxbai-embed-large:latest",
        Prompt: "What are the T-800 Strengths?",
    })
    if err != nil {
        log.Fatal("Error getting question embedding:", err)
    }

    embeddingFromQuestion := daphnia.VectorRecord{
        Prompt:    "What are the T-800 Strengths?",
        Embedding: questionEmbed.Embedding,
    }

    similarities, err := vectorStore.SearchTopNSimilarities(embeddingFromQuestion, 0.5, 2)
    if err != nil {
        log.Fatal("Error searching for similarities:", err)
    }

    fmt.Println("Top 2 relevant chunks about T-800 Strengths:")
    for _, sim := range similarities {
        fmt.Printf("Similarity: %.4f\n", sim.CosineDistance)
        fmt.Println(sim.Prompt)
        fmt.Println("---")
    }
}
```

In this example:
1. We initialize Daphnia and set up an Ollama client.
2. We read a file containing information about Terminators and split it into chunks.
3. Each chunk is embedded using Ollama and stored in Daphnia.
4. We then embed a question about T-800 strengths and use Daphnia to find the most similar chunks of information.
5. Finally, we print out the top 2 most relevant chunks of information about T-800 strengths.

This demonstrates how Daphnia can be used with Ollama to create a powerful, searchable knowledge base.

## 🧪 How It Works

1. **Vector Storage**: Daphnia uses Artemia's PrevalenceLayer to persistently store your vector records.
2. **Embedding**: When integrated with Ollama, Daphnia can use language models to create embeddings for your text data.
3. **Similarity Search**: When you search for similar vectors, Daphnia calculates the cosine distance between your query and all stored vectors.
4. **Top N Results**: Daphnia can return the top N most similar results, perfect for when you need the cream of the crop.

## 🦸‍♀️ Contributing

We welcome contributions from all water fleas and sea monkeys! If you have ideas for improvements or have found a bug, please open an issue or submit a pull request.

## 📜 License

Daphnia is released under the MIT License. Feel free to use it in your projects, but please don't use it to build Skynet. We have enough problems with terminators as it is!

---

Remember, whether you're building the next big AI application or just trying to find the perfect meme, Daphnia is here to help you swim through your data ocean with the grace of a... well, a microscopic water flea. But hey, they're pretty good swimmers for their size!
