package environ

import (
	"os"

	"github.com/CzarSimon/text-service/go/pkg/utils/logger"
)

var log = logger.GetDefaultLogger("pkg/environ").Sugar()

// Get gets environment variable with default value.
func Get(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

// MustGet gets environment variable and panics if not found
func MustGet(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Panicw("Failed to get environment variable", "name", name)
	}

	return value
}
