package parser

import (
	"fmt"
	"strconv"
	"wordlang/ast"
	"wordlang/lexer"
	"wordlang/token"
)

// Parser holds the state for parsing.
type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns   map[token.TokenType]prefixParseFn
	infixParseFns    map[token.TokenType]infixParseFn
	statementParseFns map[token.TokenType]statementParseFn // Add statementParseFns
}

type (
	//prefixParseFn   func() ast.Expression
	//infixParseFn    func(ast.Expression) ast.Expression
	statementParseFn func() ast.Statement // <--- New type for statement parsing
)

func (p *Parser) registerStatement(tokenType token.TokenType, fn statementParseFn) {
	p.statementParseFns[tokenType] = fn // <--- New function to register statement parsers
}

// New creates a new Parser.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:                 l,
		errors:            []string{},
		prefixParseFns:    make(map[token.TokenType]prefixParseFn),
		infixParseFns:     make(map[token.TokenType]infixParseFn),
		statementParseFns: make(map[token.TokenType]statementParseFn), // Initialize statementParseFns
	}

	p.registerParseFunctions() // Register parse functions

	// Read two tokens to set curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead at line %d, column %d",
		t, p.peekToken.Type, p.peekToken.Line, p.peekToken.Column)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// ParseProgram parses the entire program.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}
func (p *Parser) parseStatement() ast.Statement {
	if parseStatementFn, ok := p.statementParseFns[p.curToken.Type]; ok {
		return parseStatementFn() // Call the registered statement parser
	}

	// If no statement parser is found, default to an expression statement
	return p.parseExpressionStatement()
}

/*func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.PRINT:
		return p.parsePrintStatement()
	case token.INPUT:
		return p.parseInputStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOREACH:
		return p.parseForEachStatement()
	case token.FUNCTION: // Function definition as expression statement
		return p.parseExpressionStatement() // Parse function as expression statement
	case token.CALL: // Function call as expression statement
		return p.parseExpressionStatement() // Parse call as expression statement
	case token.EXIT:
		return p.parseExitStatement()
	default:
		return p.parseExpressionStatement() // Default to expression statement
	}
}*/

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	fmt.Println("parseLetStatement: curToken=", p.curToken, ", peekToken=", p.peekToken) // Debug print

	if !p.expectPeek(token.IDENT) {
		fmt.Println("parseLetStatement: expectPeek(IDENT) failed, peekToken=", p.peekToken) // Debug print
		return nil // Error already added by expectPeek
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	fmt.Println("parseLetStatement: after IDENT, curToken=", p.curToken, ", peekToken=", p.peekToken) // Debug print

	if !p.expectPeek(token.BE) { // Expect 'be' after variable name
		fmt.Println("parseLetStatement: expectPeek(BE) failed, peekToken=", p.peekToken) // Debug print
		return nil
	}

	p.nextToken() // Consume 'be', move to the expression
	stmt.Value = p.parseExpression(LOWEST) // Parse the value expression

	// Semicolon handling might be different in WordLang, we'll assume newline or 'end' for statement termination for now.

	return stmt
}


func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken() // Move past 'return'

	stmt.ReturnValue = p.parseExpression(LOWEST) // Parse the return value expression

	// Semicolon/newline handling similar to let statement

	return stmt
}


func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	// Semicolon/newline handling
	return stmt
}

// Operator Precedence (Word-based operators, needs careful thought for WordLang)
const (
	_ int = iota
	LOWEST
	EQUALS_PREC      // equals, not equals
	LESSGREATER_PREC // greater than, less than, etc.
	SUM_PREC         // add, subtract
	PRODUCT_PREC     // multiply, divide
	PREFIX_PREC      // not
	CALL_PREC
	INDEX_PREC
)

var precedence = map[token.TokenType]int{
	token.EQUALS:      EQUALS_PREC,
	token.NOTEQUALS:   EQUALS_PREC,
	token.GREATERTHAN: LESSGREATER_PREC,
	token.LESSTHAN:    LESSGREATER_PREC,
	token.GREATEREQUAL:LESSGREATER_PREC,
	token.LESSEQUAL:   LESSGREATER_PREC,
	token.ADD:         SUM_PREC,
	token.SUBTRACT:    SUM_PREC,
	token.MULTIPLY:    PRODUCT_PREC,
	token.DIVIDE:      PRODUCT_PREC,
	token.OR:          EQUALS_PREC, // Example precedence - adjust as needed
	token.AND:         EQUALS_PREC, // Example precedence - adjust as needed
	token.CALL:        CALL_PREC,
	token.GETITEMATINDEX: INDEX_PREC, // Example precedence
}


func (p *Parser) peekPrecedence() int {
	if p, ok := precedence[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedence[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}


func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFn := p.prefixParseFns[p.curToken.Type]
	if prefixFn == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefixFn() // Just call the prefix function and return

	// --- REMOVE ALL INFIX PARSING LOGIC ---
	// No more infix loop or infix function calls here for now.

	return leftExp
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found at line %d, column %d", t, p.curToken.Line, p.curToken.Column)
	p.errors = append(p.errors, msg)
}


func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer at line %d, column %d", p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float at line %d, column %d", p.curToken.Literal, p.curToken.Line, p.curToken.Column)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}


func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal, // Operator will be the keyword like "not"
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX_PREC) // Parse the right-hand side expression

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal, // Operator will be the keyword like "add", "equals", etc.
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence) // Parse the right-hand side expression

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	// In WordLang, we might not have parentheses for grouping in the traditional sense.
	// We could use keywords for explicit grouping if needed, but for now, we'll skip explicit grouping for this basic example.
	// If we decide to add grouping later (e.g., with "group ... end group" keywords), this is where we'd handle it.

	// For now, just return the inner expression as if there were no grouping.
	p.nextToken() // Consume the opening group keyword (if we had one)
	exp := p.parseExpression(LOWEST)
	// Expect closing group keyword (if we had one)
	return exp
}
// Change return type to ast.Statement
func (p *Parser) parseIfStatement() ast.Statement { 
	stmt := &ast.IfStatement{Token: p.curToken, ElseIfBlocks: []*ast.ElseIfBlock{}}

	p.nextToken() // Consume 'if'
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.THEN) { // Expect 'then' after condition
		return nil
	}

	stmt.ThenBlock = p.parseBlockStatement() // Parse the 'then' block

	for p.peekTokenIs(token.ELSEIF) { // Handle multiple 'elseif' blocks
		p.nextToken() // Consume 'elseif'
		elseifBlock := &ast.ElseIfBlock{}
		elseifBlock.Condition = p.parseExpression(LOWEST)
		if !p.expectPeek(token.THEN) {
			return nil
		}
		elseifBlock.Block = p.parseBlockStatement()
		stmt.ElseIfBlocks = append(stmt.ElseIfBlocks, elseifBlock)
	}


	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // Consume 'else'
		stmt.ElseBlock = p.parseBlockStatement() // Parse the 'else' block
	}

	if !p.expectPeek(token.ENDIF) { // Expect 'endif' to close the if statement
		return nil
	}

	return stmt // Still return the *ast.IfStatement, which now satisfies ast.Statement
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken, Statements: []ast.Statement{}}

	p.nextToken() // Consume '{' (though we don't have explicit braces in WordLang, this conceptually starts the block)

	for !p.curTokenIs(token.ENDIF) && !p.curTokenIs(token.ELSE) && !p.curTokenIs(token.ELSEIF) && !p.curTokenIs(token.ENDWHILE) && !p.curTokenIs(token.ENDFOREACH) && !p.curTokenIs(token.END) && !p.curTokenIs(token.EOF) { // Stop at block terminators or EOF
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	p.nextToken() // Consume 'while'
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.DO) { // Expect 'do' after condition
		return nil
	}

	stmt.Body = p.parseBlockStatement() // Parse the loop body

	if !p.expectPeek(token.ENDWHILE) { // Expect 'endwhile' to close the while loop
		return nil
	}

	return stmt
}

func (p *Parser) parseForEachStatement() *ast.ForEachStatement {
	stmt := &ast.ForEachStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) { // Expect identifier for variable name
		return nil
	}
	stmt.Variable = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.IN) { // Expect 'in' keyword
		return nil
	}

	p.nextToken() // Consume 'in'
	stmt.Iterable = p.parseExpression(LOWEST) // Parse the iterable expression (should be a list)

	if !p.expectPeek(token.DO) { // Expect 'do' before loop body
		return nil
	}

	stmt.Body = p.parseBlockStatement() // Parse the loop body

	if !p.expectPeek(token.ENDFOREACH) { // Expect 'endforeach' to close the loop
		return nil
	}

	return stmt
}

// Change return type to ast.Expression
func (p *Parser) parseFunctionStatement() ast.Expression { 
    lit := &ast.FunctionLiteral{Token: p.curToken}

    if p.peekTokenIs(token.IDENT) { // Parameters are optional for now, but if present, expect IDENTs
        p.nextToken()
        lit.Parameters = p.parseFunctionParameters()
    } else {
        lit.Parameters = []*ast.Identifier{} // No parameters
    }

    lit.Body = p.parseBlockStatement() // Parse function body

    if !p.expectPeek(token.ENDFUNCTION) && !p.expectPeek(token.END) { // Expect 'end function' or 'end' to close function definition
        return nil // Or handle error appropriately
    }

    return lit // Return the *ast.FunctionLiteral, which now satisfies ast.Expression
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if !p.curTokenIs(token.IDENT) { // No parameters case handled in parseFunctionStatement
		return identifiers
	}

	identifiers = append(identifiers, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})

	for p.peekTokenIs(token.IDENT) { // Expect more identifiers (parameters), separated by spaces (or commas if we decide to add them back minimally)
		p.nextToken()
		identifiers = append(identifiers, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	return identifiers
}


func (p *Parser) parseCallStatement() *ast.CallExpression { // Changed to CallExpression as function calls are expressions (for now)
	callExp := &ast.CallExpression{Token: p.curToken}

	p.nextToken() // Consume 'call'

	callExp.Function = p.parseExpression(CALL_PREC) // Parse function identifier or function literal

	callExp.Arguments = p.parseCallArguments()

	return callExp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.END) { // No arguments
		p.nextToken()
		return args
	}

	p.nextToken() // Move to the first argument
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.IDENT) || p.peekTokenIs(token.NUMBER) || p.peekTokenIs(token.STRING) || p.peekTokenIs(token.TRUE) || p.peekTokenIs(token.FALSE) || p.peekTokenIs(token.LIST) || p.peekTokenIs(token.GETITEMATINDEX) || p.peekTokenIs(token.CONVERTTONUMBER) || p.peekTokenIs(token.CONVERTTOSTRING){ // Check for tokens that can start an expression argument
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	return args
}

func (p *Parser) parsePrintStatement() *ast.PrintStatement {
	stmt := &ast.PrintStatement{Token: p.curToken}

	p.nextToken() // Consume 'print'

	stmt.Value = p.parseExpression(LOWEST) // Parse the expression to print

	return stmt
}

func (p *Parser) parseInputStatement() *ast.InputStatement {
	stmt := &ast.InputStatement{Token: p.curToken}

	if p.peekTokenIs(token.STRING) { // Optional prompt string
		p.nextToken()
		stmt.Prompt = p.parseStringLiteral().(*ast.StringLiteral)
	}

	return stmt
}

func (p *Parser) parseListLiteral() ast.Expression {
	listLit := &ast.ListLiteral{Token: p.curToken, Elements: []ast.Expression{}}

	if p.peekTokenIs(token.END) { // Empty list
		p.nextToken()
		return listLit
	}

	p.nextToken() // Move to the first element
	listLit.Elements = append(listLit.Elements, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.IDENT) || p.peekTokenIs(token.NUMBER) || p.peekTokenIs(token.STRING) || p.peekTokenIs(token.TRUE) || p.peekTokenIs(token.FALSE) || p.peekTokenIs(token.LIST) || p.peekTokenIs(token.GETITEMATINDEX) || p.peekTokenIs(token.CONVERTTONUMBER) || p.peekTokenIs(token.CONVERTTOSTRING){ // Check for tokens that can start an expression list element
		p.nextToken()
		listLit.Elements = append(listLit.Elements, p.parseExpression(LOWEST))
	}

	return listLit
}

func (p *Parser) parseGetItemAtIndexExpression(list ast.Expression) ast.Expression {
	getItemAtIndexExp := &ast.GetItemAtIndexExpression{Token: p.curToken, List: list}

	if !p.expectPeek(token.INDEX) {
		return nil
	}
	p.nextToken() // Consume 'index'
	getItemAtIndexExp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.FROM) {
		return nil
	}
	p.nextToken() // consume 'from'
	// List is already parsed and passed as 'list' argument to this function

	return getItemAtIndexExp
}

func (p *Parser) parseIsDefinedExpression() ast.Expression {
	isDefinedExp := &ast.IsDefinedExpression{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil // Expect identifier after 'is defined'
	}
	isDefinedExp.Identifier = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return isDefinedExp
}

func (p *Parser) parseExitStatement() *ast.ExitStatement {
	stmt := &ast.ExitStatement{Token: p.curToken}

	if !p.peekTokenIs(token.END) && !p.peekTokenIs(token.EOF) { // Optional exit code
		p.nextToken()
		stmt.Code = p.parseExpression(LOWEST)
	}

	return stmt
}

func (p *Parser) parseConvertToNumberExpression() ast.Expression {
	convExp := &ast.ConvertToNumberExpression{Token: p.curToken}
	p.nextToken() // consume 'convert to number'
	convExp.Expression = p.parseExpression(LOWEST)
	return convExp
}

func (p *Parser) parseConvertToStringExpression() ast.Expression {
	convExp := &ast.ConvertToStringExpression{Token: p.curToken}
	p.nextToken() // consume 'convert to string'
	convExp.Expression = p.parseExpression(LOWEST)
	return convExp
}



// --- Prefix and Infix Function Registration ---
func (p *Parser) registerParseFunctions() {
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	// --- Prefix Parsing Functions (Simplified) ---
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.NUMBER, p.parseIntegerLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	// REMOVE: p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionStatement) // Function literal as expression
	p.registerPrefix(token.LIST, p.parseListLiteral)
	p.registerPrefix(token.GETITEMATINDEX, p.parseGetItemAtIndexPrefix)
	p.registerPrefix(token.ISDEFINED, p.parseIsDefinedExpression)
	p.registerPrefix(token.CONVERTTONUMBER, p.parseConvertToNumberExpression)
	p.registerPrefix(token.CONVERTTOSTRING, p.parseConvertToStringExpression)


	// --- REMOVE ALL INFIX PARSING REGISTRATIONS ---
	// REMOVE: p.registerInfix(token.ADD, p.parseInfixExpression)

	// --- Statement Parsing Registrations (NEW) ---
	p.registerStatement(token.LET, p.parseLetStatement)
	p.registerStatement(token.IF, p.parseIfStatement)
	p.registerStatement(token.PRINT, p.parsePrintStatement)
	p.registerStatement(token.WHILE, p.parseWhileStatement)
	p.registerStatement(token.FOREACH, p.parseForEachStatement)
    p.registerStatement(token.RETURN, p.parseReturnStatement)
	p.registerStatement(token.EXIT, p.parseExitStatement)
	p.registerStatement(token.INPUT, p.parseInputStatement)
	//Add function call statement if applicable:  p.registerStatement(token.CALL, p.parseCallStatement)
}

func (p *Parser) parseGetItemAtIndexPrefix() ast.Expression {
	p.nextToken() // Consume 'get item at index' and move to next token which should be index expression.
	indexExp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.FROM) {
		return nil
	}
	p.nextToken() // Consume 'from'
	listExp := p.parseExpression(LOWEST)

	return &ast.GetItemAtIndexExpression{
		Token: p.curToken, // Token context might need adjustment
		List: listExp,
		Index: indexExp,
	}
}

func (p *Parser) parseGetItemAtIndexInfix(left ast.Expression) ast.Expression {
	getItemAtIndexExp := &ast.GetItemAtIndexExpression{Token: p.curToken, List: left}

	if !p.expectPeek(token.INDEX) {
		return nil
	}
	p.nextToken() // Consume 'index'
	getItemAtIndexExp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.FROM) {
		return nil
	}
	p.nextToken() // consume 'from'
	// List is already parsed as 'left' expression

	return getItemAtIndexExp
}


func (p *Parser) parseCallExpressionInfix(function ast.Expression) ast.Expression {
	callExp := &ast.CallExpression{Token: p.curToken, Function: function}
	callExp.Arguments = p.parseCallArguments()
	return callExp
}


// Errors returns parsing errors.
func (p *Parser) Errors() []string {
	return p.errors
}

