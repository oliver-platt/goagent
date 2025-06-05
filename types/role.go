package types

type Role string

// Enum values for Role
const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

// String returns the string representation of the role
func (r Role) String() string {
	return string(r)
}

// IsValid checks if the role is one of the valid enum values
func (r Role) IsValid() bool {
	switch r {
	case RoleUser, RoleAssistant, RoleSystem:
		return true
	default:
		return false
	}
}
