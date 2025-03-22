package ast

import (
	"wordlang/token"
	"strings"
)

// Node is the base interface for all nodes in the AST.
type Node interface {
	TokenLiteral() string // For debugging and testing
	String() string       // For pretty printing the AST
}

// Statement is the interface for all statement nodes.
type Statement interface {
	Node
	statementNode()
}

// Expression is the interface for all expression nodes.
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of the AST.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out string
	for _, s := range p.Statements {
		out += s.String()
	}
	return out
}
// Identifier represents an identifier (variable name, function name).
type Identifier struct {
	Token token.Token // The identifier token
	Value string
}
func (i *Identifier) expressionNode()    {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// LetStatement represents a 'let' statement.
type LetStatement struct {
	Token token.Token // The 'let' token
	Name  *Identifier
	Value Expression
}
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	return ls.TokenLiteral() + " " + ls.Name.String() + " be " + ls.Value.String()
}


// ReturnStatement represents a 'return' statement.
type ReturnStatement struct {
	Token token.Token // The 'return' token
	ReturnValue Expression
}
func (rs *ReturnStatement) statementNode()     {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	return rs.TokenLiteral() + " " + rs.ReturnValue.String()
}

// ExpressionStatement wraps an expression to be used as a statement.
type ExpressionStatement struct {
	Token      token.Token // The first token of the expression
	Expression Expression
}
func (es *ExpressionStatement) statementNode()     {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IntegerLiteral represents an integer literal.
type IntegerLiteral struct {
	Token token.Token // The number token
	Value int64
}

func (il *IntegerLiteral) expressionNode()    {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral represents a floating-point literal.
type FloatLiteral struct {
	Token token.Token // The number token
	Value float64
}

func (fl *FloatLiteral) expressionNode()    {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// StringLiteral represents a string literal.
type StringLiteral struct {
	Token token.Token // The string token
	Value string
}

func (sl *StringLiteral) expressionNode()    {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// BooleanLiteral represents a boolean literal (true or false).
type BooleanLiteral struct {
	Token token.Token // The boolean token (TRUE or FALSE)
	Value bool
}

func (bl *BooleanLiteral) expressionNode()    {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

// PrefixExpression represents a prefix operator expression (e.g., 'not condition').
type PrefixExpression struct {
	Token    token.Token // The prefix operator token (e.g., NOT)
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()    {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}

// InfixExpression represents an infix operator expression (e.g., 'add a and b').
type InfixExpression struct {
	Token    token.Token // The operator token (e.g., ADD, EQUALS)
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()    {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	return "(" + oe.Left.String() + " " + oe.Operator + " " + oe.Right.String() + ")"
}

// IfStatement represents an 'if' statement.
type IfStatement struct {
	Token token.Token // The 'if' token
	Condition Expression
	ThenBlock   *BlockStatement
	ElseIfBlocks []*ElseIfBlock // Slice to handle multiple 'elseif'
	ElseBlock   *BlockStatement
}

func (is *IfStatement) statementNode()     {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	var out string
	out += "if " + is.Condition.String() + " then " + is.ThenBlock.String()
	for _, elseifBlock := range is.ElseIfBlocks {
		out += " elseif " + elseifBlock.Condition.String() + " then " + elseifBlock.Block.String()
	}
	if is.ElseBlock != nil {
		out += " else " + is.ElseBlock.String()
	}
	out += " endif"
	return out
}

// ElseIfBlock represents an 'elseif' block within an IfStatement.
type ElseIfBlock struct {
	Condition Expression
	Block     *BlockStatement
}

// BlockStatement represents a block of statements (inside if, while, function, etc.).
type BlockStatement struct {
	Token      token.Token // The '{' token (though we don't use braces in WordLang, we can use the first token of the block)
	Statements []Statement
}

func (bs *BlockStatement) statementNode()     {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out string
	out += "{\n" // For visual representation of blocks
	for _, s := range bs.Statements {
		out += "  " + s.String() + "\n" // Indent for block content
	}
	out += "}\n" // End of block
	return out
}


// WhileStatement represents a 'while' loop.
type WhileStatement struct {
	Token     token.Token // The 'while' token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()     {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	return "while " + ws.Condition.String() + " do " + ws.Body.String() + " endwhile"
}

// ForEachStatement represents a 'for each' loop.
type ForEachStatement struct {
	Token    token.Token // The 'foreach' token
	Variable *Identifier
	Iterable Expression // Expression that should evaluate to a list
	Body     *BlockStatement
}

func (fes *ForEachStatement) statementNode()     {}
func (fes *ForEachStatement) TokenLiteral() string { return fes.Token.Literal }
func (fes *ForEachStatement) String() string {
	return "foreach " + fes.Variable.String() + " in " + fes.Iterable.String() + " do " + fes.Body.String() + " endforeach"
}

// FunctionLiteral represents a function definition.
type FunctionLiteral struct {
	Token      token.Token // The 'function' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()    {} // Functions are expressions in some contexts (e.g., function literals)
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	return "function(" + strings.Join(params, ", ") + ") " + fl.Body.String() + " end function"
}


// CallExpression represents a function call.
type CallExpression struct {
	Token     token.Token // The 'call' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()    {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	return "call " + ce.Function.String() + "(" + strings.Join(args, ", ") + ")" // Parentheses for arguments for now, might reconsider
}

// PrintStatement represents a 'print' statement.
type PrintStatement struct {
	Token token.Token // The 'print' token
	Value Expression
}

func (ps *PrintStatement) statementNode()     {}
func (ps *PrintStatement) TokenLiteral() string { return ps.Token.Literal }
func (ps *PrintStatement) String() string {
	return "print " + ps.Value.String()
}

// InputStatement represents an 'input' statement.
type InputStatement struct {
	Token token.Token // The 'input' token
	Prompt *StringLiteral // Optional prompt string
}

func (is *InputStatement) statementNode() {}
func (is *InputStatement) TokenLiteral() string { return is.Token.Literal }
func (is *InputStatement) String() string {
	if is.Prompt != nil {
		return "input " + is.Prompt.String()
	}
	return "input"
}

// ListLiteral represents a list literal.
type ListLiteral struct {
	Token    token.Token // The 'list' token
	Elements []Expression
}

func (ll *ListLiteral) expressionNode()    {}
func (ll *ListLiteral) TokenLiteral() string { return ll.Token.Literal }
func (ll *ListLiteral) String() string {
	elems := []string{}
	for _, el := range ll.Elements {
		elems = append(elems, el.String())
	}
	return "list(" + strings.Join(elems, ", ") + ")" // Parentheses for list elements for now, reconsider
}

// GetItemAtIndexExpression represents getting an item from a list at a specific index.
type GetItemAtIndexExpression struct {
	Token token.Token // The 'get item at index' token
	List Expression
	Index Expression
}

func (giae *GetItemAtIndexExpression) expressionNode()    {}
func (giae *GetItemAtIndexExpression) TokenLiteral() string { return giae.Token.Literal }
func (giae *GetItemAtIndexExpression) String() string {
	return "get item at index " + giae.Index.String() + " from " + giae.List.String()
}

// IsDefinedExpression checks if a variable is defined.
type IsDefinedExpression struct {
	Token token.Token // The 'is defined' token
	Identifier *Identifier
}

func (ide *IsDefinedExpression) expressionNode() {}
func (ide *IsDefinedExpression) TokenLiteral() string { return ide.Token.Literal }
func (ide *IsDefinedExpression) String() string {
	return "is defined " + ide.Identifier.String()
}

// ExitStatement represents the 'exit' statement.
type ExitStatement struct {
	Token token.Token // The 'exit' token
	Code Expression // Optional exit code
}

func (es *ExitStatement) statementNode() {}
func (es *ExitStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExitStatement) String() string {
	if es.Code != nil {
		return "exit " + es.Code.String()
	}
	return "exit"
}

// ConvertToNumberExpression represents converting an expression to a number.
type ConvertToNumberExpression struct {
	Token token.Token // The 'convert to number' token
	Expression Expression
}

func (ctne *ConvertToNumberExpression) expressionNode() {}
func (ctne *ConvertToNumberExpression) TokenLiteral() string { return ctne.Token.Literal }
func (ctne *ConvertToNumberExpression) String() string {
	return "convert to number " + ctne.Expression.String()
}

// ConvertToStringExpression represents converting an expression to a string.
type ConvertToStringExpression struct {
	Token token.Token // The 'convert to string' token
	Expression Expression
}

func (ctse *ConvertToStringExpression) expressionNode() {}
func (ctse *ConvertToStringExpression) TokenLiteral() string { return ctse.Token.Literal }
func (ctse *ConvertToStringExpression) String() string {
	return "convert to string " + ctse.Expression.String()
}
