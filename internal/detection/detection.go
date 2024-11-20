package detection

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var secretPatterns = []string{
	`(?i)api[_-]?key\s*=\s*['"]?[\w-]+['"]?`,
	`(?i)secret[_-]?key\s*=\s*['"]?[\w-]+['"]?`,
	`(?i)password\s*=\s*['"]?[\w-]+['"]?`,
	// Add more patterns as needed
}

// ScanFile scans a file for secrets.
func ScanFile(filePath string) ([]string, error) {
	var matches []string

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	for _, pattern := range secretPatterns {
		re := regexp.MustCompile(pattern)
		found := re.FindAllString(string(fileContent), -1)
		matches = append(matches, found...)
	}

	return matches, nil
}

// ScanDirectory scans a directory for secrets.
func ScanDirectory(dirPath string) (map[string][]string, error) {
	results := make(map[string][]string)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") { // Change this to match the file types you want to scan
			matches, err := ScanFile(path)
			if err != nil {
				return err
			}

			if len(matches) > 0 {
				results[path] = matches
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}
