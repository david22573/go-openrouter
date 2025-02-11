package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

// Define structs for the request body
type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response struct for parsing the API response
type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

var models = []string{
	"google/gemini-2.0-flash-thinking-exp:free",
}

var wg sync.WaitGroup

// Function to get the response from the OpenRouter.ai API
func PromptAI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY not found in environment variables")
	}

	var filename string
	quizTemplateFile := "data/quiz-template.json"

	flag.StringVar(&filename, "p", "", "Read file to prompt openrouter.ai for response")
	flag.Parse()

	if filename == "" {
		panic(errors.New("no filename provided"))
	}

	file, err := os.ReadFile(filename)
	qfile, err := os.ReadFile(quizTemplateFile)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	content := string(file) + "\n" + string(qfile)
	res, err := queryModel(context.Background(), apiKey, models[0], content)

	if err != nil {
		res, err = queryModel(context.Background(), apiKey, models[1], content)
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
