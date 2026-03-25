package monitor

import (
	"regexp"
	"strings"
)

// logLinePatterns are compiled once at startup and reused for every log line.
var logLinePatterns = []*regexp.Regexp{
	// MySQL general log: timestamp query
	regexp.MustCompile(`(?i)^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z\s+\d+\s+Query\s+(.+)$`),
	// PostgreSQL log
	regexp.MustCompile(`(?i)^.*LOG:\s+statement:\s+(.+)$`),
	// SQL Server
	regexp.MustCompile(`(?i)^.*exec\s+(.+)$`),
	// Generic SQL (if line starts with a DML/DDL keyword)
	regexp.MustCompile(`(?i)^(SELECT|INSERT|UPDATE|DELETE|CREATE|DROP|ALTER|MERGE|WITH)\s+.+$`),
}

// extractQueryFromLine extracts SQL query from a log line
func (p *LogProcessor) extractQueryFromLine(line string) string {
	line = strings.TrimSpace(line)

	for _, re := range logLinePatterns {
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
