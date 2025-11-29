package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

// Benchmark BEGIN TRANSACTION
func BenchmarkBeginTransaction(b *testing.B) {
	sql := `BEGIN TRANSACTION`
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

// Benchmark START TRANSACTION
func BenchmarkStartTransaction(b *testing.B) {
	sql := `START TRANSACTION`
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

// Benchmark COMMIT
func BenchmarkCommit(b *testing.B) {
	sql := `COMMIT`
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

// Benchmark COMMIT WORK
func BenchmarkCommitWork(b *testing.B) {
	sql := `COMMIT WORK`
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

// Benchmark ROLLBACK
func BenchmarkRollback(b *testing.B) {
	sql := `ROLLBACK`
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

// Benchmark ROLLBACK TO SAVEPOINT
func BenchmarkRollbackToSavepoint(b *testing.B) {
	sql := `ROLLBACK TO SAVEPOINT sp1`
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

// Benchmark SAVEPOINT
func BenchmarkSavepoint(b *testing.B) {
	sql := `SAVEPOINT my_savepoint`
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

// Benchmark RELEASE SAVEPOINT
func BenchmarkReleaseSavepoint(b *testing.B) {
	sql := `RELEASE SAVEPOINT my_savepoint`
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
