package parser

import "fmt"

type Node interface {
	String() string
	Type() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Base node implementation
type BaseNode struct{}

func (bn *BaseNode) String() string {
	return ""
}

func (bn *BaseNode) Type() string {
	return "BaseNode"
}

// SELECT Statement
type SelectStatement struct {
	BaseNode
	Distinct bool
	Top      *TopClause
	Columns  []Expression
	From     *FromClause
	Joins    []*JoinClause
	Where    Expression
	GroupBy  []Expression
	Having   Expression
	OrderBy  []*OrderByClause
	Limit    *LimitClause
}

func (ss *SelectStatement) statementNode() {}
func (ss *SelectStatement) Type() string   { return "SelectStatement" }
func (ss *SelectStatement) String() string {
	return fmt.Sprintf("SELECT Statement with %d columns", len(ss.Columns))
}

// FROM Clause
type FromClause struct {
	BaseNode
	Tables []TableReference
}

func (fc *FromClause) Type() string   { return "FromClause" }
func (fc *FromClause) String() string { return "FROM Clause" }

// Table Reference
type TableReference struct {
	BaseNode
	Schema string
	Name   string
	Alias  string
}

func (tr *TableReference) expressionNode() {}
func (tr *TableReference) Type() string    { return "TableReference" }
func (tr *TableReference) String() string {
	if tr.Schema != "" {
		return fmt.Sprintf("%s.%s", tr.Schema, tr.Name)
	}
	return tr.Name
}

// JOIN Clause
type JoinClause struct {
	BaseNode
	JoinType  string // INNER, LEFT, RIGHT, FULL
	Table     TableReference
	Condition Expression
}

func (jc *JoinClause) Type() string   { return "JoinClause" }
func (jc *JoinClause) String() string { return fmt.Sprintf("%s JOIN", jc.JoinType) }

// Column Reference
type ColumnReference struct {
	BaseNode
	Table  string
	Column string
}

func (cr *ColumnReference) expressionNode() {}
func (cr *ColumnReference) Type() string    { return "ColumnReference" }
func (cr *ColumnReference) String() string {
	if cr.Table != "" {
		return fmt.Sprintf("%s.%s", cr.Table, cr.Column)
	}
	return cr.Column
}

// Literal Expression
type Literal struct {
	BaseNode
	Value interface{}
}

func (l *Literal) expressionNode() {}
func (l *Literal) Type() string    { return "Literal" }
func (l *Literal) String() string  { return fmt.Sprintf("%v", l.Value) }

// Binary Expression (for WHERE conditions, etc.)
type BinaryExpression struct {
	BaseNode
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) expressionNode() {}
func (be *BinaryExpression) Type() string    { return "BinaryExpression" }
func (be *BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", be.Left.String(), be.Operator, be.Right.String())
}

// Function Call
type FunctionCall struct {
	BaseNode
	Name      string
	Arguments []Expression
}

func (fc *FunctionCall) expressionNode() {}
func (fc *FunctionCall) Type() string    { return "FunctionCall" }
func (fc *FunctionCall) String() string  { return fmt.Sprintf("%s(...)", fc.Name) }

// SELECT * Expression
type StarExpression struct {
	BaseNode
	Table string // optional table qualifier
}

func (se *StarExpression) expressionNode() {}
func (se *StarExpression) Type() string    { return "StarExpression" }
func (se *StarExpression) String() string {
	if se.Table != "" {
		return fmt.Sprintf("%s.*", se.Table)
	}
	return "*"
}

// Aliased Expression (expression AS alias)
type AliasedExpression struct {
	BaseNode
	Expression Expression
	Alias      string
}

func (ae *AliasedExpression) expressionNode() {}
func (ae *AliasedExpression) Type() string    { return "AliasedExpression" }
func (ae *AliasedExpression) String() string {
	if ae.Alias != "" {
		return fmt.Sprintf("%s AS %s", ae.Expression.String(), ae.Alias)
	}
	return ae.Expression.String()
}

// ORDER BY Clause
type OrderByClause struct {
	BaseNode
	Expression Expression
	Direction  string // ASC, DESC
}

func (obc *OrderByClause) Type() string { return "OrderByClause" }
func (obc *OrderByClause) String() string {
	return fmt.Sprintf("ORDER BY %s %s", obc.Expression.String(), obc.Direction)
}

// TOP Clause (SQL Server specific)
type TopClause struct {
	BaseNode
	Count   int
	Percent bool
}

func (tc *TopClause) Type() string   { return "TopClause" }
func (tc *TopClause) String() string { return fmt.Sprintf("TOP %d", tc.Count) }

// LIMIT Clause
type LimitClause struct {
	BaseNode
	Count  int
	Offset int
}

func (lc *LimitClause) Type() string   { return "LimitClause" }
func (lc *LimitClause) String() string { return fmt.Sprintf("LIMIT %d", lc.Count) }

// INSERT Statement
type InsertStatement struct {
	BaseNode
	Table   TableReference
	Columns []string         // Optional column list
	Values  [][]Expression   // For INSERT ... VALUES
	Select  *SelectStatement // For INSERT ... SELECT
}

func (is *InsertStatement) statementNode() {}
func (is *InsertStatement) Type() string   { return "InsertStatement" }
func (is *InsertStatement) String() string {
	if is.Select != nil {
		return fmt.Sprintf("INSERT INTO %s SELECT", is.Table.Name)
	}
	return fmt.Sprintf("INSERT INTO %s (%d rows)", is.Table.Name, len(is.Values))
}

// UPDATE Statement
type UpdateStatement struct {
	BaseNode
	Table   TableReference
	Set     []*Assignment
	Where   Expression
	OrderBy []*OrderByClause // MySQL/SQLite support ORDER BY in UPDATE
	Limit   *LimitClause     // MySQL/SQLite support LIMIT in UPDATE
}

func (us *UpdateStatement) statementNode() {}
func (us *UpdateStatement) Type() string   { return "UpdateStatement" }
func (us *UpdateStatement) String() string {
	return fmt.Sprintf("UPDATE %s SET %d columns", us.Table.Name, len(us.Set))
}

// Assignment for UPDATE SET clause
type Assignment struct {
	BaseNode
	Column string
	Value  Expression
}

func (a *Assignment) Type() string   { return "Assignment" }
func (a *Assignment) String() string { return fmt.Sprintf("%s = %s", a.Column, a.Value.String()) }

// DELETE Statement
type DeleteStatement struct {
	BaseNode
	From    TableReference
	Where   Expression
	OrderBy []*OrderByClause // MySQL/SQLite support ORDER BY in DELETE
	Limit   *LimitClause     // MySQL/SQLite support LIMIT in DELETE
}

func (ds *DeleteStatement) statementNode() {}
func (ds *DeleteStatement) Type() string   { return "DeleteStatement" }
func (ds *DeleteStatement) String() string {
	return fmt.Sprintf("DELETE FROM %s", ds.From.Name)
}

// Unary Expression (NOT, etc.)
type UnaryExpression struct {
	BaseNode
	Operator string
	Operand  Expression
}

func (ue *UnaryExpression) expressionNode() {}
func (ue *UnaryExpression) Type() string    { return "UnaryExpression" }
func (ue *UnaryExpression) String() string {
	return fmt.Sprintf("%s %s", ue.Operator, ue.Operand.String())
}

// IN Expression
type InExpression struct {
	BaseNode
	Expression Expression
	Values     []Expression
	Not        bool
}

func (ie *InExpression) expressionNode() {}
func (ie *InExpression) Type() string    { return "InExpression" }
func (ie *InExpression) String() string {
	if ie.Not {
		return fmt.Sprintf("%s NOT IN (...)", ie.Expression.String())
	}
	return fmt.Sprintf("%s IN (...)", ie.Expression.String())
}

// EXISTS Expression
type ExistsExpression struct {
	BaseNode
	Subquery Statement
	Not      bool
}

func (ee *ExistsExpression) expressionNode() {}
func (ee *ExistsExpression) Type() string    { return "ExistsExpression" }
func (ee *ExistsExpression) String() string {
	if ee.Not {
		return "NOT EXISTS (...)"
	}
	return "EXISTS (...)"
}

// SubqueryExpression wraps a SelectStatement to make it usable as an Expression
type SubqueryExpression struct {
	BaseNode
	Query *SelectStatement
}

func (se *SubqueryExpression) expressionNode() {}
func (se *SubqueryExpression) Type() string    { return "SubqueryExpression" }
func (se *SubqueryExpression) String() string {
	return fmt.Sprintf("(%s)", se.Query.String())
}

// CTE (Common Table Expression) - WITH clause
type CommonTableExpression struct {
	BaseNode
	Name    string
	Columns []string // Optional column names
	Query   *SelectStatement
}

func (cte *CommonTableExpression) Type() string   { return "CommonTableExpression" }
func (cte *CommonTableExpression) String() string { return fmt.Sprintf("CTE: %s", cte.Name) }

// WITH Statement (contains CTEs and main query)
type WithStatement struct {
	BaseNode
	Recursive bool
	CTEs      []*CommonTableExpression
	Query     Statement // Main query (usually SelectStatement)
}

func (ws *WithStatement) statementNode() {}
func (ws *WithStatement) Type() string   { return "WithStatement" }
func (ws *WithStatement) String() string { return fmt.Sprintf("WITH %d CTEs", len(ws.CTEs)) }

// Window Function
type WindowFunction struct {
	BaseNode
	Function   *FunctionCall // The window function (ROW_NUMBER, RANK, etc.)
	OverClause *OverClause
}

func (wf *WindowFunction) expressionNode() {}
func (wf *WindowFunction) Type() string    { return "WindowFunction" }
func (wf *WindowFunction) String() string {
	return fmt.Sprintf("%s OVER (...)", wf.Function.Name)
}

// OVER Clause for Window Functions
type OverClause struct {
	BaseNode
	PartitionBy []Expression
	OrderBy     []*OrderByClause
	Frame       *WindowFrame
}

func (oc *OverClause) Type() string   { return "OverClause" }
func (oc *OverClause) String() string { return "OVER clause" }

// Window Frame (ROWS/RANGE BETWEEN ... AND ...)
type WindowFrame struct {
	BaseNode
	FrameType string // ROWS or RANGE
	Start     *FrameBound
	End       *FrameBound
}

func (wf *WindowFrame) Type() string   { return "WindowFrame" }
func (wf *WindowFrame) String() string { return fmt.Sprintf("%s frame", wf.FrameType) }

// Frame Boundary (UNBOUNDED PRECEDING, CURRENT ROW, etc.)
type FrameBound struct {
	BaseNode
	BoundType string     // UNBOUNDED, CURRENT, or expression
	Direction string     // PRECEDING or FOLLOWING
	Offset    Expression // For expression-based bounds
}

func (fb *FrameBound) Type() string   { return "FrameBound" }
func (fb *FrameBound) String() string { return fmt.Sprintf("%s %s", fb.BoundType, fb.Direction) }

// Set Operation (UNION, INTERSECT, EXCEPT)
type SetOperation struct {
	BaseNode
	Left     Statement
	Operator string // UNION, INTERSECT, EXCEPT
	All      bool   // UNION ALL, etc.
	Right    Statement
}

func (so *SetOperation) statementNode() {}
func (so *SetOperation) Type() string   { return "SetOperation" }
func (so *SetOperation) String() string {
	op := so.Operator
	if so.All {
		op += " ALL"
	}
	return fmt.Sprintf("Set Operation: %s", op)
}

// CASE Expression
type CaseExpression struct {
	BaseNode
	Input       Expression // Optional input for simple CASE
	WhenClauses []*WhenClause
	ElseResult  Expression // Optional ELSE clause
}

func (ce *CaseExpression) expressionNode() {}
func (ce *CaseExpression) Type() string    { return "CaseExpression" }
func (ce *CaseExpression) String() string {
	return fmt.Sprintf("CASE with %d WHEN clauses", len(ce.WhenClauses))
}

// WHEN Clause in CASE expression
type WhenClause struct {
	BaseNode
	Condition Expression
	Result    Expression
}

func (wc *WhenClause) Type() string   { return "WhenClause" }
func (wc *WhenClause) String() string { return "WHEN clause" }
