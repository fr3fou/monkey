package evaluator

import (
	"github.com/fr3fou/monkey/ast"
	"github.com/fr3fou/monkey/object"
)

var (
	// NULL represents lack of value
	NULL = &object.Null{}
	// TRUE represents a true value
	TRUE = &object.Boolean{Value: true}
	// FALSE represents a false value
	FALSE = &object.Boolean{Value: false}
)

// Eval takes in an ast.Node and evaluates it
func Eval(n ast.Node) object.Object {
	switch node := n.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	}

	return nil
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixExpression(right)
	default:
		return NULL
	}
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return nil
	}

	val := right.(*object.Integer).Value
	return &object.Integer{
		Value: -val,
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}
