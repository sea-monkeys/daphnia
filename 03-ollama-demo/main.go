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

	vectorStore := daphnia.VectorStore{}
	vectorStore.Initialize("vectors.gob")

	ollamaUrl, errParse := url.Parse("http://host.docker.internal:11434")
	if errParse != nil {
		log.Fatal("😡:", errParse)
	}
	client := api.NewClient(ollamaUrl, http.DefaultClient)

	// 🖐️🤖 all data about the terminators
	terminatorsDataFile, err := os.ReadFile("./terminators.data.md")
	if err != nil {
		log.Fatal("Error reading terminators data file")
	}
	terminatorsData := string(terminatorsDataFile)

	// Split the data
	chunks := strings.Split(terminatorsData, "<!-- SPLIT -->")

	for idx, chunk := range chunks {

		req := &api.EmbeddingRequest{
			Model:  "mxbai-embed-large:latest",
			Prompt: chunk,
		}

		resp, errEmb := client.Embeddings(ctx, req)
		if errEmb != nil {
			fmt.Println("😡:", errEmb)
		}

		_, err := vectorStore.Save(daphnia.VectorRecord{
			Prompt:    chunk,
			Embedding: resp.Embedding,
			Id: 	  fmt.Sprintf("terminator-%d", idx),
		})

		if err != nil {
			fmt.Println("😡:", err)
		}

	}

	// Embbeding of the question - search for the closest chunk(s)
	reqEmbedding := &api.EmbeddingRequest{
		Model:  "mxbai-embed-large:latest",
		Prompt: "What are the T-800 Strengths?",
	}
	resp, errEmb := client.Embeddings(ctx, reqEmbedding)
	if errEmb != nil {
		fmt.Println("😡:", errEmb)
	}

	embeddingFromQuestion := daphnia.VectorRecord{
		Prompt:    "What are the T-800 Strengths?",
		Embedding: resp.Embedding,
	}
	similarities, _ := vectorStore.SearchTopNSimilarities(embeddingFromQuestion, 0.5, 2)

	for _, sim := range similarities {
		fmt.Println(sim.Prompt, sim.CosineDistance)
	}

}
