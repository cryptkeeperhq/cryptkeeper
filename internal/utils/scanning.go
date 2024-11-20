package utils

import (
	"regexp"
)

// Define a list of regex patterns for detecting secrets
var patterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)aws_secret_access_key\s*=\s*['"]?([A-Za-z0-9/+=]{40})['"]?`),
	regexp.MustCompile(`(?i)apikey\s*[:=]\s*['"]?([A-Za-z0-9]{32,})['"]?`),
	regexp.MustCompile(`(?i)password\s*[:=]\s*['"]?([^\s]+)['"]?`),
	// Add more patterns as needed
}

// ScanForSecrets scans the given text and returns potential secrets
func ScanForSecrets(text string) []string {
	var secrets []string
	for _, pattern := range patterns {
		matches := pattern.FindAllString(text, -1)
		secrets = append(secrets, matches...)
	}
	return secrets
}
