package monitor

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

// AlertLevel represents the severity of an alert
type AlertLevel int

const (
	AlertInfo AlertLevel = iota
	AlertWarning
	AlertError
	AlertCritical
)

func (a AlertLevel) String() string {
	switch a {
	case AlertInfo:
		return "INFO"
	case AlertWarning:
		return "WARNING"
	case AlertError:
		return "ERROR"
	case AlertCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// Alert represents a monitoring alert
type Alert struct {
	Level     AlertLevel
	Type      string
	Message   string
	Query     *ProcessedQuery
	Timestamp time.Time
	Metadata  map[string]interface{}
}

// AlertRule defines conditions for triggering alerts
type AlertRule interface {
	Check(pq *ProcessedQuery) *Alert
	Name() string
}

// AlertManager manages alert rules and notifications
type AlertManager struct {
	rules    []AlertRule
	handlers []AlertHandler
	mu       sync.RWMutex

	// Alert statistics
	alertCount map[AlertLevel]int64
	statsMu    sync.RWMutex
}

// AlertHandler handles alerts when they are triggered
type AlertHandler func(*Alert)

// NewAlertManager creates a new alert manager
func NewAlertManager() *AlertManager {
	return &AlertManager{
		rules:      []AlertRule{},
		handlers:   []AlertHandler{},
		alertCount: make(map[AlertLevel]int64),
	}
}

// AddRule adds an alert rule
func (am *AlertManager) AddRule(rule AlertRule) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.rules = append(am.rules, rule)
}

// AddHandler adds an alert handler
func (am *AlertManager) AddHandler(handler AlertHandler) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.handlers = append(am.handlers, handler)
}

// Check checks all rules against a processed query
func (am *AlertManager) Check(pq *ProcessedQuery) {
	am.mu.RLock()
	rules := am.rules
	handlers := am.handlers
	am.mu.RUnlock()

	for _, rule := range rules {
		if alert := rule.Check(pq); alert != nil {
			// Update statistics
			am.statsMu.Lock()
			am.alertCount[alert.Level]++
			am.statsMu.Unlock()

			// Trigger handlers
			for _, handler := range handlers {
				handler(alert)
			}
		}
	}
}

// GetAlertCounts returns the count of alerts by level
func (am *AlertManager) GetAlertCounts() map[AlertLevel]int64 {
	am.statsMu.RLock()
	defer am.statsMu.RUnlock()

	counts := make(map[AlertLevel]int64)
	for level, count := range am.alertCount {
		counts[level] = count
	}
	return counts
}

// SlowQueryRule alerts on queries exceeding a duration threshold
type SlowQueryRule struct {
	Threshold float64 // in seconds
}

func (r *SlowQueryRule) Name() string {
	return "SlowQueryRule"
}

func (r *SlowQueryRule) Check(pq *ProcessedQuery) *Alert {
	if pq.Duration >= r.Threshold {
		level := AlertWarning
		if pq.Duration >= r.Threshold*2 {
			level = AlertError
		}
		if pq.Duration >= r.Threshold*5 {
			level = AlertCritical
		}

		return &Alert{
			Level:     level,
			Type:      "SLOW_QUERY",
			Message:   fmt.Sprintf("Query took %.2fs (threshold: %.2fs)", pq.Duration, r.Threshold),
			Query:     pq,
			Timestamp: time.Now(),
			Metadata: map[string]interface{}{
				"duration":  pq.Duration,
				"threshold": r.Threshold,
			},
		}
	}
	return nil
}

// ParseErrorRule alerts on query parse failures
type ParseErrorRule struct{}

func (r *ParseErrorRule) Name() string {
	return "ParseErrorRule"
}

func (r *ParseErrorRule) Check(pq *ProcessedQuery) *Alert {
	if pq.Statement == nil && pq.Query != "" {
		return &Alert{
			Level:     AlertWarning,
			Type:      "PARSE_ERROR",
			Message:   "Failed to parse SQL query",
			Query:     pq,
			Timestamp: time.Now(),
		}
	}
	return nil
}

// OptimizationRule alerts on queries with optimization opportunities
type OptimizationRule struct {
	MinSeverity string // "low", "medium", "high"
}

func (r *OptimizationRule) Name() string {
	return "OptimizationRule"
}

func (r *OptimizationRule) Check(pq *ProcessedQuery) *Alert {
	if pq.Analysis == nil || len(pq.Analysis.EnhancedSuggestions) == 0 {
		return nil
	}

	// Check if any optimization meets severity threshold
	hasHighSeverity := false
	for _, opt := range pq.Analysis.EnhancedSuggestions {
		severity := strings.ToLower(opt.Severity)
		if r.MinSeverity == "low" ||
			(r.MinSeverity == "medium" && (severity == "medium" || severity == "high")) ||
			(r.MinSeverity == "high" && severity == "high") {
			hasHighSeverity = true
			break
		}
	}

	if hasHighSeverity {
		return &Alert{
			Level:     AlertInfo,
			Type:      "OPTIMIZATION_OPPORTUNITY",
			Message:   fmt.Sprintf("Query has %d optimization suggestions", len(pq.Analysis.EnhancedSuggestions)),
			Query:     pq,
			Timestamp: time.Now(),
			Metadata: map[string]interface{}{
				"optimization_count": len(pq.Analysis.EnhancedSuggestions),
			},
		}
	}

	return nil
}

// FullTableScanRule alerts on queries doing full table scans
type FullTableScanRule struct{}

func (r *FullTableScanRule) Name() string {
	return "FullTableScanRule"
}

func (r *FullTableScanRule) Check(pq *ProcessedQuery) *Alert {
	if pq.Analysis == nil {
		return nil
	}

	// Check for missing WHERE clause on SELECT
	if selectStmt, ok := pq.Statement.(*parser.SelectStatement); ok {
		if selectStmt.Where == nil && len(pq.Analysis.Tables) > 0 {
			return &Alert{
				Level:     AlertWarning,
				Type:      "FULL_TABLE_SCAN",
				Message:   "SELECT query without WHERE clause may cause full table scan",
				Query:     pq,
				Timestamp: time.Now(),
				Metadata: map[string]interface{}{
					"tables": pq.Analysis.Tables,
				},
			}
		}
	}

	// Check for missing WHERE clause on UPDATE/DELETE
	if _, ok := pq.Statement.(*parser.UpdateStatement); ok {
		return &Alert{
			Level:     AlertError,
			Type:      "UNSAFE_UPDATE",
			Message:   "UPDATE query without WHERE clause will affect all rows",
			Query:     pq,
			Timestamp: time.Now(),
		}
	}

	if _, ok := pq.Statement.(*parser.DeleteStatement); ok {
		return &Alert{
			Level:     AlertError,
			Type:      "UNSAFE_DELETE",
			Message:   "DELETE query without WHERE clause will remove all rows",
			Query:     pq,
			Timestamp: time.Now(),
		}
	}

	return nil
}

// ConsoleAlertHandler prints alerts to console
func ConsoleAlertHandler(alert *Alert) {
	fmt.Printf("[%s] %s: %s\n",
		alert.Level.String(),
		alert.Type,
		alert.Message)

	if alert.Query != nil {
		fmt.Printf("  Query: %s\n", truncateString(alert.Query.Query, 100))
		if alert.Query.Duration > 0 {
			fmt.Printf("  Duration: %.2fs\n", alert.Query.Duration)
		}
	}
	fmt.Println()
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
