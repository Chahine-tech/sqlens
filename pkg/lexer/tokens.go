package lexer

import "fmt"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	// Identifiers and literals
	IDENT  // table_name, column_name
	STRING // 'hello'
	NUMBER // 123, 123.45

	// SQL Keywords
	SELECT
	FROM
	WHERE
	JOIN
	INNER
	LEFT
	RIGHT
	FULL
	ON
	GROUP
	BY
	ORDER
	HAVING
	AS
	AND
	OR
	NOT
	IN
	EXISTS
	DISTINCT
	TOP
	LIMIT
	OFFSET
	UNION
	ALL
	INSERT
	INTO   // INSERT INTO
	VALUES // INSERT VALUES
	UPDATE
	SET // UPDATE SET
	DELETE
	CREATE
	DROP
	ALTER
	TABLE
	VIEW         // CREATE VIEW
	MATERIALIZED // MATERIALIZED VIEW (PostgreSQL)
	WITH         // CTE or WITH CHECK OPTION
	RECURSIVE    // Recursive CTE
	CHECK        // WITH CHECK OPTION
	OPTION       // WITH CHECK OPTION
	OVER         // Window functions
	PARTITION    // Window functions
	ROWS         // Window frame
	RANGE        // Window frame
	UNBOUNDED    // Window frame
	PRECEDING    // Window frame
	FOLLOWING    // Window frame
	CURRENT      // Window frame
	ROW          // Window frame
	INTERSECT    // Set operations
	EXCEPT       // Set operations
	CASE         // CASE expression
	WHEN         // CASE expression
	THEN         // CASE expression
	ELSE         // CASE expression
	END          // CASE/CTE end

	// DDL Keywords
	PRIMARY        // PRIMARY KEY
	FOREIGN        // FOREIGN KEY
	KEY            // KEY
	CONSTRAINT     // CONSTRAINT
	UNIQUE         // UNIQUE
	INDEX          // INDEX
	AUTO_INCREMENT // AUTO_INCREMENT (MySQL)
	AUTOINCREMENT  // AUTOINCREMENT (SQLite)
	IDENTITY       // IDENTITY (SQL Server)
	DEFAULT        // DEFAULT value
	REFERENCES     // FOREIGN KEY REFERENCES
	ADD            // ALTER TABLE ADD
	MODIFY         // ALTER TABLE MODIFY
	CHANGE         // ALTER TABLE CHANGE (MySQL)
	COLUMN         // COLUMN
	IF             // IF EXISTS/IF NOT EXISTS
	DATABASE       // DATABASE
	SCHEMA         // SCHEMA

	// Transaction Keywords
	BEGIN       // BEGIN TRANSACTION
	START       // START TRANSACTION
	COMMIT      // COMMIT
	ROLLBACK    // ROLLBACK
	SAVEPOINT   // SAVEPOINT
	RELEASE     // RELEASE SAVEPOINT
	WORK        // WORK (optional in COMMIT/ROLLBACK)
	TRANSACTION // TRANSACTION

	// Execution Plan Keywords
	EXPLAIN  // EXPLAIN
	ANALYZE  // ANALYZE (EXPLAIN ANALYZE)
	FORMAT   // FORMAT (EXPLAIN FORMAT=JSON)
	QUERY    // QUERY (EXPLAIN QUERY PLAN - SQLite)
	PLAN     // PLAN (EXPLAIN QUERY PLAN - SQLite)
	EXTENDED // EXTENDED (MySQL)

	// Stored Procedures and Functions Keywords
	PROCEDURE     // CREATE PROCEDURE
	FUNCTION      // CREATE FUNCTION
	RETURNS       // RETURNS (function return type)
	RETURN        // RETURN statement
	DECLARE       // DECLARE variables/cursors
	CURSOR        // CURSOR
	OPEN          // OPEN cursor
	FETCH         // FETCH cursor
	CLOSE         // CLOSE cursor
	INOUT         // INOUT parameter mode
	OUT           // OUT parameter mode
	LANGUAGE      // LANGUAGE (PostgreSQL)
	PLPGSQL       // PL/pgSQL (PostgreSQL)
	SQL           // SQL language
	REPLACE       // OR REPLACE
	SECURITY      // SECURITY DEFINER/INVOKER
	DEFINER       // DEFINER
	INVOKER       // INVOKER
	DETERMINISTIC // DETERMINISTIC function
	MODIFIES      // MODIFIES SQL DATA
	READS         // READS SQL DATA
	CONTAINS      // CONTAINS SQL
	NO            // NO SQL
	LOOP          // LOOP
	WHILE         // WHILE loop
	FOR           // FOR loop
	REVERSE       // REVERSE (FOR REVERSE)
	EXIT          // EXIT loop
	CONTINUE      // CONTINUE loop
	ITERATE       // ITERATE (MySQL synonym for CONTINUE)
	LABEL         // Loop label
	ELSEIF        // ELSEIF (MySQL)
	ELSIF         // ELSIF (PostgreSQL, Oracle)
	VARIADIC      // VARIADIC parameters (PostgreSQL)

	// Trigger Keywords
	TRIGGER // CREATE TRIGGER
	BEFORE  // BEFORE INSERT/UPDATE/DELETE
	AFTER   // AFTER INSERT/UPDATE/DELETE
	INSTEAD // INSTEAD OF (SQL Server, Oracle)
	OF      // OF (INSTEAD OF)
	EACH    // FOR EACH ROW/STATEMENT
	NEW     // NEW (trigger references)
	OLD     // OLD (trigger references)

	// Operators
	ASSIGN  // =
	EQ      // ==
	NOT_EQ  // !=
	LT      // <
	GT      // >
	LTE     // <=
	GTE     // >=
	LIKE    // LIKE
	BETWEEN // BETWEEN
	IS      // IS
	NULL    // NULL

	// Delimiters
	COMMA     // ,
	SEMICOLON // ;
	LPAREN    // (
	RPAREN    // )
	DOT       // .
	ASTERISK  // *
	PLUS      // +
	MINUS     // -
	SLASH     // /
	PERCENT   // %
)

var keywords = map[string]TokenType{
	"SELECT":         SELECT,
	"FROM":           FROM,
	"WHERE":          WHERE,
	"JOIN":           JOIN,
	"INNER":          INNER,
	"LEFT":           LEFT,
	"RIGHT":          RIGHT,
	"FULL":           FULL,
	"ON":             ON,
	"GROUP":          GROUP,
	"BY":             BY,
	"ORDER":          ORDER,
	"HAVING":         HAVING,
	"AS":             AS,
	"AND":            AND,
	"OR":             OR,
	"NOT":            NOT,
	"IN":             IN,
	"EXISTS":         EXISTS,
	"DISTINCT":       DISTINCT,
	"TOP":            TOP,
	"LIMIT":          LIMIT,
	"OFFSET":         OFFSET,
	"UNION":          UNION,
	"ALL":            ALL,
	"INSERT":         INSERT,
	"INTO":           INTO,
	"VALUES":         VALUES,
	"UPDATE":         UPDATE,
	"SET":            SET,
	"DELETE":         DELETE,
	"CREATE":         CREATE,
	"DROP":           DROP,
	"ALTER":          ALTER,
	"TABLE":          TABLE,
	"VIEW":           VIEW,
	"MATERIALIZED":   MATERIALIZED,
	"CHECK":          CHECK,
	"OPTION":         OPTION,
	"LIKE":           LIKE,
	"BETWEEN":        BETWEEN,
	"IS":             IS,
	"NULL":           NULL,
	"WITH":           WITH,
	"RECURSIVE":      RECURSIVE,
	"OVER":           OVER,
	"PARTITION":      PARTITION,
	"ROWS":           ROWS,
	"RANGE":          RANGE,
	"UNBOUNDED":      UNBOUNDED,
	"PRECEDING":      PRECEDING,
	"FOLLOWING":      FOLLOWING,
	"CURRENT":        CURRENT,
	"ROW":            ROW,
	"INTERSECT":      INTERSECT,
	"EXCEPT":         EXCEPT,
	"CASE":           CASE,
	"WHEN":           WHEN,
	"THEN":           THEN,
	"ELSE":           ELSE,
	"END":            END,
	"PRIMARY":        PRIMARY,
	"FOREIGN":        FOREIGN,
	"KEY":            KEY,
	"CONSTRAINT":     CONSTRAINT,
	"UNIQUE":         UNIQUE,
	"INDEX":          INDEX,
	"AUTO_INCREMENT": AUTO_INCREMENT,
	"AUTOINCREMENT":  AUTOINCREMENT,
	"IDENTITY":       IDENTITY,
	"DEFAULT":        DEFAULT,
	"REFERENCES":     REFERENCES,
	"ADD":            ADD,
	"MODIFY":         MODIFY,
	"CHANGE":         CHANGE,
	"COLUMN":         COLUMN,
	"IF":             IF,
	"DATABASE":       DATABASE,
	"SCHEMA":         SCHEMA,
	"BEGIN":          BEGIN,
	"START":          START,
	"COMMIT":         COMMIT,
	"ROLLBACK":       ROLLBACK,
	"SAVEPOINT":      SAVEPOINT,
	"RELEASE":        RELEASE,
	"WORK":           WORK,
	"TRANSACTION":    TRANSACTION,
	"EXPLAIN":        EXPLAIN,
	"ANALYZE":        ANALYZE,
	"FORMAT":         FORMAT,
	"QUERY":          QUERY,
	"PLAN":           PLAN,
	"EXTENDED":       EXTENDED,
	"PROCEDURE":      PROCEDURE,
	"FUNCTION":       FUNCTION,
	"RETURNS":        RETURNS,
	"RETURN":         RETURN,
	"DECLARE":        DECLARE,
	"CURSOR":         CURSOR,
	"OPEN":           OPEN,
	"FETCH":          FETCH,
	"CLOSE":          CLOSE,
	"INOUT":          INOUT,
	"OUT":            OUT,
	"LANGUAGE":       LANGUAGE,
	"PLPGSQL":        PLPGSQL,
	"SQL":            SQL,
	"REPLACE":        REPLACE,
	"SECURITY":       SECURITY,
	"DEFINER":        DEFINER,
	"INVOKER":        INVOKER,
	"DETERMINISTIC":  DETERMINISTIC,
	"MODIFIES":       MODIFIES,
	"READS":          READS,
	"CONTAINS":       CONTAINS,
	"NO":             NO,
	"LOOP":           LOOP,
	"WHILE":          WHILE,
	"FOR":            FOR,
	"REVERSE":        REVERSE,
	"EXIT":           EXIT,
	"CONTINUE":       CONTINUE,
	"ITERATE":        ITERATE,
	"LABEL":          LABEL,
	"ELSEIF":         ELSEIF,
	"ELSIF":          ELSIF,
	"VARIADIC":       VARIADIC,
	"TRIGGER":        TRIGGER,
	"BEFORE":         BEFORE,
	"AFTER":          AFTER,
	"INSTEAD":        INSTEAD,
	"OF":             OF,
	"EACH":           EACH,
	"NEW":            NEW,
	"OLD":            OLD,
}

// tokenNames maps TokenType to string representation
var tokenNames = map[TokenType]string{
	ILLEGAL:         "ILLEGAL",
	EOF:             "EOF",
	IDENT:           "IDENT",
	STRING:          "STRING",
	NUMBER:          "NUMBER",
	SELECT:          "SELECT",
	FROM:            "FROM",
	WHERE:           "WHERE",
	JOIN:            "JOIN",
	INNER:           "INNER",
	LEFT:            "LEFT",
	RIGHT:           "RIGHT",
	FULL:            "FULL",
	ON:              "ON",
	GROUP:           "GROUP",
	BY:              "BY",
	ORDER:           "ORDER",
	HAVING:          "HAVING",
	AS:              "AS",
	AND:             "AND",
	OR:              "OR",
	NOT:             "NOT",
	IN:              "IN",
	EXISTS:          "EXISTS",
	DISTINCT:        "DISTINCT",
	TOP:             "TOP",
	LIMIT:           "LIMIT",
	OFFSET:          "OFFSET",
	UNION:           "UNION",
	ALL:             "ALL",
	INSERT:          "INSERT",
	INTO:            "INTO",
	VALUES:          "VALUES",
	UPDATE:          "UPDATE",
	SET:             "SET",
	DELETE:          "DELETE",
	CREATE:          "CREATE",
	DROP:            "DROP",
	ALTER:           "ALTER",
	TABLE:           "TABLE",
	VIEW:            "VIEW",
	MATERIALIZED:    "MATERIALIZED",
	CHECK:           "CHECK",
	OPTION:          "OPTION",
	ASSIGN:          "ASSIGN",
	EQ:              "EQ",
	NOT_EQ:          "NOT_EQ",
	LT:              "LT",
	GT:              "GT",
	LTE:             "LTE",
	GTE:             "GTE",
	LIKE:            "LIKE",
	BETWEEN:         "BETWEEN",
	IS:              "IS",
	NULL:            "NULL",
	WITH:            "WITH",
	RECURSIVE:       "RECURSIVE",
	OVER:            "OVER",
	PARTITION:       "PARTITION",
	ROWS:            "ROWS",
	RANGE:           "RANGE",
	UNBOUNDED:       "UNBOUNDED",
	PRECEDING:       "PRECEDING",
	FOLLOWING:       "FOLLOWING",
	CURRENT:         "CURRENT",
	ROW:             "ROW",
	INTERSECT:       "INTERSECT",
	EXCEPT:          "EXCEPT",
	CASE:            "CASE",
	WHEN:            "WHEN",
	THEN:            "THEN",
	ELSE:            "ELSE",
	END:             "END",
	COMMA:           "COMMA",
	SEMICOLON:       "SEMICOLON",
	LPAREN:          "LPAREN",
	RPAREN:          "RPAREN",
	DOT:             "DOT",
	ASTERISK:        "ASTERISK",
	PLUS:            "PLUS",
	MINUS:           "MINUS",
	SLASH:           "SLASH",
	PERCENT:         "PERCENT",
	BEGIN:           "BEGIN",
	START:           "START",
	COMMIT:          "COMMIT",
	ROLLBACK:        "ROLLBACK",
	SAVEPOINT:       "SAVEPOINT",
	RELEASE:         "RELEASE",
	WORK:            "WORK",
	TRANSACTION:     "TRANSACTION",
	EXPLAIN:         "EXPLAIN",
	ANALYZE:         "ANALYZE",
	FORMAT:          "FORMAT",
	QUERY:           "QUERY",
	PLAN:            "PLAN",
	EXTENDED:        "EXTENDED",
	PROCEDURE:       "PROCEDURE",
	FUNCTION:        "FUNCTION",
	RETURNS:         "RETURNS",
	RETURN:          "RETURN",
	DECLARE:         "DECLARE",
	CURSOR:          "CURSOR",
	OPEN:            "OPEN",
	FETCH:           "FETCH",
	CLOSE:           "CLOSE",
	INOUT:           "INOUT",
	OUT:             "OUT",
	LANGUAGE:        "LANGUAGE",
	PLPGSQL:         "PLPGSQL",
	SQL:             "SQL",
	REPLACE:         "REPLACE",
	SECURITY:        "SECURITY",
	DEFINER:         "DEFINER",
	INVOKER:         "INVOKER",
	DETERMINISTIC:   "DETERMINISTIC",
	MODIFIES:        "MODIFIES",
	READS:           "READS",
	CONTAINS:        "CONTAINS",
	NO:              "NO",
	LOOP:            "LOOP",
	WHILE:           "WHILE",
	FOR:             "FOR",
	REVERSE:         "REVERSE",
	EXIT:            "EXIT",
	CONTINUE:        "CONTINUE",
	ITERATE:         "ITERATE",
	LABEL:           "LABEL",
	ELSEIF:          "ELSEIF",
	ELSIF:           "ELSIF",
	VARIADIC:        "VARIADIC",
	TRIGGER:         "TRIGGER",
	BEFORE:          "BEFORE",
	AFTER:           "AFTER",
	INSTEAD:         "INSTEAD",
	OF:              "OF",
	EACH:            "EACH",
	NEW:             "NEW",
	OLD:             "OLD",
	PRIMARY:         "PRIMARY",
	FOREIGN:         "FOREIGN",
	KEY:             "KEY",
	CONSTRAINT:      "CONSTRAINT",
	UNIQUE:          "UNIQUE",
	INDEX:           "INDEX",
	AUTO_INCREMENT:  "AUTO_INCREMENT",
	AUTOINCREMENT:   "AUTOINCREMENT",
	IDENTITY:        "IDENTITY",
	DEFAULT:         "DEFAULT",
	REFERENCES:      "REFERENCES",
	ADD:             "ADD",
	MODIFY:          "MODIFY",
	CHANGE:          "CHANGE",
	COLUMN:          "COLUMN",
	IF:              "IF",
	DATABASE:        "DATABASE",
	SCHEMA:          "SCHEMA",
}

type Token struct {
	Type     TokenType
	Literal  string
	Position int
	Line     int
	Column   int
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %d, Literal: %s, Position: %d, Line: %d, Column: %d}",
		t.Type, t.Literal, t.Position, t.Line, t.Column)
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// TokenTypeString returns the string representation of a token type
func (tt TokenType) String() string {
	if name, ok := tokenNames[tt]; ok {
		return name
	}
	return "UNKNOWN"
}
