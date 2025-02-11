package main

import (
	"fmt"

	"github.com/david22573/go-openrouter/src/client"
)

func main() {
	fmt.Println(client.PromptAI("Hello, how are you today?"))
}
