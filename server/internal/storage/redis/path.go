package redis

import (
	"encoding/json"
	"strings"
)

func jsonPath(parts []string) string {
	var path strings.Builder
	path.WriteByte('$')
	for _, part := range parts {
		path.WriteByte('[')
		path.WriteString(quotePathKey(part))
		path.WriteByte(']')
	}

	return path.String()
}

func quoteKey(key string) string {
	encoded, _ := json.Marshal(key)
	return string(encoded)
}

func quotePathKey(key string) string {
	escaped := strings.ReplaceAll(key, "\\", "\\\\")
	return "\"" + strings.ReplaceAll(escaped, "\"", "\\\"") + "\""
}
