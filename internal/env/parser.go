package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ParseFile reads a .env file and returns a map of key-value pairs.
// Lines starting with '#' and empty lines are ignored.
func ParseFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open env file %q: %w", path, err)
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid syntax on line %d: %q", lineNum, line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		result[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading env file: %w", err)
	}
	return result, nil
}

// WriteFile serialises a key-value map to a .env file, sorted alphabetically.
func WriteFile(path string, data map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create env file %q: %w", path, err)
	}
	defer f.Close()

	keys := sortedKeys(data)
	w := bufio.NewWriter(f)
	for _, k := range keys {
		fmt.Fprintf(w, "%s=%s\n", k, data[k])
	}
	return w.Flush()
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// simple insertion sort — files are small
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}
