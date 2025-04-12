package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type visitorExpression struct {
	Err   error
	Stack []ast.Node
	nodes []ast.Node
}

func (v *visitorExpression) Visit(current ast.Node) ast.Visitor {
	if current != nil {
		// Process after entered
		v.nodes = append(v.nodes, current)
		return v
	}
	if err := appendNode(v); err != nil {
		// Process before exited
		v.Err = err
		return nil
	}
	return nil
}

func appendNode(visitor *visitorExpression) error {
	if len(visitor.nodes) < 1 {
		return nil
	}
	var current ast.Node
	current, visitor.nodes = visitor.nodes[len(visitor.nodes)-1], visitor.nodes[:len(visitor.nodes)-1]
	// fmt.Printf("onExit: %+v\n", current)

	switch n := current.(type) {
	case *ast.BinaryExpr:
		// fmt.Printf("Bin : %v %s %v\n", n.X, n.Op, n.Y)
		if n.Op != token.LAND && n.Op != token.LOR {
			return fmt.Errorf("Unsupported op: %s", n.Op)
		}
		visitor.Stack = append(visitor.Stack, n)
		return nil
	case *ast.BasicLit:
		// fmt.Printf("Lit : %s %s\n", n.Kind, n.Value)
		visitor.Stack = append(visitor.Stack, n)
		return nil
	case *ast.CallExpr:
		// fmt.Printf("Call: %s %v\n", n.Fun, n.Args)
		visitor.Stack = append(visitor.Stack, n)
		return nil
	case *ast.Ident:
		if n.Name == "true" || n.Name == "false" {
			visitor.Stack = append(visitor.Stack, n)
		}
		return nil
	}
	return nil
}

// ParseExpression は式を解析し、逆ポーランド記法っぽいやつを返す
// この逆ポーランド記法は、通常のポーランド記法に関数が含まれているものである
func parseExpression(input string) (*[]ast.Node, error) {
	// Generate AST
	root, err := parser.ParseExpr(input)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	visitor := visitorExpression{
		Err: nil,
	}
	// Traverse to generate Reverse Polish Notation
	ast.Walk(&visitor, root)
	if visitor.Err != nil {
		return nil, visitor.Err
	}
	return &visitor.Stack, nil
}
