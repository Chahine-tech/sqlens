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
	switch tt {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case SELECT:
		return "SELECT"
	case FROM:
		return "FROM"
	case WHERE:
		return "WHERE"
	case JOIN:
		return "JOIN"
	case INNER:
		return "INNER"
	case LEFT:
		return "LEFT"
	case RIGHT:
		return "RIGHT"
	case FULL:
		return "FULL"
	case ON:
		return "ON"
	case GROUP:
		return "GROUP"
	case BY:
		return "BY"
	case ORDER:
		return "ORDER"
	case HAVING:
		return "HAVING"
	case AS:
		return "AS"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case NOT:
		return "NOT"
	case IN:
		return "IN"
	case EXISTS:
		return "EXISTS"
	case DISTINCT:
		return "DISTINCT"
	case TOP:
		return "TOP"
	case LIMIT:
		return "LIMIT"
	case OFFSET:
		return "OFFSET"
	case UNION:
		return "UNION"
	case ALL:
		return "ALL"
	case INSERT:
		return "INSERT"
	case INTO:
		return "INTO"
	case VALUES:
		return "VALUES"
	case UPDATE:
		return "UPDATE"
	case SET:
		return "SET"
	case DELETE:
		return "DELETE"
	case CREATE:
		return "CREATE"
	case DROP:
		return "DROP"
	case ALTER:
		return "ALTER"
	case TABLE:
		return "TABLE"
	case VIEW:
		return "VIEW"
	case MATERIALIZED:
		return "MATERIALIZED"
	case CHECK:
		return "CHECK"
	case OPTION:
		return "OPTION"
	case ASSIGN:
		return "ASSIGN"
	case EQ:
		return "EQ"
	case NOT_EQ:
		return "NOT_EQ"
	case LT:
		return "LT"
	case GT:
		return "GT"
	case LTE:
		return "LTE"
	case GTE:
		return "GTE"
	case LIKE:
		return "LIKE"
	case BETWEEN:
		return "BETWEEN"
	case IS:
		return "IS"
	case NULL:
		return "NULL"
	case WITH:
		return "WITH"
	case RECURSIVE:
		return "RECURSIVE"
	case OVER:
		return "OVER"
	case PARTITION:
		return "PARTITION"
	case ROWS:
		return "ROWS"
	case RANGE:
		return "RANGE"
	case UNBOUNDED:
		return "UNBOUNDED"
	case PRECEDING:
		return "PRECEDING"
	case FOLLOWING:
		return "FOLLOWING"
	case CURRENT:
		return "CURRENT"
	case ROW:
		return "ROW"
	case INTERSECT:
		return "INTERSECT"
	case EXCEPT:
		return "EXCEPT"
	case CASE:
		return "CASE"
	case WHEN:
		return "WHEN"
	case THEN:
		return "THEN"
	case ELSE:
		return "ELSE"
	case END:
		return "END"
	case COMMA:
		return "COMMA"
	case SEMICOLON:
		return "SEMICOLON"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case DOT:
		return "DOT"
	case ASTERISK:
		return "ASTERISK"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case SLASH:
		return "SLASH"
	case PERCENT:
		return "PERCENT"
	case BEGIN:
		return "BEGIN"
	case START:
		return "START"
	case COMMIT:
		return "COMMIT"
	case ROLLBACK:
		return "ROLLBACK"
	case SAVEPOINT:
		return "SAVEPOINT"
	case RELEASE:
		return "RELEASE"
	case WORK:
		return "WORK"
	case TRANSACTION:
		return "TRANSACTION"
	case EXPLAIN:
		return "EXPLAIN"
	case ANALYZE:
		return "ANALYZE"
	case FORMAT:
		return "FORMAT"
	case QUERY:
		return "QUERY"
	case PLAN:
		return "PLAN"
	case EXTENDED:
		return "EXTENDED"
	case PROCEDURE:
		return "PROCEDURE"
	case FUNCTION:
		return "FUNCTION"
	case RETURNS:
		return "RETURNS"
	case RETURN:
		return "RETURN"
	case DECLARE:
		return "DECLARE"
	case CURSOR:
		return "CURSOR"
	case OPEN:
		return "OPEN"
	case FETCH:
		return "FETCH"
	case CLOSE:
		return "CLOSE"
	case INOUT:
		return "INOUT"
	case OUT:
		return "OUT"
	case LANGUAGE:
		return "LANGUAGE"
	case PLPGSQL:
		return "PLPGSQL"
	case SQL:
		return "SQL"
	case REPLACE:
		return "REPLACE"
	case SECURITY:
		return "SECURITY"
	case DEFINER:
		return "DEFINER"
	case INVOKER:
		return "INVOKER"
	case DETERMINISTIC:
		return "DETERMINISTIC"
	case MODIFIES:
		return "MODIFIES"
	case READS:
		return "READS"
	case CONTAINS:
		return "CONTAINS"
	case NO:
		return "NO"
	case LOOP:
		return "LOOP"
	case WHILE:
		return "WHILE"
	case FOR:
		return "FOR"
	case REVERSE:
		return "REVERSE"
	case EXIT:
		return "EXIT"
	case CONTINUE:
		return "CONTINUE"
	case ITERATE:
		return "ITERATE"
	case LABEL:
		return "LABEL"
	case ELSEIF:
		return "ELSEIF"
	case ELSIF:
		return "ELSIF"
	case VARIADIC:
		return "VARIADIC"
	case TRIGGER:
		return "TRIGGER"
	case BEFORE:
		return "BEFORE"
	case AFTER:
		return "AFTER"
	case INSTEAD:
		return "INSTEAD"
	case OF:
		return "OF"
	case EACH:
		return "EACH"
	case NEW:
		return "NEW"
	case OLD:
		return "OLD"
	default:
		return "UNKNOWN"
	}
}
