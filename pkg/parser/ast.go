package parser

import (
	"fmt"
	"strings"
)

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
	Schema   string
	Name     string
	Alias    string
	Subquery *SelectStatement // For derived tables: (SELECT ...) AS alias
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
	Columns []string  // Optional column names
	Query   Statement // Can be SelectStatement or SetOperation (for recursive CTEs with UNION)
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

// DDL Statements

// CREATE TABLE Statement
type CreateTableStatement struct {
	BaseNode
	Table       TableReference
	Columns     []*ColumnDefinition
	Constraints []*TableConstraint
	IfNotExists bool
}

func (cts *CreateTableStatement) statementNode() {}
func (cts *CreateTableStatement) Type() string   { return "CreateTableStatement" }
func (cts *CreateTableStatement) String() string {
	return fmt.Sprintf("CREATE TABLE %s", cts.Table.Name)
}

// Column Definition
type ColumnDefinition struct {
	BaseNode
	Name          string
	DataType      string
	Length        int // For VARCHAR(255), etc.
	Precision     int // For DECIMAL(10,2)
	Scale         int // For DECIMAL(10,2)
	NotNull       bool
	PrimaryKey    bool
	Unique        bool
	AutoIncrement bool
	Default       Expression
	References    *ForeignKeyReference // For inline FOREIGN KEY
}

func (cd *ColumnDefinition) Type() string   { return "ColumnDefinition" }
func (cd *ColumnDefinition) String() string { return fmt.Sprintf("%s %s", cd.Name, cd.DataType) }

// Table Constraint (PRIMARY KEY, FOREIGN KEY, UNIQUE, CHECK)
type TableConstraint struct {
	BaseNode
	Name           string // Optional constraint name
	ConstraintType string // PRIMARY_KEY, FOREIGN_KEY, UNIQUE, CHECK
	Columns        []string
	References     *ForeignKeyReference // For FOREIGN KEY
	Check          Expression           // For CHECK constraint
}

func (tc *TableConstraint) Type() string   { return "TableConstraint" }
func (tc *TableConstraint) String() string { return fmt.Sprintf("%s constraint", tc.ConstraintType) }

// Foreign Key Reference
type ForeignKeyReference struct {
	BaseNode
	Table    string
	Columns  []string
	OnDelete string // CASCADE, SET NULL, etc.
	OnUpdate string
}

func (fkr *ForeignKeyReference) Type() string   { return "ForeignKeyReference" }
func (fkr *ForeignKeyReference) String() string { return fmt.Sprintf("REFERENCES %s", fkr.Table) }

// DROP Statement (TABLE, DATABASE, INDEX)
type DropStatement struct {
	BaseNode
	ObjectType string // TABLE, DATABASE, INDEX
	ObjectName string
	IfExists   bool
	Cascade    bool
}

func (ds *DropStatement) statementNode() {}
func (ds *DropStatement) Type() string   { return "DropStatement" }
func (ds *DropStatement) String() string {
	return fmt.Sprintf("DROP %s %s", ds.ObjectType, ds.ObjectName)
}

// ALTER TABLE Statement
type AlterTableStatement struct {
	BaseNode
	Table  TableReference
	Action *AlterAction
}

func (ats *AlterTableStatement) statementNode() {}
func (ats *AlterTableStatement) Type() string   { return "AlterTableStatement" }
func (ats *AlterTableStatement) String() string {
	return fmt.Sprintf("ALTER TABLE %s", ats.Table.Name)
}

// ALTER Action
type AlterAction struct {
	BaseNode
	ActionType string            // ADD, DROP, MODIFY, CHANGE, RENAME
	Column     *ColumnDefinition // For ADD/MODIFY
	ColumnName string            // For DROP/CHANGE
	NewColumn  *ColumnDefinition // For CHANGE
	Constraint *TableConstraint  // For ADD constraint
}

func (aa *AlterAction) Type() string   { return "AlterAction" }
func (aa *AlterAction) String() string { return fmt.Sprintf("ALTER %s", aa.ActionType) }

// CREATE INDEX Statement
type CreateIndexStatement struct {
	BaseNode
	IndexName   string
	Table       TableReference
	Columns     []string
	Unique      bool
	IfNotExists bool
}

func (cis *CreateIndexStatement) statementNode() {}
func (cis *CreateIndexStatement) Type() string   { return "CreateIndexStatement" }
func (cis *CreateIndexStatement) String() string {
	return fmt.Sprintf("CREATE INDEX %s ON %s", cis.IndexName, cis.Table.Name)
}

// Transaction Statements

// BEGIN/START TRANSACTION Statement
type BeginTransactionStatement struct {
	BaseNode
	UseStart bool // true if START TRANSACTION, false if BEGIN
}

func (bts *BeginTransactionStatement) statementNode() {}
func (bts *BeginTransactionStatement) Type() string   { return "BeginTransactionStatement" }
func (bts *BeginTransactionStatement) String() string {
	if bts.UseStart {
		return "START TRANSACTION"
	}
	return "BEGIN TRANSACTION"
}

// COMMIT Statement
type CommitStatement struct {
	BaseNode
	Work bool // true if COMMIT WORK
}

func (cs *CommitStatement) statementNode() {}
func (cs *CommitStatement) Type() string   { return "CommitStatement" }
func (cs *CommitStatement) String() string {
	if cs.Work {
		return "COMMIT WORK"
	}
	return "COMMIT"
}

// ROLLBACK Statement
type RollbackStatement struct {
	BaseNode
	Work        bool   // true if ROLLBACK WORK
	ToSavepoint string // Optional: ROLLBACK TO SAVEPOINT name
}

func (rs *RollbackStatement) statementNode() {}
func (rs *RollbackStatement) Type() string   { return "RollbackStatement" }
func (rs *RollbackStatement) String() string {
	if rs.ToSavepoint != "" {
		return fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", rs.ToSavepoint)
	}
	if rs.Work {
		return "ROLLBACK WORK"
	}
	return "ROLLBACK"
}

// SAVEPOINT Statement
type SavepointStatement struct {
	BaseNode
	Name string
}

func (ss *SavepointStatement) statementNode() {}
func (ss *SavepointStatement) Type() string   { return "SavepointStatement" }
func (ss *SavepointStatement) String() string {
	return fmt.Sprintf("SAVEPOINT %s", ss.Name)
}

// RELEASE SAVEPOINT Statement
type ReleaseSavepointStatement struct {
	BaseNode
	Name string
}

func (rss *ReleaseSavepointStatement) statementNode() {}
func (rss *ReleaseSavepointStatement) Type() string   { return "ReleaseSavepointStatement" }
func (rss *ReleaseSavepointStatement) String() string {
	return fmt.Sprintf("RELEASE SAVEPOINT %s", rss.Name)
}

// EXPLAIN Statement
type ExplainStatement struct {
	BaseNode
	Statement Statement         // The statement to explain
	Analyze   bool              // EXPLAIN ANALYZE
	Format    string            // FORMAT (JSON, XML, TEXT, etc.)
	Options   map[string]string // Dialect-specific options
}

func (es *ExplainStatement) statementNode() {}
func (es *ExplainStatement) Type() string   { return "ExplainStatement" }
func (es *ExplainStatement) String() string {
	if es.Analyze {
		return fmt.Sprintf("EXPLAIN ANALYZE %s", es.Statement.Type())
	}
	return fmt.Sprintf("EXPLAIN %s", es.Statement.Type())
}

// ============================================================================
// VIEW Statements
// ============================================================================

// CreateViewStatement represents CREATE VIEW or CREATE MATERIALIZED VIEW
type CreateViewStatement struct {
	BaseNode
	OrReplace    bool              // CREATE OR REPLACE VIEW
	Materialized bool              // MATERIALIZED VIEW (PostgreSQL)
	IfNotExists  bool              // IF NOT EXISTS
	ViewName     TableReference    // View name (can have schema)
	Columns      []string          // Optional column list
	SelectStmt   *SelectStatement  // The SELECT query
	WithCheck    bool              // WITH CHECK OPTION
	Options      map[string]string // Dialect-specific options (e.g., SECURITY DEFINER)
}

func (cvs *CreateViewStatement) statementNode() {}
func (cvs *CreateViewStatement) Type() string   { return "CreateViewStatement" }
func (cvs *CreateViewStatement) String() string {
	var result string
	if cvs.OrReplace {
		result = "CREATE OR REPLACE "
	} else {
		result = "CREATE "
	}

	if cvs.Materialized {
		result += "MATERIALIZED "
	}

	result += "VIEW "

	if cvs.IfNotExists {
		result += "IF NOT EXISTS "
	}

	result += cvs.ViewName.String()

	if len(cvs.Columns) > 0 {
		result += fmt.Sprintf(" (%s)", strings.Join(cvs.Columns, ", "))
	}

	result += " AS " + cvs.SelectStmt.String()

	if cvs.WithCheck {
		result += " WITH CHECK OPTION"
	}

	return result
}

// ============================================================================
// Stored Procedures and Functions
// ============================================================================

// DataTypeDefinition represents a data type with optional size/precision
type DataTypeDefinition struct {
	BaseNode
	Name      string // VARCHAR, INT, DECIMAL, etc.
	Length    int    // For VARCHAR(255), CHAR(10), etc.
	Precision int    // For DECIMAL(10,2), NUMERIC(8,3)
	Scale     int    // For DECIMAL(10,2), NUMERIC(8,3)
	IsArray   bool   // For array types (PostgreSQL)
}

func (dtd *DataTypeDefinition) Type() string { return "DataTypeDefinition" }
func (dtd *DataTypeDefinition) String() string {
	if dtd.Length > 0 {
		return fmt.Sprintf("%s(%d)", dtd.Name, dtd.Length)
	}
	if dtd.Precision > 0 {
		return fmt.Sprintf("%s(%d,%d)", dtd.Name, dtd.Precision, dtd.Scale)
	}
	if dtd.IsArray {
		return fmt.Sprintf("%s[]", dtd.Name)
	}
	return dtd.Name
}

// CreateProcedureStatement represents CREATE PROCEDURE
type CreateProcedureStatement struct {
	BaseNode
	Name            string
	Parameters      []*ProcedureParameter
	Body            *ProcedureBody
	Language        string            // SQL, PLPGSQL, etc.
	SecurityDefiner bool              // SECURITY DEFINER vs INVOKER
	Options         map[string]string // Dialect-specific options
	OrReplace       bool              // CREATE OR REPLACE
	IfNotExists     bool              // IF NOT EXISTS
}

func (cps *CreateProcedureStatement) statementNode() {}
func (cps *CreateProcedureStatement) Type() string   { return "CreateProcedureStatement" }
func (cps *CreateProcedureStatement) String() string {
	return fmt.Sprintf("CREATE PROCEDURE %s (%d parameters)", cps.Name, len(cps.Parameters))
}

// CreateFunctionStatement represents CREATE FUNCTION
type CreateFunctionStatement struct {
	BaseNode
	Name            string
	Parameters      []*ProcedureParameter
	ReturnType      *DataTypeDefinition
	Body            *ProcedureBody
	Language        string            // SQL, PLPGSQL, etc.
	Deterministic   bool              // DETERMINISTIC
	SecurityDefiner bool              // SECURITY DEFINER vs INVOKER
	Options         map[string]string // Dialect-specific options
	OrReplace       bool              // CREATE OR REPLACE
	IfNotExists     bool              // IF NOT EXISTS
}

func (cfs *CreateFunctionStatement) statementNode() {}
func (cfs *CreateFunctionStatement) Type() string   { return "CreateFunctionStatement" }
func (cfs *CreateFunctionStatement) String() string {
	return fmt.Sprintf("CREATE FUNCTION %s (%d parameters) RETURNS %s", cfs.Name, len(cfs.Parameters), cfs.ReturnType.Name)
}

// ProcedureParameter represents a parameter in a procedure/function
type ProcedureParameter struct {
	BaseNode
	Name       string
	Mode       string              // IN, OUT, INOUT
	DataType   *DataTypeDefinition // Parameter type
	Default    Expression          // Default value
	IsVariadic bool                // VARIADIC (PostgreSQL)
}

func (pp *ProcedureParameter) Type() string { return "ProcedureParameter" }
func (pp *ProcedureParameter) String() string {
	return fmt.Sprintf("%s %s %s", pp.Mode, pp.Name, pp.DataType.Name)
}

// ProcedureBody represents the body of a procedure/function
type ProcedureBody struct {
	BaseNode
	Statements []Statement     // List of statements in the body
	Variables  []*VariableDecl // DECLARE variables
	Cursors    []*CursorDecl   // DECLARE cursors
}

func (pb *ProcedureBody) Type() string { return "ProcedureBody" }
func (pb *ProcedureBody) String() string {
	return fmt.Sprintf("Procedure Body (%d statements)", len(pb.Statements))
}

// VariableDecl represents a variable declaration (DECLARE)
type VariableDecl struct {
	BaseNode
	Name     string
	DataType *DataTypeDefinition
	Default  Expression
}

func (vd *VariableDecl) statementNode() {}
func (vd *VariableDecl) Type() string   { return "VariableDecl" }
func (vd *VariableDecl) String() string {
	return fmt.Sprintf("DECLARE %s %s", vd.Name, vd.DataType.Name)
}

// CursorDecl represents a cursor declaration
type CursorDecl struct {
	BaseNode
	Name  string
	Query *SelectStatement
}

func (cd *CursorDecl) statementNode() {}
func (cd *CursorDecl) Type() string   { return "CursorDecl" }
func (cd *CursorDecl) String() string { return fmt.Sprintf("DECLARE CURSOR %s", cd.Name) }

// IfStatement represents IF...THEN...ELSE
type IfStatement struct {
	BaseNode
	Condition  Expression
	ThenBlock  []Statement
	ElseIfList []*ElseIfBlock
	ElseBlock  []Statement
}

func (is *IfStatement) statementNode() {}
func (is *IfStatement) Type() string   { return "IfStatement" }
func (is *IfStatement) String() string { return "IF Statement" }

// ElseIfBlock represents ELSEIF/ELSIF block
type ElseIfBlock struct {
	BaseNode
	Condition Expression
	Block     []Statement
}

func (eib *ElseIfBlock) Type() string   { return "ElseIfBlock" }
func (eib *ElseIfBlock) String() string { return "ELSEIF Block" }

// WhileStatement represents WHILE loop
type WhileStatement struct {
	BaseNode
	Condition Expression
	Block     []Statement
	Label     string // Loop label
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) Type() string   { return "WhileStatement" }
func (ws *WhileStatement) String() string { return "WHILE Loop" }

// LoopStatement represents LOOP...END LOOP
type LoopStatement struct {
	BaseNode
	Block []Statement
	Label string // Loop label
}

func (ls *LoopStatement) statementNode() {}
func (ls *LoopStatement) Type() string   { return "LoopStatement" }
func (ls *LoopStatement) String() string { return "LOOP" }

// ForStatement represents FOR loop
type ForStatement struct {
	BaseNode
	Variable  string
	Start     Expression
	End       Expression
	Step      Expression // Optional
	Block     []Statement
	Label     string // Loop label
	IsReverse bool   // FOR ... IN REVERSE (PostgreSQL)
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) Type() string   { return "ForStatement" }
func (fs *ForStatement) String() string { return "FOR Loop" }

// CaseStatement represents CASE statement (procedural, not expression)
type CaseStatement struct {
	BaseNode
	Expression Expression   // Optional: CASE expr WHEN...
	WhenList   []*WhenBlock // WHEN conditions
	ElseBlock  []Statement  // ELSE block
}

func (cs *CaseStatement) statementNode() {}
func (cs *CaseStatement) Type() string   { return "CaseStatement" }
func (cs *CaseStatement) String() string { return "CASE Statement" }

// WhenBlock represents a WHEN block in CASE
type WhenBlock struct {
	BaseNode
	Condition Expression
	Block     []Statement
}

func (wb *WhenBlock) Type() string   { return "WhenBlock" }
func (wb *WhenBlock) String() string { return "WHEN Block" }

// ReturnStatement represents RETURN
type ReturnStatement struct {
	BaseNode
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) Type() string   { return "ReturnStatement" }
func (rs *ReturnStatement) String() string { return "RETURN" }

// AssignmentStatement represents variable assignment (SET or :=)
type AssignmentStatement struct {
	BaseNode
	Variable string
	Value    Expression
}

func (as *AssignmentStatement) statementNode() {}
func (as *AssignmentStatement) Type() string   { return "AssignmentStatement" }
func (as *AssignmentStatement) String() string { return fmt.Sprintf("SET %s = ...", as.Variable) }

// OpenCursorStatement represents OPEN cursor
type OpenCursorStatement struct {
	BaseNode
	CursorName string
}

func (ocs *OpenCursorStatement) statementNode() {}
func (ocs *OpenCursorStatement) Type() string   { return "OpenCursorStatement" }
func (ocs *OpenCursorStatement) String() string { return fmt.Sprintf("OPEN %s", ocs.CursorName) }

// FetchStatement represents FETCH cursor
type FetchStatement struct {
	BaseNode
	CursorName string
	Variables  []string // INTO variables
}

func (fs *FetchStatement) statementNode() {}
func (fs *FetchStatement) Type() string   { return "FetchStatement" }
func (fs *FetchStatement) String() string { return fmt.Sprintf("FETCH %s", fs.CursorName) }

// CloseStatement represents CLOSE cursor
type CloseStatement struct {
	BaseNode
	CursorName string
}

func (cs *CloseStatement) statementNode() {}
func (cs *CloseStatement) Type() string   { return "CloseStatement" }
func (cs *CloseStatement) String() string { return fmt.Sprintf("CLOSE %s", cs.CursorName) }

// ExitStatement represents EXIT/BREAK (loop control)
type ExitStatement struct {
	BaseNode
	Label     string     // Loop label to exit
	Condition Expression // WHEN condition (PostgreSQL)
}

func (es *ExitStatement) statementNode() {}
func (es *ExitStatement) Type() string   { return "ExitStatement" }
func (es *ExitStatement) String() string { return "EXIT" }

// ContinueStatement represents CONTINUE/ITERATE (loop control)
type ContinueStatement struct {
	BaseNode
	Label     string     // Loop label
	Condition Expression // WHEN condition (PostgreSQL)
}

func (cs *ContinueStatement) statementNode() {}
func (cs *ContinueStatement) Type() string   { return "ContinueStatement" }
func (cs *ContinueStatement) String() string { return "CONTINUE" }
