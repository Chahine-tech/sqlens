package plan

import (
	"encoding/json"
	"fmt"
)

// PlanAnalyzer analyzes execution plans and provides optimization suggestions
type PlanAnalyzer struct {
	dialect string
}

// NewPlanAnalyzer creates a new plan analyzer
func NewPlanAnalyzer(dialect string) *PlanAnalyzer {
	return &PlanAnalyzer{
		dialect: dialect,
	}
}

// AnalyzePlan performs comprehensive analysis on an execution plan
func (pa *PlanAnalyzer) AnalyzePlan(plan *ExecutionPlan) *PlanAnalysis {
	analysis := &PlanAnalysis{
		Plan:            plan,
		Issues:          make([]*PlanIssue, 0),
		Recommendations: make([]*Recommendation, 0),
	}

	// Calculate statistics
	plan.CalculateStatistics()

	// Find bottlenecks
	bottlenecks := plan.FindBottlenecks()
	for _, bottleneck := range bottlenecks {
		analysis.Issues = append(analysis.Issues, &PlanIssue{
			Severity:    bottleneck.Severity,
			Type:        "BOTTLENECK",
			Description: bottleneck.Issue,
			Node:        bottleneck.Node,
			ImpactScore: bottleneck.ImpactScore,
		})

		analysis.Recommendations = append(analysis.Recommendations, &Recommendation{
			Type:        "OPTIMIZATION",
			Description: bottleneck.Recommendation,
			Priority:    pa.calculatePriority(bottleneck.Severity, bottleneck.ImpactScore),
		})
	}

	// Analyze plan structure
	pa.analyzeStructure(plan.RootNode, analysis)

	// Calculate overall score
	analysis.PerformanceScore = pa.calculatePerformanceScore(plan, analysis)

	return analysis
}

// analyzeStructure recursively analyzes the plan structure
func (pa *PlanAnalyzer) analyzeStructure(node *PlanNode, analysis *PlanAnalysis) {
	if node == nil {
		return
	}

	// Check for missing indexes
	if node.IsFullTableScan() && node.Rows != nil && node.Rows.Estimated > 100 {
		analysis.Issues = append(analysis.Issues, &PlanIssue{
			Severity:    "WARNING",
			Type:        "MISSING_INDEX",
			Description: fmt.Sprintf("Full table scan on '%s' with %d estimated rows", node.Table, node.Rows.Estimated),
			Node:        node,
			ImpactScore: float64(node.Rows.Estimated) / 1000.0,
		})

		analysis.Recommendations = append(analysis.Recommendations, &Recommendation{
			Type:        "INDEX",
			Description: fmt.Sprintf("Consider adding an index on table '%s' for columns used in filters or joins", node.Table),
			Priority:    "MEDIUM",
		})
	}

	// Check for inefficient joins
	if node.IsJoin() {
		pa.analyzeJoin(node, analysis)
	}

	// Check for sort operations
	if node.NodeType == NodeTypeSort || node.NodeType == NodeTypeQuickSort {
		if node.Rows != nil && node.Rows.Estimated > 10000 {
			analysis.Issues = append(analysis.Issues, &PlanIssue{
				Severity:    "INFO",
				Type:        "EXPENSIVE_SORT",
				Description: fmt.Sprintf("Sorting %d estimated rows", node.Rows.Estimated),
				Node:        node,
				ImpactScore: float64(node.Rows.Estimated) / 10000.0,
			})

			analysis.Recommendations = append(analysis.Recommendations, &Recommendation{
				Type:        "OPTIMIZATION",
				Description: "Consider adding an index to avoid sorting, or limit the result set before sorting",
				Priority:    "LOW",
			})
		}
	}

	// Recursively analyze children
	for _, child := range node.Children {
		pa.analyzeStructure(child, analysis)
	}
}

// analyzeJoin analyzes join operations
func (pa *PlanAnalyzer) analyzeJoin(node *PlanNode, analysis *PlanAnalysis) {
	if node.NodeType == NodeTypeNestedLoop {
		// Nested loop joins can be inefficient with large datasets
		if node.Rows != nil && node.Rows.Estimated > 5000 {
			analysis.Issues = append(analysis.Issues, &PlanIssue{
				Severity:    "WARNING",
				Type:        "INEFFICIENT_JOIN",
				Description: fmt.Sprintf("Nested loop join with %d estimated rows", node.Rows.Estimated),
				Node:        node,
				ImpactScore: float64(node.Rows.Estimated) / 5000.0,
			})

			analysis.Recommendations = append(analysis.Recommendations, &Recommendation{
				Type:        "JOIN_OPTIMIZATION",
				Description: "Consider using hash join or merge join for better performance with large datasets. Ensure appropriate indexes exist on join columns.",
				Priority:    "HIGH",
			})
		}
	}

	// Check for Cartesian products (joins without conditions)
	if node.Condition == "" && len(node.Children) >= 2 {
		analysis.Issues = append(analysis.Issues, &PlanIssue{
			Severity:    "CRITICAL",
			Type:        "CARTESIAN_PRODUCT",
			Description: "Join without condition detected - possible Cartesian product",
			Node:        node,
			ImpactScore: 10.0,
		})

		analysis.Recommendations = append(analysis.Recommendations, &Recommendation{
			Type:        "QUERY_REWRITE",
			Description: "Add explicit join conditions to avoid Cartesian product",
			Priority:    "CRITICAL",
		})
	}
}

// calculatePerformanceScore calculates an overall performance score (0-100)
func (pa *PlanAnalyzer) calculatePerformanceScore(plan *ExecutionPlan, analysis *PlanAnalysis) float64 {
	score := 100.0

	// Deduct points for issues
	for _, issue := range analysis.Issues {
		switch issue.Severity {
		case "CRITICAL":
			score -= 20.0
		case "WARNING":
			score -= 10.0
		case "INFO":
			score -= 5.0
		}
	}

	// Deduct points for high costs
	if plan.TotalCost > 10000 {
		score -= 10.0
	}

	// Deduct points for full table scans
	if plan.Statistics != nil {
		fullTableScanRatio := float64(plan.Statistics.FullTableScans) / float64(plan.Statistics.TotalNodes)
		score -= fullTableScanRatio * 20.0
	}

	// Ensure score doesn't go below 0
	if score < 0 {
		score = 0
	}

	return score
}

// calculatePriority calculates recommendation priority based on severity and impact
func (pa *PlanAnalyzer) calculatePriority(severity string, impact float64) string {
	if severity == "CRITICAL" || impact > 5.0 {
		return "CRITICAL"
	}
	if severity == "WARNING" || impact > 2.0 {
		return "HIGH"
	}
	if severity == "INFO" || impact > 1.0 {
		return "MEDIUM"
	}
	return "LOW"
}

// PlanAnalysis contains the results of plan analysis
type PlanAnalysis struct {
	Plan             *ExecutionPlan    `json:"plan"`
	Issues           []*PlanIssue      `json:"issues"`
	Recommendations  []*Recommendation `json:"recommendations"`
	PerformanceScore float64           `json:"performance_score"`
}

// PlanIssue represents an issue found in the execution plan
type PlanIssue struct {
	Severity    string    `json:"severity"` // CRITICAL, WARNING, INFO
	Type        string    `json:"type"`     // BOTTLENECK, MISSING_INDEX, etc.
	Description string    `json:"description"`
	Node        *PlanNode `json:"node,omitempty"`
	ImpactScore float64   `json:"impact_score"`
}

// Recommendation represents an optimization recommendation
type Recommendation struct {
	Type        string `json:"type"` // INDEX, OPTIMIZATION, QUERY_REWRITE, etc.
	Description string `json:"description"`
	Priority    string `json:"priority"` // CRITICAL, HIGH, MEDIUM, LOW
}

// ParseJSONPlan parses a JSON execution plan
// This is used for parsing output from EXPLAIN FORMAT=JSON
func ParseJSONPlan(jsonData []byte, dialect string) (*ExecutionPlan, error) {
	// The format varies by dialect, so we need to handle each one
	switch dialect {
	case "mysql":
		return parseMySQLJSONPlan(jsonData)
	case "postgresql", "postgres":
		return parsePostgreSQLJSONPlan(jsonData)
	case "sqlserver":
		return parseSQLServerXMLPlan(jsonData) // SQL Server uses XML
	case "sqlite":
		return parseSQLiteTextPlan(jsonData)
	default:
		return nil, fmt.Errorf("unsupported dialect for JSON plan parsing: %s", dialect)
	}
}

// parseMySQLJSONPlan parses MySQL JSON execution plan format
func parseMySQLJSONPlan(jsonData []byte) (*ExecutionPlan, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse MySQL JSON plan: %w", err)
	}

	plan := &ExecutionPlan{
		Dialect: "mysql",
	}

	// MySQL format: {"query_block": {...}}
	if queryBlock, ok := data["query_block"].(map[string]interface{}); ok {
		plan.RootNode = parseMySQLNode(queryBlock)
	}

	return plan, nil
}

// parseMySQLNode parses a MySQL plan node
func parseMySQLNode(data map[string]interface{}) *PlanNode {
	node := &PlanNode{
		Extra: make(map[string]any),
	}

	// Extract cost information
	if costInfo, ok := data["cost_info"].(map[string]interface{}); ok {
		node.Cost = &Cost{}
		if queryCost, ok := costInfo["query_cost"].(float64); ok {
			node.Cost.TotalCost = queryCost
		}
	}

	// Extract table information
	if table, ok := data["table"].(map[string]interface{}); ok {
		if tableName, ok := table["table_name"].(string); ok {
			node.Table = tableName
		}
		if accessType, ok := table["access_type"].(string); ok {
			node.Operation = accessType
			switch accessType {
			case "ALL":
				node.NodeType = NodeTypeFullTableScan
			case "index":
				node.NodeType = NodeTypeIndexScan
			case "range":
				node.NodeType = NodeTypeRangeScan
			default:
				node.NodeType = NodeTypeSeqScan
			}
		}
		if rows, ok := table["rows_examined_per_scan"].(float64); ok {
			node.Rows = &RowEstimate{
				Estimated: int64(rows),
			}
		}
	}

	return node
}

// parsePostgreSQLJSONPlan parses PostgreSQL JSON execution plan format
func parsePostgreSQLJSONPlan(jsonData []byte) (*ExecutionPlan, error) {
	var data []map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse PostgreSQL JSON plan: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("empty PostgreSQL plan")
	}

	plan := &ExecutionPlan{
		Dialect: "postgresql",
	}

	// PostgreSQL format: [{"Plan": {...}}]
	if planData, ok := data[0]["Plan"].(map[string]interface{}); ok {
		plan.RootNode = parsePostgreSQLNode(planData)

		if totalCost, ok := planData["Total Cost"].(float64); ok {
			plan.TotalCost = totalCost
		}
		if planRows, ok := planData["Plan Rows"].(float64); ok {
			plan.EstimatedRows = int64(planRows)
		}
	}

	return plan, nil
}

// parsePostgreSQLNode parses a PostgreSQL plan node
func parsePostgreSQLNode(data map[string]interface{}) *PlanNode {
	node := &PlanNode{
		Extra: make(map[string]any),
	}

	// Extract node type
	if nodeType, ok := data["Node Type"].(string); ok {
		node.Operation = nodeType
		switch nodeType {
		case "Seq Scan":
			node.NodeType = NodeTypeSeqScan
		case "Index Scan", "Index Only Scan":
			node.NodeType = NodeTypeIndexScan
		case "Bitmap Heap Scan", "Bitmap Index Scan":
			node.NodeType = NodeTypeBitmapScan
		case "Nested Loop":
			node.NodeType = NodeTypeNestedLoop
		case "Hash Join":
			node.NodeType = NodeTypeHashJoin
		case "Merge Join":
			node.NodeType = NodeTypeMergeJoin
		case "Aggregate":
			node.NodeType = NodeTypeAggregate
		case "Sort":
			node.NodeType = NodeTypeSort
		default:
			node.NodeType = NodeType(nodeType)
		}
	}

	// Extract cost
	if startupCost, ok := data["Startup Cost"].(float64); ok {
		if node.Cost == nil {
			node.Cost = &Cost{}
		}
		node.Cost.StartupCost = startupCost
	}
	if totalCost, ok := data["Total Cost"].(float64); ok {
		if node.Cost == nil {
			node.Cost = &Cost{}
		}
		node.Cost.TotalCost = totalCost
	}

	// Extract row estimates
	if planRows, ok := data["Plan Rows"].(float64); ok {
		if node.Rows == nil {
			node.Rows = &RowEstimate{}
		}
		node.Rows.Estimated = int64(planRows)
	}
	if actualRows, ok := data["Actual Rows"].(float64); ok {
		if node.Rows == nil {
			node.Rows = &RowEstimate{}
		}
		node.Rows.Actual = int64(actualRows)
	}

	// Extract table name
	if relationName, ok := data["Relation Name"].(string); ok {
		node.Table = relationName
	}

	// Extract index name
	if indexName, ok := data["Index Name"].(string); ok {
		node.Index = indexName
	}

	// Parse children
	if plans, ok := data["Plans"].([]interface{}); ok {
		node.Children = make([]*PlanNode, 0, len(plans))
		for _, childData := range plans {
			if childMap, ok := childData.(map[string]interface{}); ok {
				node.Children = append(node.Children, parsePostgreSQLNode(childMap))
			}
		}
	}

	return node
}

// parseSQLServerXMLPlan parses SQL Server XML execution plan format
func parseSQLServerXMLPlan(xmlData []byte) (*ExecutionPlan, error) {
	// TODO: Implement SQL Server XML plan parsing
	// SQL Server uses XML format which requires XML parsing
	return &ExecutionPlan{
		Dialect:  "sqlserver",
		Warnings: []string{"SQL Server XML plan parsing not yet implemented"},
	}, nil
}

// parseSQLiteTextPlan parses SQLite text execution plan format
func parseSQLiteTextPlan(textData []byte) (*ExecutionPlan, error) {
	// TODO: Implement SQLite text plan parsing
	// SQLite uses a simple text format
	return &ExecutionPlan{
		Dialect:  "sqlite",
		Warnings: []string{"SQLite text plan parsing not yet implemented"},
	}, nil
}
