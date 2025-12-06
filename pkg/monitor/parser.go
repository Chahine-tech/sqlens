package monitor

import (
	"regexp"
	"strings"
)

// extractQueryFromLine extracts SQL query from a log line
func (p *LogProcessor) extractQueryFromLine(line string) string {
	// Remove timestamp prefix if present
	line = strings.TrimSpace(line)

	// Common log patterns
	patterns := []string{
		// MySQL general log: timestamp query
		`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z\s+\d+\s+Query\s+(.+)$`,
		// PostgreSQL log
		`^.*LOG:\s+statement:\s+(.+)$`,
		// SQL Server
		`^.*exec\s+(.+)$`,
		// Generic SQL (if line starts with SELECT, INSERT, UPDATE, DELETE, etc.)
		`^(SELECT|INSERT|UPDATE|DELETE|CREATE|DROP|ALTER|MERGE|WITH)\s+.+$`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		if matches := re.FindStringSubmatch(line); len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	// If the line itself looks like SQL, return it
	upperLine := strings.ToUpper(strings.TrimSpace(line))
	if strings.HasPrefix(upperLine, "SELECT") ||
		strings.HasPrefix(upperLine, "INSERT") ||
		strings.HasPrefix(upperLine, "UPDATE") ||
		strings.HasPrefix(upperLine, "DELETE") ||
		strings.HasPrefix(upperLine, "CREATE") ||
		strings.HasPrefix(upperLine, "DROP") ||
		strings.HasPrefix(upperLine, "ALTER") ||
		strings.HasPrefix(upperLine, "MERGE") ||
		strings.HasPrefix(upperLine, "WITH") {
		return strings.TrimSpace(line)
	}

	return ""
}
