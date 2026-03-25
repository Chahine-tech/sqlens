package tests

import (
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/plan"
)

func makeSeqScanNode(table string, rows int64) *plan.PlanNode {
	return &plan.PlanNode{
		NodeType:  plan.NodeTypeSeqScan,
		Operation: "Seq Scan",
		Table:     table,
		Cost:      &plan.Cost{StartupCost: 0, TotalCost: float64(rows) * 0.1},
		Rows:      &plan.RowEstimate{Estimated: rows},
	}
}

func makeIndexScanNode(table string, rows int64) *plan.PlanNode {
	return &plan.PlanNode{
		NodeType:  plan.NodeTypeIndexScan,
		Operation: "Index Scan",
		Table:     table,
		Cost:      &plan.Cost{StartupCost: 0.5, TotalCost: float64(rows) * 0.01},
		Rows:      &plan.RowEstimate{Estimated: rows},
	}
}

// --- PlanNode methods ---

func TestPlanNode_IsFullTableScan(t *testing.T) {
	cases := []struct {
		nt       plan.NodeType
		expected bool
	}{
		{plan.NodeTypeSeqScan, true},
		{plan.NodeTypeTableScan, true},
		{plan.NodeTypeFullTableScan, true},
		{plan.NodeTypeIndexScan, false},
		{plan.NodeTypeHashJoin, false},
	}
	for _, tc := range cases {
		n := &plan.PlanNode{NodeType: tc.nt}
		if n.IsFullTableScan() != tc.expected {
			t.Errorf("NodeType %s: IsFullTableScan() = %v, want %v", tc.nt, n.IsFullTableScan(), tc.expected)
		}
	}
}

func TestPlanNode_IsIndexScan(t *testing.T) {
	scanTypes := []plan.NodeType{
		plan.NodeTypeIndexScan,
		plan.NodeTypeIndexOnlyScan,
		plan.NodeTypeBitmapScan,
		plan.NodeTypeClusteredIndexScan,
		plan.NodeTypeNonClusteredIndexScan,
		plan.NodeTypeRangeScan,
	}
	for _, nt := range scanTypes {
		n := &plan.PlanNode{NodeType: nt}
		if !n.IsIndexScan() {
			t.Errorf("NodeType %s should be an index scan", nt)
		}
	}
	n := &plan.PlanNode{NodeType: plan.NodeTypeSeqScan}
	if n.IsIndexScan() {
		t.Error("SeqScan should not be an index scan")
	}
}

func TestPlanNode_IsJoin(t *testing.T) {
	joinTypes := []plan.NodeType{
		plan.NodeTypeNestedLoop,
		plan.NodeTypeHashJoin,
		plan.NodeTypeMergeJoin,
	}
	for _, nt := range joinTypes {
		n := &plan.PlanNode{NodeType: nt}
		if !n.IsJoin() {
			t.Errorf("NodeType %s should be a join", nt)
		}
	}
	n := &plan.PlanNode{NodeType: plan.NodeTypeSeqScan}
	if n.IsJoin() {
		t.Error("SeqScan should not be a join")
	}
}

func TestPlanNode_IsExpensive(t *testing.T) {
	cheap := &plan.PlanNode{Cost: &plan.Cost{TotalCost: 500}}
	if cheap.IsExpensive() {
		t.Error("cost 500 should not be expensive")
	}
	expensive := &plan.PlanNode{Cost: &plan.Cost{TotalCost: 1001}}
	if !expensive.IsExpensive() {
		t.Error("cost 1001 should be expensive")
	}
	noCost := &plan.PlanNode{}
	if noCost.IsExpensive() {
		t.Error("node without cost should not be expensive")
	}
}

func TestPlanNode_Depth(t *testing.T) {
	leaf := &plan.PlanNode{NodeType: plan.NodeTypeSeqScan}
	if leaf.Depth() != 1 {
		t.Errorf("leaf depth = %d, want 1", leaf.Depth())
	}

	parent := &plan.PlanNode{
		NodeType: plan.NodeTypeHashJoin,
		Children: []*plan.PlanNode{
			{NodeType: plan.NodeTypeSeqScan},
			{NodeType: plan.NodeTypeIndexScan},
		},
	}
	if parent.Depth() != 2 {
		t.Errorf("parent depth = %d, want 2", parent.Depth())
	}

	// 3 levels deep
	root := &plan.PlanNode{
		NodeType: plan.NodeTypeNestedLoop,
		Children: []*plan.PlanNode{parent},
	}
	if root.Depth() != 3 {
		t.Errorf("root depth = %d, want 3", root.Depth())
	}
}

// --- ExecutionPlan.CalculateStatistics ---

func TestExecutionPlan_CalculateStatistics(t *testing.T) {
	ep := &plan.ExecutionPlan{
		RootNode: &plan.PlanNode{
			NodeType:  plan.NodeTypeHashJoin,
			Operation: "Hash Join",
			Rows:      &plan.RowEstimate{Estimated: 500},
			Children: []*plan.PlanNode{
				makeSeqScanNode("users", 1000),
				makeIndexScanNode("orders", 200),
			},
		},
	}
	ep.CalculateStatistics()
	stats := ep.Statistics
	if stats == nil {
		t.Fatal("expected statistics to be set")
	}
	if stats.TotalNodes != 3 {
		t.Errorf("TotalNodes = %d, want 3", stats.TotalNodes)
	}
	if stats.JoinNodes != 1 {
		t.Errorf("JoinNodes = %d, want 1", stats.JoinNodes)
	}
	if stats.FullTableScans != 1 {
		t.Errorf("FullTableScans = %d, want 1", stats.FullTableScans)
	}
	if stats.IndexScans != 1 {
		t.Errorf("IndexScans = %d, want 1", stats.IndexScans)
	}
	if stats.MaxDepth != 2 {
		t.Errorf("MaxDepth = %d, want 2", stats.MaxDepth)
	}
}

func TestExecutionPlan_CalculateStatistics_NilRoot(t *testing.T) {
	ep := &plan.ExecutionPlan{}
	ep.CalculateStatistics() // should not panic
}

func TestExecutionPlan_CalculateStatistics_WithActualRows(t *testing.T) {
	ep := &plan.ExecutionPlan{
		RootNode: &plan.PlanNode{
			NodeType: plan.NodeTypeSeqScan,
			Rows:     &plan.RowEstimate{Estimated: 100, Actual: 95},
		},
	}
	ep.CalculateStatistics()
	if ep.Statistics.TotalActualRows != 95 {
		t.Errorf("TotalActualRows = %d, want 95", ep.Statistics.TotalActualRows)
	}
}

// --- ExecutionPlan.FindBottlenecks ---

func TestFindBottlenecks_FullTableScan(t *testing.T) {
	ep := &plan.ExecutionPlan{
		RootNode: makeSeqScanNode("large_table", 5000),
	}
	bottlenecks := ep.FindBottlenecks()
	if len(bottlenecks) == 0 {
		t.Error("expected bottleneck for large full table scan")
	}
	found := false
	for _, b := range bottlenecks {
		if b.Issue == "Full table scan on large table" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'Full table scan on large table' bottleneck")
	}
}

func TestFindBottlenecks_SmallFullTableScan_NoBottleneck(t *testing.T) {
	ep := &plan.ExecutionPlan{
		RootNode: makeSeqScanNode("small_table", 50),
	}
	bottlenecks := ep.FindBottlenecks()
	for _, b := range bottlenecks {
		if b.Issue == "Full table scan on large table" {
			t.Error("should not flag small table as bottleneck")
		}
	}
}

func TestFindBottlenecks_ExpensiveNode(t *testing.T) {
	ep := &plan.ExecutionPlan{
		RootNode: &plan.PlanNode{
			NodeType: plan.NodeTypeHashJoin,
			Cost:     &plan.Cost{TotalCost: 9999},
		},
	}
	bottlenecks := ep.FindBottlenecks()
	found := false
	for _, b := range bottlenecks {
		if b.Issue == "Expensive operation detected" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'Expensive operation detected' bottleneck")
	}
}

func TestFindBottlenecks_NestedLoopHighCardinality(t *testing.T) {
	ep := &plan.ExecutionPlan{
		RootNode: &plan.PlanNode{
			NodeType: plan.NodeTypeNestedLoop,
			Rows:     &plan.RowEstimate{Estimated: 50000},
		},
	}
	bottlenecks := ep.FindBottlenecks()
	found := false
	for _, b := range bottlenecks {
		if b.Issue == "Nested loop join with high cardinality" {
			found = true
		}
	}
	if !found {
		t.Error("expected nested loop bottleneck")
	}
}

func TestFindBottlenecks_PoorEstimation(t *testing.T) {
	ep := &plan.ExecutionPlan{
		RootNode: &plan.PlanNode{
			NodeType: plan.NodeTypeSeqScan,
			Rows:     &plan.RowEstimate{Estimated: 10, Actual: 5000},
		},
	}
	bottlenecks := ep.FindBottlenecks()
	found := false
	for _, b := range bottlenecks {
		if b.Issue == "Poor cardinality estimation" {
			found = true
		}
	}
	if !found {
		t.Error("expected poor cardinality estimation bottleneck")
	}
}

func TestFindBottlenecks_NilRoot(t *testing.T) {
	ep := &plan.ExecutionPlan{}
	bottlenecks := ep.FindBottlenecks()
	if len(bottlenecks) != 0 {
		t.Errorf("expected no bottlenecks for nil root, got %d", len(bottlenecks))
	}
}

// --- PlanAnalyzer ---

func TestPlanAnalyzer_AnalyzePlan(t *testing.T) {
	ep := &plan.ExecutionPlan{
		RootNode: makeSeqScanNode("products", 10000),
	}
	pa := plan.NewPlanAnalyzer("postgresql")
	analysis := pa.AnalyzePlan(ep)
	if analysis == nil {
		t.Fatal("expected non-nil analysis")
	}
	if analysis.PerformanceScore < 0 || analysis.PerformanceScore > 100 {
		t.Errorf("PerformanceScore %f out of range [0,100]", analysis.PerformanceScore)
	}
	if len(analysis.Issues) == 0 {
		t.Error("expected at least one issue for large full table scan")
	}
}

func TestPlanAnalyzer_AnalyzePlan_GoodPlan(t *testing.T) {
	ep := &plan.ExecutionPlan{
		RootNode: makeIndexScanNode("users", 10),
	}
	pa := plan.NewPlanAnalyzer("mysql")
	analysis := pa.AnalyzePlan(ep)
	if analysis == nil {
		t.Fatal("expected non-nil analysis")
	}
	// Good plan should have high score
	if analysis.PerformanceScore < 50 {
		t.Errorf("expected high performance score for index scan, got %f", analysis.PerformanceScore)
	}
}

// --- ParseJSONPlan ---

func TestParseJSONPlan_PostgreSQL(t *testing.T) {
	jsonPlan := `[{"Plan": {"Node Type": "Seq Scan", "Relation Name": "users", "Startup Cost": 0.0, "Total Cost": 100.0, "Plan Rows": 1000, "Plan Width": 50}}]`
	ep, err := plan.ParseJSONPlan([]byte(jsonPlan), "postgresql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ep == nil {
		t.Fatal("expected non-nil plan")
	}
	if ep.RootNode == nil {
		t.Fatal("expected non-nil root node")
	}
}

func TestParseJSONPlan_MySQL(t *testing.T) {
	jsonPlan := `{"query_block": {"select_id": 1, "cost_info": {"query_cost": "50.00"}, "table": {"table_name": "orders", "access_type": "ALL", "rows_examined_per_scan": 500}}}`
	ep, err := plan.ParseJSONPlan([]byte(jsonPlan), "mysql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ep == nil {
		t.Fatal("expected non-nil plan")
	}
}

func TestParseJSONPlan_InvalidJSON(t *testing.T) {
	_, err := plan.ParseJSONPlan([]byte("{invalid"), "postgresql")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

// --- ParseJSONPlan with sqlserver (internally calls parseSQLServerXMLPlan) ---

func TestParseJSONPlan_SQLServer_XML(t *testing.T) {
	xmlPlan := `<?xml version="1.0"?>
<ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.5" Build="15.0.2000.5">
  <BatchSequence>
    <Batch>
      <Statements>
        <StmtSimple StatementText="SELECT * FROM users" StatementType="SELECT">
          <QueryPlan>
            <RelOp NodeId="0" PhysicalOp="Clustered Index Scan" EstimateRows="1000" EstimateTotalCost="1.5">
            </RelOp>
          </QueryPlan>
        </StmtSimple>
      </Statements>
    </Batch>
  </BatchSequence>
</ShowPlanXML>`
	ep, err := plan.ParseJSONPlan([]byte(xmlPlan), "sqlserver")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ep == nil {
		t.Fatal("expected non-nil plan")
	}
}

func TestParseJSONPlan_SQLServer_InvalidXML(t *testing.T) {
	_, err := plan.ParseJSONPlan([]byte("<invalid>"), "sqlserver")
	if err == nil {
		t.Error("expected error for invalid XML")
	}
}

// --- ParseJSONPlan with sqlite (internally calls parseSQLiteTextPlan) ---

func TestParseJSONPlan_SQLite(t *testing.T) {
	textPlan := "0|0|0|SCAN TABLE users\n1|0|0|USE TEMP B-TREE FOR ORDER BY\n"
	ep, err := plan.ParseJSONPlan([]byte(textPlan), "sqlite")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ep == nil {
		t.Fatal("expected non-nil plan")
	}
}

func TestParseJSONPlan_SQLite_Empty(t *testing.T) {
	ep, err := plan.ParseJSONPlan([]byte(""), "sqlite")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ep == nil {
		t.Fatal("expected non-nil plan")
	}
}

func TestParseJSONPlan_UnsupportedDialect(t *testing.T) {
	_, err := plan.ParseJSONPlan([]byte("{}"), "oracle")
	if err == nil {
		t.Error("expected error for unsupported dialect")
	}
}
