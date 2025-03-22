package lexer

import (
	"unicode"
	"wordlang/token"
)

// Lexer holds the state for lexing.
type Lexer struct {
	input        string
	position     int     // current position in input (points to current char)
	readPosition int     // next reading position in input (after current char)
	ch           byte    // current char under examination
	line         int     // current line number
	column       int     // current column number
	errors       []string // Lexer errors
}

// New creates a new Lexer.
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 1}
	l.readChar() // Initialize lexer
	return l
}

// readChar reads the next character and advances the position.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for "NUL" character, signals EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++ // Increment column on character read
}

// peekChar looks at the next character without advancing.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case 0:
		tok = newToken(token.EOF, l.ch)
		tok.Line = l.line
		tok.Column = l.column
	case '#':
		l.readComment()
		return l.NextToken() // Skip comment
	case '"':
		tok = l.readString()
	default:
		if unicode.IsLetter(rune(l.ch)) {
			ident := l.readIdentifier()
			// Check for multi-word keywords *immediately* after reading an identifier
			switch ident {
			case "greater":
				if l.peekKeyword("or") { // Check for "greater or"
					l.readIdentifier() // Consume "or"
					if l.peekKeyword("equal") { // Check for "greater or equal"
						l.readIdentifier() // Consume "equal"
						return token.Token{Type: token.GREATEREQUAL, Literal: "greater or equal", Line: l.line, Column: l.column - len("greater or equal") + 1}
					}
					return token.Token{Type: token.OR, Literal: "or", Line: l.line, Column: l.column - len("or") + 1} // Just "greater or" is treated as "or" keyword (might need refinement)
				} else if l.peekKeyword("than"){ // Check for "greater than"
					l.readIdentifier() // Consume "than"
					return token.Token{Type: token.GREATERTHAN, Literal: "greater than", Line: l.line, Column: l.column - len("greater than") + 1}
				}
				return token.Token{Type: token.GREATERTHAN, Literal: "greater", Line: l.line, Column: l.column - len("greater") + 1} // Just "greater" is treated as "greater than" keyword (might need refinement)
			case "less":
				if l.peekKeyword("or") { // Check for "less or"
					l.readIdentifier() // Consume "or"
					if l.peekKeyword("equal") { // Check for "less or equal"
						l.readIdentifier() // Consume "equal"
						return token.Token{Type: token.LESSEQUAL, Literal: "less or equal", Line: l.line, Column: l.column - len("less or equal") + 1}
					}
					return token.Token{Type: token.OR, Literal: "or", Line: l.line, Column: l.column - len("or") + 1} // Just "less or" is treated as "or" keyword (might need refinement)
				} else if l.peekKeyword("than"){ // Check for "less than"
					l.readIdentifier() // Consume "than"
					return token.Token{Type: token.LESSTHAN, Literal: "less than", Line: l.line, Column: l.column - len("less than") + 1}
				}
				return token.Token{Type: token.LESSTHAN, Literal: "less", Line: l.line, Column: l.column - len("less") + 1} // Just "less" is treated as "less than" keyword (might need refinement)
			case "end":
				if l.peekKeyword("if") {
					l.readIdentifier()
					return token.Token{Type: token.ENDIF, Literal: "endif", Line: l.line, Column: l.column - len("endif") + 1}
				} else if l.peekKeyword("while") {
					l.readIdentifier()
					return token.Token{Type: token.ENDWHILE, Literal: "endwhile", Line: l.line, Column: l.column - len("endwhile") + 1}
				} else if l.peekKeyword("foreach") {
					l.readIdentifier()
					return token.Token{Type: token.ENDFOREACH, Literal: "endforeach", Line: l.line, Column: l.column - len("endforeach") + 1}
				} else if l.peekKeyword("function") {
					l.readIdentifier()
					return token.Token{Type: token.ENDFUNCTION, Literal: "end function", Line: l.line, Column: l.column - len("end function") + 1}
				}
				return token.Token{Type: token.END, Literal: "end", Line: l.line, Column: l.column - len("end") + 1} // Just "end"
			case "get":
				if l.peekKeyword("item") {
					l.readIdentifier()
					if l.peekKeyword("at") {
						l.readIdentifier()
						if l.peekKeyword("index") {
							l.readIdentifier()
							return token.Token{Type: token.GETITEMATINDEX, Literal: "get item at index", Line: l.line, Column: l.column - len("get item at index") + 1}
						}
					}
				}
				return token.Token{Type: token.GETITEMATINDEX, Literal: "get", Line: l.line, Column: l.column - len("get") + 1} // Just "get" - might need refinement
			case "is":
				if l.peekKeyword("defined") {
					l.readIdentifier()
					return token.Token{Type: token.ISDEFINED, Literal: "is defined", Line: l.line, Column: l.column - len("is defined") + 1}
				}
				return token.Token{Type: token.ISDEFINED, Literal: "is", Line: l.line, Column: l.column - len("is") + 1} // Just "is" - might need refinement
			case "convert":
				if l.peekKeyword("to") {
					l.readIdentifier()
					if l.peekKeyword("number") {
						l.readIdentifier()
						return token.Token{Type: token.CONVERTTONUMBER, Literal: "convert to number", Line: l.line, Column: l.column - len("convert to number") + 1}
					} else if l.peekKeyword("string") {
						l.readIdentifier()
						return token.Token{Type: token.CONVERTTOSTRING, Literal: "convert to string", Line: l.line, Column: l.column - len("convert to string") + 1}
					}
				}
				return token.Token{Type: token.CONVERTTONUMBER, Literal: "convert", Line: l.line, Column: l.column - len("convert") + 1} // Just "convert" - might need refinement
			}


			tokType := token.LookupIdent(ident)
			return token.Token{Type: tokType, Literal: ident, Line: l.line, Column: l.column - len(ident) + 1}
		} else if unicode.IsDigit(rune(l.ch)) {
			return l.readNumber()
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
			tok.Line = l.line
			tok.Column = l.column
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) peekKeyword(keyword string) bool {
	currentPos := l.position
	currentReadPos := l.readPosition
	currentColumn := l.column
	currentChar := l.ch

	l.skipWhitespace() // Skip any whitespace before the potential keyword

	startPos := l.position
	for unicode.IsLetter(rune(l.ch)) {
		l.readChar()
	}
	peekedWord := l.input[startPos:l.position]

	l.position = currentPos
	l.readPosition = currentReadPos
	l.column = currentColumn
	l.ch = currentChar // Restore lexer state

	return peekedWord == keyword
}


func (l *Lexer) readIdentifier() string {
	startPos := l.position
	for unicode.IsLetter(rune(l.ch)) || unicode.IsDigit(rune(l.ch)) || l.ch == '_' { // Removed space from identifier chars
		l.readChar()
	}
	return l.input[startPos:l.position]
}


func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.ch)) {
		if l.ch == '\n' {
			l.line++
			l.column = 0 // Reset column on newline
		}
		l.readChar()
	}
}

// func (l *Lexer) readIdentifier() token.Token {
	// startPos := l.position
	// for unicode.IsLetter(rune(l.ch)) || unicode.IsDigit(rune(l.ch)) || l.ch == '_' || unicode.IsSpace(rune(l.ch)){ // Allow spaces in multi-word keywords
		// l.readChar()
	// }
	// literal := l.input[startPos:l.position]
	// tokType := token.LookupIdent(literal)
	// return token.Token{Type: tokType, Literal: literal, Line: l.line, Column: l.column - len(literal) + 1}
// }

func (l *Lexer) readNumber() token.Token {
    startPos := l.position
    for unicode.IsDigit(rune(l.ch)) || l.ch == '.' {
        l.readChar()
    }
    return token.Token{Type: token.NUMBER, Literal: l.input[startPos:l.position], Line: l.line, Column: l.column - len(l.input[startPos:l.position]) + 1}
}

func (l *Lexer) readString() token.Token {
	startPos := l.position + 1 // Skip the opening quote
	l.readChar() // Move past the opening quote
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	literal := l.input[startPos:l.position]
	return token.Token{Type: token.STRING, Literal: literal, Line: l.line, Column: l.column - len(literal) -1 } // Adjust column to start of string content
}

func (l *Lexer) readComment() token.Token {
	startPos := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	literal := l.input[startPos:l.position]
	return token.Token{Type: token.COMMENT, Literal: literal, Line: l.line, Column: l.column - len(literal) + 1}
}


func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// Errors returns the list of lexer errors.
func (l *Lexer) Errors() []string {
	return l.errors
}
