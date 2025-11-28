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
	WITH      // CTE
	RECURSIVE // Recursive CTE
	OVER      // Window functions
	PARTITION // Window functions
	ROWS      // Window frame
	RANGE     // Window frame
	UNBOUNDED // Window frame
	PRECEDING // Window frame
	FOLLOWING // Window frame
	CURRENT   // Window frame
	ROW       // Window frame
	INTERSECT // Set operations
	EXCEPT    // Set operations
	CASE      // CASE expression
	WHEN      // CASE expression
	THEN      // CASE expression
	ELSE      // CASE expression
	END       // CASE/CTE end

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
	"SELECT":    SELECT,
	"FROM":      FROM,
	"WHERE":     WHERE,
	"JOIN":      JOIN,
	"INNER":     INNER,
	"LEFT":      LEFT,
	"RIGHT":     RIGHT,
	"FULL":      FULL,
	"ON":        ON,
	"GROUP":     GROUP,
	"BY":        BY,
	"ORDER":     ORDER,
	"HAVING":    HAVING,
	"AS":        AS,
	"AND":       AND,
	"OR":        OR,
	"NOT":       NOT,
	"IN":        IN,
	"EXISTS":    EXISTS,
	"DISTINCT":  DISTINCT,
	"TOP":       TOP,
	"LIMIT":     LIMIT,
	"OFFSET":    OFFSET,
	"UNION":     UNION,
	"ALL":       ALL,
	"INSERT":    INSERT,
	"INTO":      INTO,
	"VALUES":    VALUES,
	"UPDATE":    UPDATE,
	"SET":       SET,
	"DELETE":    DELETE,
	"CREATE":    CREATE,
	"DROP":      DROP,
	"ALTER":     ALTER,
	"TABLE":     TABLE,
	"LIKE":      LIKE,
	"BETWEEN":   BETWEEN,
	"IS":        IS,
	"NULL":      NULL,
	"WITH":      WITH,
	"RECURSIVE": RECURSIVE,
	"OVER":      OVER,
	"PARTITION": PARTITION,
	"ROWS":      ROWS,
	"RANGE":     RANGE,
	"UNBOUNDED": UNBOUNDED,
	"PRECEDING": PRECEDING,
	"FOLLOWING": FOLLOWING,
	"CURRENT":   CURRENT,
	"ROW":       ROW,
	"INTERSECT": INTERSECT,
	"EXCEPT":    EXCEPT,
	"CASE":      CASE,
	"WHEN":      WHEN,
	"THEN":      THEN,
	"ELSE":      ELSE,
	"END":       END,
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
	default:
		return "UNKNOWN"
	}
}
