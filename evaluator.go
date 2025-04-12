package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
)

type element interface {
	Eval(args ...interface{}) bool
}

type elementOpType int

const (
	elementOpAnd elementOpType = iota
	elementOpOr
)

type elementOp struct {
	t elementOpType
}

func (e *elementOp) Eval(args ...interface{}) bool {
	l := args[0]
	r := args[1]
	if e.t == elementOpAnd {
		return l.(bool) && r.(bool)
	}
	return l.(bool) || r.(bool)
}

func (e *elementOp) String() string {
	if e.t == elementOpAnd {
		return "and"
	}
	return "or"
}

type elementTrue struct {
	name string
}

func (e *elementTrue) Eval(args ...interface{}) bool {
	return true
}

func (e *elementTrue) String() string {
	return fmt.Sprintf("%s: true", e.name)
}

type elementFalse struct {
	name string
}

func (e *elementFalse) Eval(args ...interface{}) bool {
	return false
}

func (e *elementFalse) String() string {
	return fmt.Sprintf("%s: false", e.name)
}

func callFilter(f reflect.Value, in []reflect.Value) (result bool, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			result = false
			err = fmt.Errorf("%s", err1)
		}
	}()
	out := f.Call(in)
	if len(out) != 1 {
		return false, fmt.Errorf("ErrInvalidFilter: Number of return values must be 1")
	}
	if out[0].Kind() != reflect.Bool {
		return false, fmt.Errorf("ErrInvalidFilter: Returned type is not bool")
	}
	ret := out[0].Bool()
	return ret, nil
}

// evalFunctionは逆ポーランド記法に含まれている関数の部分だけを評価する。
// 関数だった部分がブール値に置き換わる。
// 関数の出力はブール値と二値演算（AND, OR）を含む逆ポーランド記法。
func evalFunction(post *Post, nodes *[]ast.Node) (*[]element, error) {
	ret := []element{}
	args := []reflect.Value{reflect.ValueOf(post)}
	for _, node := range *nodes {
		switch n := node.(type) {
		case *ast.CallExpr:
			funcName := fmt.Sprintf("%s", n.Fun)
			// fmt.Printf("Call %s with %+v\n", n.Fun, args)
			filter, exists := Filters[funcName]
			if !exists {
				return nil, fmt.Errorf("Function not found: %s", funcName)
			}
			filterValue := reflect.ValueOf(filter)
			result, err := callFilter(filterValue, args)
			if err != nil {
				return nil, err
			}
			args = []reflect.Value{reflect.ValueOf(post)}
			if result {
				ret = append(ret, &elementTrue{name: funcName})
			} else {
				ret = append(ret, &elementFalse{name: funcName})
			}
		case *ast.BasicLit:
			// args = append(args, n)
			// fmt.Println(n.Kind, n.Value)
			var v reflect.Value
			switch n.Kind {
			case token.STRING:
				value, _ := strconv.Unquote(n.Value)
				v = reflect.ValueOf(value)
			case token.INT:
				i, err := strconv.Atoi(n.Value)
				if err != nil {
					return nil, fmt.Errorf("%w: %+v", err, n)
				}
				v = reflect.ValueOf(i)
			case token.FLOAT:
				i, err := strconv.ParseFloat(n.Value, 64)
				if err != nil {
					return nil, fmt.Errorf("%w: %+v", err, n)
				}
				v = reflect.ValueOf(i)
			default:
				return nil, fmt.Errorf("ErrEvalUnsupportedLiteral: %s %s", n.Kind, n.Value)
			}
			args = append(args, v)
		case *ast.BinaryExpr:
			var t elementOpType
			if n.Op.String() == "&&" {
				t = elementOpAnd
			} else if n.Op.String() == "||" {
				t = elementOpOr
			} else {
				return nil, fmt.Errorf("Unsupported op: %s", n.Op)
			}
			ret = append(
				ret,
				&elementOp{
					t: t,
				},
			)
		case *ast.Ident:
			// args = append(args, n)
			if n.Name != "true" && n.Name != "false" {
				return nil, fmt.Errorf("ErrEvalUnsupportedIdent: %s", n.Name)
			}
			b, err := strconv.ParseBool(n.Name)
			if err != nil {
				return nil, fmt.Errorf("%w: %+v", err, n)
			}
			args = append(args, reflect.ValueOf(b))
		}
	}
	return &ret, nil
}

// ブール値と二値演算子を含む逆ポーランド記法を評価する。
func evalElements(elems *[]element) (bool, error) {
	stack := []bool{}
	for _, elem := range *elems {
		switch e := elem.(type) {
		case *elementOp:
			var poped1 bool
			poped1, stack = stack[len(stack)-1], stack[:len(stack)-1]
			var poped2 bool
			poped2, stack = stack[len(stack)-1], stack[:len(stack)-1]
			var result bool
			if e.t == elementOpAnd {
				result = poped1 && poped2
			}
			if e.t == elementOpOr {
				result = poped1 || poped2
			}
			stack = append(stack, result)
		case *elementTrue:
			stack = append(stack, true)
		case *elementFalse:
			stack = append(stack, false)
		default:
			return false, fmt.Errorf("Invalid element: %+v", e)
		}
	}
	if len(stack) != 1 {
		return false, fmt.Errorf("Invalid computation: %+v", elems)
	}
	return stack[0], nil
}
