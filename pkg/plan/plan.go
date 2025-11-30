package plan

import "time"

// ExecutionPlan represents a query execution plan
type ExecutionPlan struct {
	Query         string          `json:"query"`
	Dialect       string          `json:"dialect"`
	RootNode      *PlanNode       `json:"root_node"`
	TotalCost     float64         `json:"total_cost"`
	EstimatedRows int64           `json:"estimated_rows"`
	ActualRows    int64           `json:"actual_rows,omitempty"`    // For EXPLAIN ANALYZE
	ExecutionTime time.Duration   `json:"execution_time,omitempty"` // For EXPLAIN ANALYZE
	Warnings      []string        `json:"warnings,omitempty"`
	Statistics    *PlanStatistics `json:"statistics,omitempty"`
}

// PlanNode represents a single node in the execution plan tree
type PlanNode struct {
	NodeType      NodeType       `json:"node_type"`
	Operation     string         `json:"operation"` // e.g., "Seq Scan", "Index Scan", "Hash Join"
	Table         string         `json:"table,omitempty"`
	Index         string         `json:"index,omitempty"`
	Condition     string         `json:"condition,omitempty"`
	Cost          *Cost          `json:"cost,omitempty"`
	Rows          *RowEstimate   `json:"rows,omitempty"`
	OutputColumns []string       `json:"output_columns,omitempty"`
	Children      []*PlanNode    `json:"children,omitempty"`
	Extra         map[string]any `json:"extra,omitempty"` // Dialect-specific fields
}

// NodeType represents the type of plan node
type NodeType string

const (
	// Scan operations
	NodeTypeSeqScan       NodeType = "SEQ_SCAN"
	NodeTypeIndexScan     NodeType = "INDEX_SCAN"
	NodeTypeIndexOnlyScan NodeType = "INDEX_ONLY_SCAN"
	NodeTypeBitmapScan    NodeType = "BITMAP_SCAN"

	// Join operations
	NodeTypeNestedLoop NodeType = "NESTED_LOOP"
	NodeTypeHashJoin   NodeType = "HASH_JOIN"
	NodeTypeMergeJoin  NodeType = "MERGE_JOIN"

	// Aggregation operations
	NodeTypeAggregate     NodeType = "AGGREGATE"
	NodeTypeGroupBy       NodeType = "GROUP_BY"
	NodeTypeHashAggregate NodeType = "HASH_AGGREGATE"

	// Sorting operations
	NodeTypeSort      NodeType = "SORT"
	NodeTypeQuickSort NodeType = "QUICKSORT"

	// Other operations
	NodeTypeFilter      NodeType = "FILTER"
	NodeTypeLimit       NodeType = "LIMIT"
	NodeTypeUnion       NodeType = "UNION"
	NodeTypeIntersect   NodeType = "INTERSECT"
	NodeTypeExcept      NodeType = "EXCEPT"
	NodeTypeSubquery    NodeType = "SUBQUERY"
	NodeTypeMaterialize NodeType = "MATERIALIZE"
	NodeTypeCTE         NodeType = "CTE"

	// SQL Server specific
	NodeTypeClusteredIndexScan    NodeType = "CLUSTERED_INDEX_SCAN"
	NodeTypeNonClusteredIndexScan NodeType = "NONCLUSTERED_INDEX_SCAN"
	NodeTypeTableScan             NodeType = "TABLE_SCAN"

	// MySQL specific
	NodeTypeFullTableScan NodeType = "FULL_TABLE_SCAN"
	NodeTypeRangeScan     NodeType = "RANGE_SCAN"
)

// Cost represents the cost estimates for a plan node
type Cost struct {
	StartupCost float64 `json:"startup_cost"` // Cost to get first row
	TotalCost   float64 `json:"total_cost"`   // Total cost to get all rows
	CPUCost     float64 `json:"cpu_cost,omitempty"`
	IOCost      float64 `json:"io_cost,omitempty"`
	NetworkCost float64 `json:"network_cost,omitempty"`
}

// RowEstimate represents row count estimates
type RowEstimate struct {
	Estimated int64   `json:"estimated"`          // Estimated rows
	Actual    int64   `json:"actual,omitempty"`   // Actual rows (EXPLAIN ANALYZE)
	Width     int     `json:"width,omitempty"`    // Average row width in bytes
	Accuracy  float64 `json:"accuracy,omitempty"` // Estimation accuracy (actual/estimated)
}

// PlanStatistics contains overall plan statistics
type PlanStatistics struct {
	TotalNodes         int               `json:"total_nodes"`
	ScanNodes          int               `json:"scan_nodes"`
	JoinNodes          int               `json:"join_nodes"`
	IndexScans         int               `json:"index_scans"`
	FullTableScans     int               `json:"full_table_scans"`
	TotalEstimatedRows int64             `json:"total_estimated_rows"`
	TotalActualRows    int64             `json:"total_actual_rows,omitempty"`
	MaxDepth           int               `json:"max_depth"`
	BottleneckNodes    []*BottleneckInfo `json:"bottleneck_nodes,omitempty"`
}

// BottleneckInfo represents a performance bottleneck in the plan
type BottleneckInfo struct {
	Node           *PlanNode `json:"node"`
	Issue          string    `json:"issue"`
	Severity       string    `json:"severity"` // CRITICAL, WARNING, INFO
	ImpactScore    float64   `json:"impact_score"`
	Recommendation string    `json:"recommendation"`
}

// IsExpensive returns true if the node is considered expensive
func (n *PlanNode) IsExpensive() bool {
	if n.Cost == nil {
		return false
	}
	// Heuristic: nodes with total cost > 1000 are expensive
	return n.Cost.TotalCost > 1000
}

// IsFullTableScan returns true if the node performs a full table scan
func (n *PlanNode) IsFullTableScan() bool {
	return n.NodeType == NodeTypeSeqScan ||
		n.NodeType == NodeTypeTableScan ||
		n.NodeType == NodeTypeFullTableScan
}

// IsIndexScan returns true if the node uses an index
func (n *PlanNode) IsIndexScan() bool {
	return n.NodeType == NodeTypeIndexScan ||
		n.NodeType == NodeTypeIndexOnlyScan ||
		n.NodeType == NodeTypeBitmapScan ||
		n.NodeType == NodeTypeClusteredIndexScan ||
		n.NodeType == NodeTypeNonClusteredIndexScan ||
		n.NodeType == NodeTypeRangeScan
}

// IsJoin returns true if the node is a join operation
func (n *PlanNode) IsJoin() bool {
	return n.NodeType == NodeTypeNestedLoop ||
		n.NodeType == NodeTypeHashJoin ||
		n.NodeType == NodeTypeMergeJoin
}

// Depth returns the depth of this node in the plan tree
func (n *PlanNode) Depth() int {
	if len(n.Children) == 0 {
		return 1
	}

	maxChildDepth := 0
	for _, child := range n.Children {
		childDepth := child.Depth()
		if childDepth > maxChildDepth {
			maxChildDepth = childDepth
		}
	}

	return maxChildDepth + 1
}

// CalculateStatistics calculates statistics for the execution plan
func (p *ExecutionPlan) CalculateStatistics() {
	stats := &PlanStatistics{
		BottleneckNodes: make([]*BottleneckInfo, 0),
	}

	p.collectStatistics(p.RootNode, stats, 1)
	p.Statistics = stats
}

// collectStatistics recursively collects statistics from plan nodes
func (p *ExecutionPlan) collectStatistics(node *PlanNode, stats *PlanStatistics, depth int) {
	if node == nil {
		return
	}

	stats.TotalNodes++

	if depth > stats.MaxDepth {
		stats.MaxDepth = depth
	}

	if node.IsFullTableScan() {
		stats.FullTableScans++
		stats.ScanNodes++
	}

	if node.IsIndexScan() {
		stats.IndexScans++
		stats.ScanNodes++
	}

	if node.IsJoin() {
		stats.JoinNodes++
	}

	if node.Rows != nil {
		stats.TotalEstimatedRows += node.Rows.Estimated
		if node.Rows.Actual > 0 {
			stats.TotalActualRows += node.Rows.Actual
		}
	}

	// Recursively process children
	for _, child := range node.Children {
		p.collectStatistics(child, stats, depth+1)
	}
}

// FindBottlenecks identifies performance bottlenecks in the plan
func (p *ExecutionPlan) FindBottlenecks() []*BottleneckInfo {
	bottlenecks := make([]*BottleneckInfo, 0)
	p.findBottlenecksRecursive(p.RootNode, &bottlenecks)
	return bottlenecks
}

// findBottlenecksRecursive recursively finds bottlenecks
func (p *ExecutionPlan) findBottlenecksRecursive(node *PlanNode, bottlenecks *[]*BottleneckInfo) {
	if node == nil {
		return
	}

	// Check for full table scans on potentially large tables
	if node.IsFullTableScan() && node.Rows != nil && node.Rows.Estimated > 1000 {
		*bottlenecks = append(*bottlenecks, &BottleneckInfo{
			Node:           node,
			Issue:          "Full table scan on large table",
			Severity:       "WARNING",
			ImpactScore:    float64(node.Rows.Estimated) / 1000.0,
			Recommendation: "Consider adding an index on the filtered columns",
		})
	}

	// Check for expensive operations
	if node.IsExpensive() {
		*bottlenecks = append(*bottlenecks, &BottleneckInfo{
			Node:           node,
			Issue:          "Expensive operation detected",
			Severity:       "WARNING",
			ImpactScore:    node.Cost.TotalCost / 1000.0,
			Recommendation: "Review query structure and indexes",
		})
	}

	// Check for nested loop joins with high row estimates
	if node.NodeType == NodeTypeNestedLoop && node.Rows != nil && node.Rows.Estimated > 10000 {
		*bottlenecks = append(*bottlenecks, &BottleneckInfo{
			Node:           node,
			Issue:          "Nested loop join with high cardinality",
			Severity:       "CRITICAL",
			ImpactScore:    float64(node.Rows.Estimated) / 10000.0,
			Recommendation: "Consider using hash join or merge join instead, or add appropriate indexes",
		})
	}

	// Check for poor estimation accuracy
	if node.Rows != nil && node.Rows.Actual > 0 {
		accuracy := float64(node.Rows.Actual) / float64(node.Rows.Estimated)
		if accuracy > 10.0 || accuracy < 0.1 {
			*bottlenecks = append(*bottlenecks, &BottleneckInfo{
				Node:           node,
				Issue:          "Poor cardinality estimation",
				Severity:       "INFO",
				ImpactScore:    accuracy,
				Recommendation: "Update table statistics or consider query rewrite",
			})
		}
	}

	// Recursively check children
	for _, child := range node.Children {
		p.findBottlenecksRecursive(child, bottlenecks)
	}
}
