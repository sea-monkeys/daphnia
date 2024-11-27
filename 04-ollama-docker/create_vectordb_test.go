package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/ollama/ollama/api"
	"github.com/sea-monkeys/daphnia"
)

// Command represents the structure of each JSON object
type Command struct {
	Input       string `json:"input"`
	Instruction string `json:"instruction"`
	Output      string `json:"output"`
}

func TestGenerate(t *testing.T) {
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

	// Read the JSON file
	content, err := os.ReadFile("06102023.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Create a slice to store the commands
	var commands []Command

	// Unmarshal JSON data into the slice
	err = json.Unmarshal(content, &commands)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// Iterate through each command and display its contents
	// and save the embeddings to the vector store
	for i, cmd := range commands {
		fmt.Printf("Command #%d:\n", i+1)
		fmt.Printf("Input: %s\n", cmd.Input)
		fmt.Printf("Instruction: %s\n", cmd.Instruction)
		fmt.Printf("Output: %s\n", cmd.Output)
		fmt.Println("----------------------------------------")

		chunk := fmt.Sprintf("INPUT:\n%s\nOUTPUT:\n%s", cmd.Input, cmd.Output)

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
			Id:        fmt.Sprintf("Command #%d:\n", i+1),
		})

		if err != nil {
			fmt.Println("😡:", err)
		}

	}

	// Print total number of commands
	fmt.Printf("Total commands processed: %d\n", len(commands))

	// 2415
	if len(commands) != 2415 {
		t.Errorf("Expected 2415 commands, got %d", len(commands))
	}

}
