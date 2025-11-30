package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
	"github.com/Chahine-tech/sql-parser-go/pkg/plan"
)

func BenchmarkExplainSimpleSelect(b *testing.B) {
	sql := "EXPLAIN SELECT * FROM users WHERE id = 1"
	d := dialect.GetDialect("mysql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExplainAnalyze(b *testing.B) {
	sql := "EXPLAIN ANALYZE SELECT u.name, COUNT(o.id) FROM users u JOIN orders o ON u.id = o.user_id GROUP BY u.name"
	d := dialect.GetDialect("postgresql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExplainWithFormat(b *testing.B) {
	sql := "EXPLAIN FORMAT=JSON SELECT * FROM orders WHERE status = 'pending'"
	d := dialect.GetDialect("mysql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExplainPostgreSQLOptions(b *testing.B) {
	sql := "EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) SELECT * FROM users WHERE email LIKE '%@example.com'"
	d := dialect.GetDialect("postgresql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExplainComplexQuery(b *testing.B) {
	sql := `EXPLAIN ANALYZE
		SELECT u.id, u.name, COUNT(o.id) as order_count, SUM(o.total) as total_spent
		FROM users u
		LEFT JOIN orders o ON u.id = o.user_id
		WHERE u.created_at > '2023-01-01'
		  AND u.status = 'active'
		GROUP BY u.id, u.name
		HAVING COUNT(o.id) > 5
		ORDER BY total_spent DESC
		LIMIT 20`
	d := dialect.GetDialect("postgresql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := parser.NewWithDialect(context.Background(), sql, d)
		_, err := p.ParseStatement()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPlanAnalyzer(b *testing.B) {
	// Create a sample execution plan
	executionPlan := &plan.ExecutionPlan{
		Query:   "SELECT * FROM users WHERE id = 1",
		Dialect: "postgresql",
		RootNode: &plan.PlanNode{
			NodeType:  plan.NodeTypeIndexScan,
			Operation: "Index Scan",
			Table:     "users",
			Index:     "users_pkey",
			Cost: &plan.Cost{
				StartupCost: 0.29,
				TotalCost:   8.30,
			},
			Rows: &plan.RowEstimate{
				Estimated: 1,
			},
		},
		TotalCost:     8.30,
		EstimatedRows: 1,
	}

	analyzer := plan.NewPlanAnalyzer("postgresql")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = analyzer.AnalyzePlan(executionPlan)
	}
}

func BenchmarkPlanBottleneckDetection(b *testing.B) {
	// Create a plan with potential bottlenecks
	executionPlan := &plan.ExecutionPlan{
		Query:   "SELECT * FROM large_table WHERE status = 'active'",
		Dialect: "mysql",
		RootNode: &plan.PlanNode{
			NodeType:  plan.NodeTypeFullTableScan,
			Operation: "Full Table Scan",
			Table:     "large_table",
			Cost: &plan.Cost{
				TotalCost: 50000.0,
			},
			Rows: &plan.RowEstimate{
				Estimated: 1000000,
			},
			Children: []*plan.PlanNode{
				{
					NodeType:  plan.NodeTypeNestedLoop,
					Operation: "Nested Loop Join",
					Cost: &plan.Cost{
						TotalCost: 25000.0,
					},
					Rows: &plan.RowEstimate{
						Estimated: 50000,
					},
				},
			},
		},
		TotalCost:     75000.0,
		EstimatedRows: 1000000,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = executionPlan.FindBottlenecks()
	}
}

func BenchmarkPlanStatisticsCalculation(b *testing.B) {
	// Create a complex plan tree
	executionPlan := &plan.ExecutionPlan{
		Query:   "Complex multi-table join",
		Dialect: "postgresql",
		RootNode: &plan.PlanNode{
			NodeType:  plan.NodeTypeHashJoin,
			Operation: "Hash Join",
			Cost: &plan.Cost{
				TotalCost: 1500.0,
			},
			Rows: &plan.RowEstimate{
				Estimated: 10000,
			},
			Children: []*plan.PlanNode{
				{
					NodeType:  plan.NodeTypeSeqScan,
					Operation: "Seq Scan",
					Table:     "table1",
					Cost: &plan.Cost{
						TotalCost: 500.0,
					},
					Rows: &plan.RowEstimate{
						Estimated: 5000,
					},
				},
				{
					NodeType:  plan.NodeTypeIndexScan,
					Operation: "Index Scan",
					Table:     "table2",
					Cost: &plan.Cost{
						TotalCost: 800.0,
					},
					Rows: &plan.RowEstimate{
						Estimated: 8000,
					},
				},
			},
		},
		TotalCost:     1500.0,
		EstimatedRows: 10000,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		executionPlan.CalculateStatistics()
	}
}
