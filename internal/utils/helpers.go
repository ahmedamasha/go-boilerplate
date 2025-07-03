package utils

// Helper functions for working with metadata maps

// GetInt safely extracts an integer value from metadata map
func GetInt(metadata map[string]interface{}, key string) int {
	if val, ok := metadata[key].(int); ok {
		return val
	}
	return 0
}

// GetBool safely extracts a boolean value from metadata map
func GetBool(metadata map[string]interface{}, key string) bool {
	if val, ok := metadata[key].(bool); ok {
		return val
	}
	return false
}

// GetString safely extracts a string value from metadata map
func GetString(metadata map[string]interface{}, key string) string {
	if val, ok := metadata[key].(string); ok {
		return val
	}
	return ""
}

// GetFloat safely extracts a float64 value from metadata map
func GetFloat(metadata map[string]interface{}, key string) float64 {
	if val, ok := metadata[key].(float64); ok {
		return val
	}
	return 0.0
}
