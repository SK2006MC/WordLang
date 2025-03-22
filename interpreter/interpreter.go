package interpreter

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"wordlang/ast"
	"wordlang/object"
)

// Environment holds variable bindings.
type Environment struct {
	store map[string]object.Object
	outer *Environment // For scopes (not implemented yet in this basic version)
}

// NewEnvironment creates a new environment.
func NewEnvironment() *Environment {
	s := make(map[string]object.Object)
	return &Environment{store: s, outer: nil}
}

// Get retrieves a variable from the environment.
func (e *Environment) Get(name string) (object.Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil { // Scope lookup (not fully implemented yet)
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set sets a variable in the environment.
func (e *Environment) Set(name string, val object.Object) object.Object {
	e.store[name] = val
	return val
}


// Eval evaluates an AST node.
func Eval(node ast.Node, env *Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node)
	case *ast.FloatLiteral:
		return evalFloatLiteral(node)
	case *ast.StringLiteral:
		return evalStringLiteral(node)
	case *ast.BooleanLiteral:
		return evalBooleanLiteral(node)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfStatement:
		return evalIfStatement(node, env)
	case *ast.WhileStatement:
		return evalWhileStatement(node, env)
	case *ast.ForEachStatement:
		return evalForEachStatement(node, env)
	case *ast.LetStatement:
		return evalLetStatement(node, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env) // Placeholder, needs actual return value handling
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.PrintStatement:
		return evalPrintStatement(node, env)
	case *ast.InputStatement:
		return evalInputStatement(node, env)
	case *ast.ListLiteral:
		return evalListLiteral(node, env)
	case *ast.GetItemAtIndexExpression:
		return evalGetItemAtIndexExpression(node, env)
	case *ast.IsDefinedExpression:
		return evalIsDefinedExpression(node, env)
	case *ast.ExitStatement:
		return evalExitStatement(node, env)
	case *ast.ConvertToNumberExpression:
		return evalConvertToNumberExpression(node, env)
	case *ast.ConvertToStringExpression:
		return evalConvertToStringExpression(node, env)
	default:
		return object.NewError("Eval: Node type not handled: %T", node)
	}
}

func evalProgram(program *ast.Program, env *Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		if returnValue, ok := result.(*object.ReturnValue); ok { // Basic return statement handling
			return returnValue.Value
		}

		if errObj, ok := result.(*object.Error); ok {
			return errObj // Propagate errors
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalIntegerLiteral(il *ast.IntegerLiteral) object.Object {
	return &object.Integer{Value: il.Value}
}

func evalFloatLiteral(fl *ast.FloatLiteral) object.Object {
	return &object.Float{Value: fl.Value}
}

func evalStringLiteral(sl *ast.StringLiteral) object.Object {
	return &object.String{Value: sl.Value}
}

func evalBooleanLiteral(bl *ast.BooleanLiteral) object.Object {
	return nativeBoolToBooleanObject(bl.Value)
}

func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return object.TRUE
	}
	return object.FALSE
}

func evalPrefixExpression(pe *ast.PrefixExpression, env *Environment) object.Object {
	right := Eval(pe.Right, env)
	if isError(right) {
		return right
	}

	switch pe.Operator {
	case "not":
		return evalNotOperatorExpression(right)
	default:
		return object.NewError("Eval: Unknown prefix operator: %s%s", pe.Operator, right.Type())
	}
}

func evalNotOperatorExpression(right object.Object) object.Object {
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NULL:
		return object.TRUE // 'not null' is true
	default:
		return object.FALSE // 'not <anything else>' is false (for simplicity in this example)
	}
}

func evalInfixExpression(ie *ast.InfixExpression, env *Environment) object.Object {
	left := Eval(ie.Left, env)
	if isError(left) {
		return left
	}

	right := Eval(ie.Right, env)
	if isError(right) {
		return right
	}

	switch ie.Operator {
	case "add":
		return evalAddInfixExpression(ie.Operator, left, right)
	case "subtract":
		return evalSubtractInfixExpression(ie.Operator, left, right)
	case "multiply":
		return evalMultiplyInfixExpression(ie.Operator, left, right)
	case "divide":
		return evalDivideInfixExpression(ie.Operator, left, right)
	case "equals":
		return evalEqualsInfixExpression(ie.Operator, left, right)
	case "notequals":
		return evalNotEqualsInfixExpression(ie.Operator, left, right)
	case "greater":
		return evalGreaterThanInfixExpression(ie.Operator, left, right)
	case "less":
		return evalLessThanInfixExpression(ie.Operator, left, right)
	case "greater or equal":
		return evalGreaterOrEqualInfixExpression(ie.Operator, left, right)
	case "less or equal":
		return evalLessOrEqualInfixExpression(ie.Operator, left, right)
	case "and":
		return evalAndInfixExpression(ie.Operator, left, right)
	case "or":
		return evalOrInfixExpression(ie.Operator, left, right)
	default:
		return object.NewError("Eval: Unknown infix operator: %s %s %s", left.Type(), ie.Operator, right.Type())
	}
}

func evalAddInfixExpression(operator string, left, right object.Object) object.Object {
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Integer).Value
		rightVal := right.(*object.Integer).Value
		return &object.Integer{Value: leftVal + rightVal}
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := right.(*object.Float).Value
		return &object.Float{Value: leftVal + rightVal}
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := float64(right.(*object.Integer).Value)
		return &object.Float{Value: leftVal + rightVal}
	}
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := float64(left.(*object.Integer).Value)
		rightVal := right.(*object.Float).Value
		return &object.Float{Value: leftVal + rightVal}
	}
	if left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ {
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return &object.String{Value: leftVal + rightVal} // String concatenation
	}
	return object.NewError("Eval: Type mismatch for '%s' operator: %s %s %s", operator, left.Type(), operator, right.Type())
}

func evalSubtractInfixExpression(operator string, left, right object.Object) object.Object {
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Integer).Value
		rightVal := right.(*object.Integer).Value
		return &object.Integer{Value: leftVal - rightVal}
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := right.(*object.Float).Value
		return &object.Float{Value: leftVal - rightVal}
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := float64(right.(*object.Integer).Value)
		return &object.Float{Value: leftVal - rightVal}
	}
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := float64(left.(*object.Integer).Value)
		rightVal := right.(*object.Float).Value
		return &object.Float{Value: leftVal - rightVal}
	}
	return object.NewError("Eval: Type mismatch for '%s' operator: %s %s %s", operator, left.Type(), operator, right.Type())
}

func evalMultiplyInfixExpression(operator string, left, right object.Object) object.Object {
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Integer).Value
		rightVal := right.(*object.Integer).Value
		return &object.Integer{Value: leftVal * rightVal}
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := right.(*object.Float).Value
		return &object.Float{Value: leftVal * rightVal}
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := float64(right.(*object.Integer).Value)
		return &object.Float{Value: leftVal * rightVal}
	}
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := float64(left.(*object.Integer).Value)
		rightVal := right.(*object.Float).Value
		return &object.Float{Value: leftVal * rightVal}
	}
	return object.NewError("Eval: Type mismatch for '%s' operator: %s %s %s", operator, left.Type(), operator, right.Type())
}

func evalDivideInfixExpression(operator string, left, right object.Object) object.Object {
	if right.(*object.Integer).Value == 0 && right.Type() == object.INTEGER_OBJ || right.(*object.Float).Value == 0 && right.Type() == object.FLOAT_OBJ{
		return object.NewError("Eval: Division by zero error")
	}
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Integer).Value
		rightVal := right.(*object.Integer).Value
		return &object.Integer{Value: leftVal / rightVal}
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := right.(*object.Float).Value
		return &object.Float{Value: leftVal / rightVal}
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := float64(right.(*object.Integer).Value)
		return &object.Float{Value: leftVal / rightVal}
	}
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := float64(left.(*object.Integer).Value)
		rightVal := right.(*object.Float).Value
		return &object.Float{Value: leftVal / rightVal}
	}
	return object.NewError("Eval: Type mismatch for '%s' operator: %s %s %s", operator, left.Type(), operator, right.Type())
}


func evalEqualsInfixExpression(operator string, left, right object.Object) object.Object {
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Integer).Value
		rightVal := right.(*object.Integer).Value
		return nativeBoolToBooleanObject(leftVal == rightVal)
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := right.(*object.Float).Value
		return nativeBoolToBooleanObject(leftVal == rightVal)
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := float64(right.(*object.Integer).Value)
		return nativeBoolToBooleanObject(leftVal == rightVal)
	}
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := float64(left.(*object.Integer).Value)
		rightVal := right.(*object.Float).Value
		return nativeBoolToBooleanObject(leftVal == rightVal)
	}
	if left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ {
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return nativeBoolToBooleanObject(leftVal == rightVal)
	}
	if left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ {
		leftVal := left.(*object.Boolean).Value
		rightVal := right.(*object.Boolean).Value
		return nativeBoolToBooleanObject(leftVal == rightVal)
	}
	return nativeBoolToBooleanObject(left == right) // Object reference equality if types differ
}

func evalNotEqualsInfixExpression(operator string, left, right object.Object) object.Object {
	equalsResult := evalEqualsInfixExpression("equals", left, right) // Reuse equals logic
	if isError(equalsResult) {
		return equalsResult
	}
	return evalNotOperatorExpression(equalsResult) // Invert the result of equals
}

func evalGreaterThanInfixExpression(operator string, left, right object.Object) object.Object {
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Integer).Value
		rightVal := right.(*object.Integer).Value
		return nativeBoolToBooleanObject(leftVal > rightVal)
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := right.(*object.Float).Value
		return nativeBoolToBooleanObject(leftVal > rightVal)
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := float64(right.(*object.Integer).Value)
		return nativeBoolToBooleanObject(leftVal > rightVal)
	}
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := float64(left.(*object.Integer).Value)
		rightVal := right.(*object.Float).Value
		return nativeBoolToBooleanObject(leftVal > rightVal)
	}
	return object.NewError("Eval: Type mismatch for '%s' operator: %s %s %s", operator, left.Type(), operator, right.Type())
}

func evalLessThanInfixExpression(operator string, left, right object.Object) object.Object {
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Integer).Value
		rightVal := right.(*object.Integer).Value
		return nativeBoolToBooleanObject(leftVal < rightVal)
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := right.(*object.Float).Value
		return nativeBoolToBooleanObject(leftVal < rightVal)
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := float64(right.(*object.Integer).Value)
		return nativeBoolToBooleanObject(leftVal < rightVal)
	}
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := float64(left.(*object.Integer).Value)
		rightVal := right.(*object.Float).Value
		return nativeBoolToBooleanObject(leftVal < rightVal)
	}
	return object.NewError("Eval: Type mismatch for '%s' operator: %s %s %s", operator, left.Type(), operator, right.Type())
}

func evalGreaterOrEqualInfixExpression(operator string, left, right object.Object) object.Object {
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Integer).Value
		rightVal := right.(*object.Integer).Value
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := right.(*object.Float).Value
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := float64(right.(*object.Integer).Value)
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	}
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := float64(left.(*object.Integer).Value)
		rightVal := right.(*object.Float).Value
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	}
	return object.NewError("Eval: Type mismatch for '%s' operator: %s %s %s", operator, left.Type(), operator, right.Type())
}

func evalLessOrEqualInfixExpression(operator string, left, right object.Object) object.Object {
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Integer).Value
		rightVal := right.(*object.Integer).Value
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := right.(*object.Float).Value
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ {
		leftVal := left.(*object.Float).Value
		rightVal := float64(right.(*object.Integer).Value)
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	}
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ {
		leftVal := float64(left.(*object.Integer).Value)
		rightVal := right.(*object.Float).Value
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	}
	return object.NewError("Eval: Type mismatch for '%s' operator: %s %s %s", operator, left.Type(), operator, right.Type())
}

func evalAndInfixExpression(operator string, left, right object.Object) object.Object {
	leftBool := isTruthy(left)
	rightBool := isTruthy(right)
	return nativeBoolToBooleanObject(leftBool && rightBool)
}

func evalOrInfixExpression(operator string, left, right object.Object) object.Object {
	leftBool := isTruthy(left)
	rightBool := isTruthy(right)
	return nativeBoolToBooleanObject(leftBool || rightBool)
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case object.NULL:
		return false
	case object.TRUE:
		return true
	case object.FALSE:
		return false
	default:
		return true // Everything else is considered truthy for simplicity
	}
}

func evalIfStatement(is *ast.IfStatement, env *Environment) object.Object {
	condition := Eval(is.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(is.ThenBlock, env)
	} else {
		for _, elseifBlock := range is.ElseIfBlocks {
			elseifCondition := Eval(elseifBlock.Condition, env)
			if isError(elseifCondition) {
				return elseifCondition
			}
			if isTruthy(elseifCondition) {
				return Eval(elseifBlock.Block, env)
			}
		}
		if is.ElseBlock != nil {
			return Eval(is.ElseBlock, env)
		}
	}

	return object.NULL // No 'else' or condition not met, returns null
}

func evalWhileStatement(ws *ast.WhileStatement, env *Environment) object.Object {
	var result object.Object = object.NULL // Default return value

	for {
		condition := Eval(ws.Condition, env)
		if isError(condition) {
			return condition
		}
		if !isTruthy(condition) {
			break // Exit loop if condition is false
		}

		blockResult := Eval(ws.Body, env) // Execute loop body
		if blockResult != nil && blockResult.Type() == object.RETURN_VALUE_OBJ {
			return blockResult // Handle return statements inside loops
		}
		if isError(blockResult) {
			return blockResult // Propagate errors
		}
		result = blockResult // Keep track of last evaluated value in the block (though might not be needed for 'while')
	}

	return result
}

func evalForEachStatement(fes *ast.ForEachStatement, env *Environment) object.Object {
	iterable := Eval(fes.Iterable, env)
	if isError(iterable) {
		return iterable
	}

	listObj, ok := iterable.(*object.List)
	if !ok {
		return object.NewError("Eval: 'for each' loop requires a list as iterable, got %s", iterable.Type())
	}

	var result object.Object = object.NULL // Default return value

	for _, element := range listObj.Elements {
		currentEnv := NewEnclosedEnvironment(env) // Create new scope for each iteration
		currentEnv.Set(fes.Variable.Value, element)    // Bind loop variable
		blockResult := Eval(fes.Body, currentEnv)        // Execute loop body in new scope

		if blockResult != nil && blockResult.Type() == object.RETURN_VALUE_OBJ {
			return blockResult // Handle return statements inside loops
		}
		if isError(blockResult) {
			return blockResult // Propagate errors
		}
		result = blockResult // Keep track of last evaluated value in the block
	}

	return result
}

// NewEnclosedEnvironment creates a new environment enclosed by outer environment.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}


func evalLetStatement(ls *ast.LetStatement, env *Environment) object.Object {
	val := Eval(ls.Value, env)
	if isError(val) {
		return val
	}
	env.Set(ls.Name.Value, val) // Store in the environment
	return val
}

func evalReturnStatement(rs *ast.ReturnStatement, env *Environment) object.Object {
	val := Eval(rs.ReturnValue, env)
	if isError(val) {
		return val
	}
	return &object.ReturnValue{Value: val} // Wrap in ReturnValue object
}

func evalIdentifier(node *ast.Identifier, env *Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return object.NewError("Eval: Identifier not found: %s", node.Value)
	}
	return val
}

func evalPrintStatement(ps *ast.PrintStatement, env *Environment) object.Object {
	value := Eval(ps.Value, env)
	if isError(value) {
		return value
	}
	fmt.Println(value.Inspect()) // Use Inspect for string representation
	return object.NULL
}

func evalInputStatement(is *ast.InputStatement, env *Environment) object.Object {
	var prompt string
	if is.Prompt != nil {
		prompt = is.Prompt.Value
		fmt.Print(prompt)
	}
	var input string
	fmt.Scanln(&input) // Read a line of input
	return &object.String{Value: input}
}

func evalListLiteral(ll *ast.ListLiteral, env *Environment) object.Object {
	elements := evalExpressions(ll.Elements, env)
	if len(elements) > 0 && isError(elements[0]) { // Check for error in first element eval
		return elements[0]
	}
	return &object.List{Elements: elements}
}

func evalExpressions(exps []ast.Expression, env *Environment) []object.Object {
	var results []object.Object
	for _, exp := range exps {
		evaluated := Eval(exp, env)
		if isError(evaluated) {
			return []object.Object{evaluated} // Return error immediately
		}
		results = append(results, evaluated)
	}
	return results
}

func evalGetItemAtIndexExpression(giae *ast.GetItemAtIndexExpression, env *Environment) object.Object {
	listObj := Eval(giae.List, env)
	if isError(listObj) {
		return listObj
	}
	list, ok := listObj.(*object.List)
	if !ok {
		return object.NewError("Eval: 'get item at index' expected a list, got %s", listObj.Type())
	}

	indexObj := Eval(giae.Index, env)
	if isError(indexObj) {
		return indexObj
	}
	index, ok := indexObj.(*object.Integer)
	if !ok {
		return object.NewError("Eval: 'get item at index' index must be a number, got %s", indexObj.Type())
	}

	if index.Value < 0 || index.Value >= int64(len(list.Elements)) {
		return object.NewError("Eval: Index out of bounds: %d, list length: %d", index.Value, len(list.Elements))
	}

	return list.Elements[index.Value]
}

func evalIsDefinedExpression(ide *ast.IsDefinedExpression, env *Environment) object.Object {
	_, ok := env.Get(ide.Identifier.Value)
	return nativeBoolToBooleanObject(ok) // Returns true if defined, false otherwise
}

func evalExitStatement(es *ast.ExitStatement, env *Environment) object.Object {
	code := 0 // Default exit code
	if es.Code != nil {
		codeObj := Eval(es.Code, env)
		if isError(codeObj) {
			fmt.Println(codeObj.Inspect()) // Print error before exiting
			code = 1 // Error exit code in case of evaluation error
		} else if intCode, ok := codeObj.(*object.Integer); ok {
			code = int(intCode.Value)
		} else {
			fmt.Println(object.NewError("Eval: Exit code must be an integer, got %s", codeObj.Type()).Inspect())
			code = 1
		}
	}
	os.Exit(code)
	return object.NULL // Should not reach here
}

func evalConvertToNumberExpression(ctne *ast.ConvertToNumberExpression, env *Environment) object.Object {
	expValue := Eval(ctne.Expression, env)
	if isError(expValue) {
		return expValue
	}

	switch value := expValue.(type) {
	case *object.String:
		floatVal, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return object.NewError("Eval: Cannot convert string '%s' to number: %s", value.Value, err.Error())
		}
		if strings.Contains(value.Value, ".") {
			return &object.Float{Value: floatVal}
		}
		intVal := int64(floatVal) // Truncate to integer if no decimal point in original string
		return &object.Integer{Value: intVal}
	case *object.Integer:
		return value // Already a number
	case *object.Float:
		return value // Already a number
	default:
		return object.NewError("Eval: Cannot convert type %s to number", expValue.Type())
	}
}

func evalConvertToStringExpression(ctse *ast.ConvertToStringExpression, env *Environment) object.Object {
	expValue := Eval(ctse.Expression, env)
	if isError(expValue) {
		return expValue
	}
	return &object.String{Value: expValue.Inspect()} // Use Inspect() to get string representation
}


func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
