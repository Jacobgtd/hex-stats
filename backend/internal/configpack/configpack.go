package configpack

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// getSecretsDir returns the secrets directory path.
// It gets the secrets folder from SECRETS_DIR env key, defaulting to current_directory/secrets if not set.
// Returns an error if the current working directory cannot be determined.
func getSecretsDir() (string, error) {
	secretsDir := os.Getenv("SECRETS_DIR")
	if secretsDir == "" {
		// Default to current directory/secrets
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current working directory: %w", err)
		}
		secretsDir = filepath.Join(cwd, "secrets")
	}
	return secretsDir, nil
}

// Load loads environment variables from a file in the secrets directory.
// If the same variable is provided multiple times, later values overwrite earlier ones.
func Load(filename string) error {
	secretsDir, err := getSecretsDir()
	if err != nil {
		return err
	}

	filePath := filepath.Join(secretsDir, filename)

	// Load the file and overwrite existing env vars
	if err := godotenv.Load(filePath); err != nil {
		return fmt.Errorf("failed to load env file %s: %w", filePath, err)
	}

	return nil
}

// String returns the value of the environment variable key.
// Returns an error if the variable is not found.
func String(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("environment variable %s not found", key)
	}
	return value, nil
}

// StringOrDefault returns the value of the environment variable key, or defaultValue if not found.
func StringOrDefault(key string, defaultValue string) string {
	value, err := String(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// Int returns the value of the environment variable key as an integer.
// Returns an error if the variable is not found or cannot be parsed as an integer.
func Int(key string) (int, error) {
	value, err := String(key)
	if err != nil {
		return 0, err
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("failed to parse environment variable %s as int: %w", key, err)
	}

	return intValue, nil
}

// IntOrDefault returns the value of the environment variable key as an integer,
// or defaultValue if not found or cannot be parsed as an integer.
func IntOrDefault(key string, defaultValue int) int {
	value, err := String(key)
	if err != nil {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

// Bool returns the value of the environment variable key as a boolean.
// Returns an error if the variable is not found or cannot be parsed as a boolean.
func Bool(key string) (bool, error) {
	value, err := String(key)
	if err != nil {
		return false, err
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("failed to parse environment variable %s as bool: %w", key, err)
	}

	return boolValue, nil
}

// BoolOrDefault returns the value of the environment variable key as a boolean,
// or defaultValue if not found or cannot be parsed as a boolean.
func BoolOrDefault(key string, defaultValue bool) bool {
	value, err := String(key)
	if err != nil {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}

// LoadFile loads the contents of a file from the secrets directory into a string.
// Returns an error if the file cannot be read.
func LoadFile(filename string) (string, error) {
	secretsDir, err := getSecretsDir()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(secretsDir, filename)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return string(content), nil
}

func GetPath(filename string) (string, error) {
	secretsDir, err := getSecretsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(secretsDir, filename), nil
}

// StringSlice parses a comma-separated string into a slice of strings.
// Whitespace around each value is trimmed.
func StringSlice(value string) []string {
	if value == "" {
		return []string{}
	}
	parts := strings.Split(value, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// StringSliceFromEnv returns a comma-separated environment variable as a slice of strings.
// Returns an error if the variable is not found.
func StringSliceFromEnv(key string) ([]string, error) {
	value, err := String(key)
	if err != nil {
		return nil, err
	}
	return StringSlice(value), nil
}

// StringSliceFromEnvOrDefault returns a comma-separated environment variable as a slice of strings,
// or a default slice if not found or empty.
func StringSliceFromEnvOrDefault(key string, defaultValue []string) []string {
	value, err := String(key)
	if err != nil || value == "" {
		return defaultValue
	}
	return StringSlice(value)
}
