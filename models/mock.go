package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/oliver-platt/goagent/v2/types"
)

// MockModel is a simple implementation for testing and development
type MockModel struct {
	name string
}

// NewMockModel creates a new mock model
func NewMockModel() *MockModel {
	return &MockModel{
		name: "mock-model",
	}
}

// Generate implements the Model interface with simple response logic
func (m *MockModel) Generate(ctx context.Context, messages []types.Message) (string, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	if len(messages) == 0 {
		return "", fmt.Errorf("no messages provided")
	}

	// Find the last user message
	var userMessage string
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == types.RoleUser {
			userMessage = messages[i].Content
			break
		}
	}

	if userMessage == "" {
		return "I didn't receive any user input.", nil
	}

	// Simple response logic based on keywords
	userLower := strings.ToLower(userMessage)

	switch {
	case strings.Contains(userLower, "hello") || strings.Contains(userLower, "hi"):
		return "Hello! How can I help you today?", nil

	case strings.Contains(userLower, "weather"):
		return "I'd love to help with weather information, but I don't have access to weather data yet.", nil

	case strings.Contains(userLower, "time"):
		return "I don't have access to the current time, but I can help with other questions!", nil

	case strings.Contains(userLower, "math") || strings.Contains(userLower, "calculate"):
		return "I can help with math! Try asking me specific calculations like '2 + 2'.", nil

	case strings.Contains(userLower, "2 + 2") || strings.Contains(userLower, "2+2"):
		return "2 + 2 = 4", nil

	case strings.Contains(userLower, "go") && strings.Contains(userLower, "programming"):
		return "Go is a fantastic programming language! It's simple, fast, and has great concurrency support.", nil

	case strings.Contains(userLower, "thank"):
		return "You're welcome! Happy to help.", nil

	case strings.Contains(userLower, "bye") || strings.Contains(userLower, "goodbye"):
		return "Goodbye! Have a great day!", nil

	default:
		return fmt.Sprintf("I understand you said: \"%s\". I'm a simple mock model, so my responses are limited, but I'm here to help!", userMessage), nil
	}
}

// Name returns the model's name
func (m *MockModel) Name() string {
	return m.name
}

// SetName allows updating the model's name (useful for testing)
func (m *MockModel) SetName(name string) {
	m.name = name
}
