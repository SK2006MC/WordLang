package token

// TokenType is a string representation of a token's type.
type TokenType string

// Token represents a token in our language.
type Token struct {
	Type    TokenType
	Literal string
	Line    int // For error reporting
	Column  int // For error reporting
}

// List of Token Types (Keywords and Symbols as Keywords in WordLang)
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + Literals
	IDENT  = "IDENT" // e.g., variable names, function names
	NUMBER = "NUMBER"
	STRING = "STRING"
	TRUE   = "TRUE"
	FALSE  = "FALSE"

	// Keywords
	LET      = "LET"
	FUNCTION = "FUNCTION"
	CALL     = "CALL"
	IF       = "IF"
	ELSE     = "ELSE"
	ELSEIF   = "ELSEIF"
	ENDIF    = "ENDIF"
	WHILE    = "WHILE"
	ENDWHILE = "ENDWHILE"
	FOREACH  = "FOREACH"
	IN       = "IN"
	ENDFOREACH = "ENDFOREACH"
	PRINT    = "PRINT"
	INPUT    = "INPUT"
	ADD      = "ADD"
	SUBTRACT = "SUBTRACT"
	MULTIPLY = "MULTIPLY"
	DIVIDE   = "DIVIDE"
	AND      = "AND"
	OR       = "OR"
	NOT      = "NOT"
	EQUALS   = "EQUALS"
	NOTEQUALS = "NOTEQUALS"
	GREATERTHAN = "GREATERTHAN"
	LESSTHAN    = "LESSTHAN"
	GREATEREQUAL = "GREATEREQUAL"
	LESSEQUAL    = "LESSEQUAL"
	THEN     = "THEN"
	DO       = "DO"
	END      = "END" // Generic 'end' keyword for blocks
	LIST     = "LIST"
	FROM       = "FROM"
	INDEX      = "INDEX"
	ISDEFINED  = "ISDEFINED"
	EXIT       = "EXIT"
	RETURN     = "RETURN"
	CONVERTTONUMBER = "CONVERTTONUMBER"
	CONVERTTOSTRING = "CONVERTTOSTRING"
	BE         = "BE"        // Add BE token type
	ENDFUNCTION = "ENDFUNCTION" // Add ENDFUNCTION token type


	// Punctuation (minimal, but we might keep # for comments)
	COMMENT = "COMMENT"
	HASH    = "#"
	SPACE   = "SPACE" // For handling whitespace significance later
	NEWLINE = "NEWLINE"
)

var keywords = map[string]TokenType{
	"let":               LET,
	"function":          FUNCTION,
	"call":              CALL,
	"if":                IF,
	"else":              ELSE,
	"elseif":            ELSEIF,
	"endif":             ENDIF,
	"while":             WHILE,
	"endwhile":          ENDWHILE,
	"foreach":           FOREACH,
	"in":                IN,
	"endforeach":        ENDFOREACH,
	"input":             INPUT,
	"add":               ADD,
	"sub":               SUBTRACT,
	"mult":              MULTIPLY,
	"div":               DIVIDE,
	"and":               AND,
	"or":                OR,
	"not":               NOT,
	"equals":            EQUALS,
	"notequals":         NOTEQUALS,
	"greater":           GREATERTHAN,
	"less":              LESSTHAN, // Shortened for brevity in keywords
	"greater or equal":  GREATEREQUAL,
	"less or equal":     LESSEQUAL, // Shortened for brevity
	"then":              THEN,
	"do":                DO,
	"end":               END,
	"list":              LIST,
	"get item at index": GETITEMATINDEX,
	"from":              FROM,
	"index":             INDEX,
	"isdefined":         ISDEFINED,
	"exit":              EXIT,
	"return":            RETURN,
	"convert to number": CONVERTTONUMBER,
	"convert to string": CONVERTTOSTRING,
	"true":              TRUE,
	"false":             FALSE,
	"be":                BE,        // Add "be" keyword
	"endfunction":       ENDFUNCTION, // Add "end function" keyword
}

// LookupIdent checks if the identifier is a keyword.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}