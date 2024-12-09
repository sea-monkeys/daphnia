package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
	"github.com/sea-monkeys/daphnia"
)



func main() {

	// Create a context
	ctx := context.Background()

	// Initialize the vector store
	vectorStore := daphnia.VectorStore{}
	vectorStore.Initialize("06102023.gob")

	ollamaUrl, errParse := url.Parse("http://host.docker.internal:11434")
	if errParse != nil {
		log.Fatal("😡:", errParse)
	}
	client := api.NewClient(ollamaUrl, http.DefaultClient)

	// Search for the closest chunk(s) to the question
	question := "How to list the running containers?"

	// Embbeding of the question - search for the closest chunk(s)
	reqEmbedding := &api.EmbeddingRequest{
		Model:  "mxbai-embed-large:latest",
		Prompt: question,
	}
	resp, errEmb := client.Embeddings(ctx, reqEmbedding)
	if errEmb != nil {
		fmt.Println("😡:", errEmb)
	}

	// Vector Record from the question
	embeddingFromQuestion := daphnia.VectorRecord{
		Prompt:    reqEmbedding.Prompt,
		Embedding: resp.Embedding,
	}

	similarities, _ := vectorStore.SearchTopNSimilarities(embeddingFromQuestion, 0.5, 2)

	fmt.Println("Question:", question)
	fmt.Println("Similarities:")

	for _, similarity := range similarities {
		fmt.Println()
		fmt.Println("Cosine similarity:", similarity.CosineSimilarity)

		fmt.Println(similarity.Prompt)
	}

}
