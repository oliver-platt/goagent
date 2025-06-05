package goagent

import (
	"context"
	"fmt"

	"github.com/oliver-platt/goagent/v2/models"
	"github.com/oliver-platt/goagent/v2/types"
)

// Agent represents an AI agent with a system prompt and model
type Agent struct {
	SystemPrompt string
	Model        models.Model
}

// NewAgent creates a new agent with the given system prompt and model
func NewAgent(systemPrompt string, model models.Model) *Agent {
	return &Agent{
		SystemPrompt: systemPrompt,
		Model:        model,
	}
}

// Run executes the agent with user input and returns the response
func (a *Agent) Run(ctx context.Context, userInput string) (string, error) {
	if a.Model == nil {
		return "", fmt.Errorf("agent has no model")
	}

	if userInput == "" {
		return "", fmt.Errorf("user input cannot be empty")
	}

	// Build the messages for the model
	messages := a.buildMessages(userInput)

	// Generate response using the model
	return a.Model.Generate(ctx, messages)
}

// buildMessages constructs the message array for the model
func (a *Agent) buildMessages(userInput string) []types.Message {
	var messages []types.Message

	// Add system message if we have one
	if a.SystemPrompt != "" {
		messages = append(messages, types.Message{
			Role:    types.RoleSystem,
			Content: a.SystemPrompt,
		})
	}

	// Add user message
	messages = append(messages, types.Message{
		Role:    types.RoleUser,
		Content: userInput,
	})

	return messages
}

// GetModelName returns the name of the underlying model
func (a *Agent) GetModelName() string {
	if a.Model == nil {
		return "none"
	}
	return a.Model.Name()
}

// SetSystemPrompt updates the agent's system prompt
func (a *Agent) SetSystemPrompt(prompt string) {
	a.SystemPrompt = prompt
}

// SetModel updates the agent's model
func (a *Agent) SetModel(model models.Model) {
	a.Model = model
}
