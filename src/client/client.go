package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var models = []string{
	"google/gemini-2.0-flash-thinking-exp:free",
}

var wg sync.WaitGroup

// Function to get the response from the OpenRouter.ai API
func PromptAI(question string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY not found in environment variables")
	}

	res, err := queryModel(context.Background(), apiKey, models[0], question)

	if err != nil {
		res, err = queryModel(context.Background(), apiKey, models[1], question)
		if err != nil {
			log.Fatal("Error: ", err)
		}
	}
	wg.Wait()
	return res
}

// Function to make the HTTP request
func queryModel(ctx context.Context, apiKey, model string, content string) (string, error) {
	wg.Add(1)
	defer wg.Done()
	url := "https://openrouter.ai/api/v1/chat/completions"
	requestBody := RequestBody{
		Model: model,
		Messages: []Message{
			{Role: "user", Content: content},
		},
	}

	jsonValue, err := json.Marshal(requestBody)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	} else {
		return "No choices in response", errors.New("no choices in response")
	}
}
