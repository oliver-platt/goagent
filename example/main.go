package main

import (
	"context"
	"fmt"

	"github.com/oliver-platt/goagent/v2"
	"github.com/oliver-platt/goagent/v2/models"
)

func main() {

	model, err := models.NewGeminiModel()

	if err != nil {
		fmt.Printf("❌ Failed to create Gemini model: %v\n", err)
		return
	}

	systemPrompt := "You are a helpful assistant. Be friendly and concise in your responses."
	agent := goagent.NewAgent(systemPrompt, model)

	ctx := context.Background()

	input := "Tell me a joke"
	response, err := agent.Run(ctx, input)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	}
	fmt.Printf("✅ Agent Response: %s\n", response)
}
