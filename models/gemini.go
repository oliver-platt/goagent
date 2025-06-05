package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/oliver-platt/goagent/v2/types"
)

// GeminiModel implements the Model interface for Google's Gemini API
type GeminiModel struct {
	apiKey     string
	modelName  string
	httpClient *http.Client
	baseURL    string
}

// NewGeminiModel creates a new Gemini model instance
func NewGeminiModel() (*GeminiModel, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
	}

	return &GeminiModel{
		apiKey:    apiKey,
		modelName: "gemini-1.5-flash",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://generativelanguage.googleapis.com/v1beta",
	}, nil
}

// NewGeminiModelWithConfig creates a Gemini model with custom configuration
func NewGeminiModelWithConfig(apiKey string, timeout time.Duration) *GeminiModel {
	return &GeminiModel{
		apiKey:    apiKey,
		modelName: "gemini-1.5-flash",
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL: "https://generativelanguage.googleapis.com/v1beta",
	}
}

// Generate implements the Model interface
func (g *GeminiModel) Generate(ctx context.Context, messages []types.Message) (string, error) {
	// Convert our messages to Gemini format
	geminiMessages, err := g.convertMessages(messages)
	if err != nil {
		return "", fmt.Errorf("failed to convert messages: %w", err)
	}

	// Build the request payload
	payload := GeminiRequest{
		Contents: geminiMessages,
		GenerationConfig: GenerationConfig{
			Temperature:     0.7,
			TopK:            40,
			TopP:            0.95,
			MaxOutputTokens: 2048,
		},
	}

	// Make the API request
	response, err := g.makeRequest(ctx, payload)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}

	// Extract the response text
	if len(response.Candidates) == 0 {
		return "", fmt.Errorf("no response candidates returned")
	}

	if len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response parts returned")
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}

// Name returns the model name
func (g *GeminiModel) Name() string {
	return g.modelName
}

// convertMessages converts our message format to Gemini's format
func (g *GeminiModel) convertMessages(messages []types.Message) ([]GeminiMessage, error) {
	var geminiMessages []GeminiMessage
	var systemInstructions []string

	for _, msg := range messages {
		switch msg.Role {
		case types.RoleSystem:
			// Gemini handles system messages differently - we'll prepend them to the first user message
			systemInstructions = append(systemInstructions, msg.Content)
		case types.RoleUser:
			content := msg.Content
			// Prepend system instructions to the first user message
			if len(systemInstructions) > 0 {
				content = strings.Join(systemInstructions, "\n\n") + "\n\n" + content
				systemInstructions = nil // Clear after using
			}
			geminiMessages = append(geminiMessages, GeminiMessage{
				Role: "user",
				Parts: []GeminiPart{
					{Text: content},
				},
			})
		case types.RoleAssistant:
			geminiMessages = append(geminiMessages, GeminiMessage{
				Role: "model", // Gemini uses "model" instead of "assistant"
				Parts: []GeminiPart{
					{Text: msg.Content},
				},
			})
		}
	}

	if len(geminiMessages) == 0 {
		return nil, fmt.Errorf("no valid messages to send")
	}

	return geminiMessages, nil
}

// makeRequest sends the request to the Gemini API
func (g *GeminiModel) makeRequest(ctx context.Context, payload GeminiRequest) (*GeminiResponse, error) {
	// Marshal the payload
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Build the URL
	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", g.baseURL, g.modelName, g.apiKey)

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var geminiResponse GeminiResponse
	if err := json.Unmarshal(body, &geminiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &geminiResponse, nil
}

// Gemini API types

type GeminiRequest struct {
	Contents         []GeminiMessage  `json:"contents"`
	GenerationConfig GenerationConfig `json:"generationConfig,omitempty"`
}

type GeminiMessage struct {
	Role  string       `json:"role"`
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GenerationConfig struct {
	Temperature     float64 `json:"temperature,omitempty"`
	TopK            int     `json:"topK,omitempty"`
	TopP            float64 `json:"topP,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}
