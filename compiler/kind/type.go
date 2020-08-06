package kind

import "strings"

// IsArray tests for an array kind.
func IsArray(ty string) bool {
	return strings.HasPrefix(ty, "[]")
}

// IsMap tests for a map kind.
func IsMap(ty string) bool {
	return strings.HasPrefix(ty, "{}")
}

// IsFunc tests for a function literal.
func IsFunc(ty string) bool {
	return strings.HasPrefix(ty, "func ")
}

// ElementType returns the element type of a value of an array or map.
func ElementType(ty string) string {
	return ty[2:]
}
