package types

import (
	"fmt"
	"strings"
)

// Message represents a single message in a conversation
type Message struct {
	Role    Role   `json:"role"`    // The role of the message sender
	Content string `json:"content"` // The message content
}

// NewMessage creates a new message with the given role and content
func NewMessage(role Role, content string) Message {
	return Message{
		Role:    role,
		Content: content,
	}
}

// NewUserMessage creates a new user message
func NewUserMessage(content string) Message {
	return Message{
		Role:    RoleUser,
		Content: content,
	}
}

// NewAssistantMessage creates a new assistant message
func NewAssistantMessage(content string) Message {
	return Message{
		Role:    RoleAssistant,
		Content: content,
	}
}

// NewSystemMessage creates a new system message
func NewSystemMessage(content string) Message {
	return Message{
		Role:    RoleSystem,
		Content: content,
	}
}

// IsUser returns true if the message is from a user
func (m Message) IsUser() bool {
	return m.Role == RoleUser
}

// IsAssistant returns true if the message is from an assistant
func (m Message) IsAssistant() bool {
	return m.Role == RoleAssistant
}

// IsSystem returns true if the message is a system message
func (m Message) IsSystem() bool {
	return m.Role == RoleSystem
}

// IsEmpty returns true if the message has no content
func (m Message) IsEmpty() bool {
	return strings.TrimSpace(m.Content) == ""
}

// String returns a formatted string representation of the message
func (m Message) String() string {
	return fmt.Sprintf("[%s] %s", strings.ToUpper(string(m.Role)), m.Content)
}

// Validate checks if the message has valid role and content
func (m Message) Validate() error {
	// Check if role is valid using the enum's IsValid method
	if !m.Role.IsValid() {
		return fmt.Errorf("invalid role '%s', must be one of: %s, %s, %s",
			m.Role, RoleUser, RoleAssistant, RoleSystem)
	}

	// Check if content is not empty (except for system messages which can be empty)
	if m.IsEmpty() && !m.IsSystem() {
		return fmt.Errorf("message content cannot be empty for role '%s'", m.Role)
	}

	return nil
}

// Truncate returns a truncated version of the message content
func (m Message) Truncate(maxLen int) Message {
	if len(m.Content) <= maxLen {
		return m
	}

	truncated := m.Content[:maxLen-3] + "..."
	return Message{
		Role:    m.Role,
		Content: truncated,
	}
}

// WordCount returns the number of words in the message content
func (m Message) WordCount() int {
	if m.IsEmpty() {
		return 0
	}

	words := strings.Fields(m.Content)
	return len(words)
}
