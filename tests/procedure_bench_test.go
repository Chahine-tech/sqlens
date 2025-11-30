package tests

import (
	"context"
	"testing"

	"github.com/Chahine-tech/sql-parser-go/pkg/dialect"
	"github.com/Chahine-tech/sql-parser-go/pkg/parser"
)

func BenchmarkSimpleProcedure(b *testing.B) {
	sql := `CREATE PROCEDURE get_user(IN user_id INT)
BEGIN
	RETURN;
END`
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

func BenchmarkProcedureWithMultipleParameters(b *testing.B) {
	sql := `CREATE PROCEDURE calculate_total(
		IN product_id INT,
		IN quantity INT,
		IN discount DECIMAL(5,2),
		OUT subtotal DECIMAL(10,2),
		OUT tax DECIMAL(10,2),
		OUT total DECIMAL(10,2)
	)
BEGIN
	RETURN;
END`
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

func BenchmarkSimpleFunction(b *testing.B) {
	sql := `CREATE FUNCTION calculate_tax(amount DECIMAL(10,2))
RETURNS DECIMAL(10,2)
DETERMINISTIC
BEGIN
	RETURN amount;
END`
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

func BenchmarkCreateOrReplaceProcedure(b *testing.B) {
	sql := `CREATE OR REPLACE PROCEDURE update_stats()
LANGUAGE plpgsql
AS
BEGIN
	RETURN;
END`
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

func BenchmarkProcedureWithVariables(b *testing.B) {
	sql := `CREATE PROCEDURE process_order(order_id INT)
BEGIN
	DECLARE total DECIMAL(10,2);
	DECLARE status VARCHAR(20);
	DECLARE count INT;
	RETURN;
END`
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

func BenchmarkProcedureWithCursor(b *testing.B) {
	sql := `CREATE PROCEDURE list_users()
BEGIN
	DECLARE user_cursor CURSOR FOR SELECT id, name FROM users;
	RETURN;
END`
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

func BenchmarkComplexFunction(b *testing.B) {
	sql := `CREATE FUNCTION calculate_discount(
		base_price DECIMAL(10,2),
		customer_tier VARCHAR(20),
		quantity INT,
		promo_code VARCHAR(50)
	)
RETURNS DECIMAL(10,2)
DETERMINISTIC
BEGIN
	DECLARE discount DECIMAL(5,2);
	DECLARE final_price DECIMAL(10,2);
	RETURN final_price;
END`
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

func BenchmarkDataTypeParsing(b *testing.B) {
	sql := `CREATE PROCEDURE test_types(
		p1 VARCHAR(255),
		p2 DECIMAL(10,2),
		p3 INT,
		p4 BIGINT,
		p5 TEXT,
		p6 DATETIME
	)
BEGIN
	RETURN;
END`
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
