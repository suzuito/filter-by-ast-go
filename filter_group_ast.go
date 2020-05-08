package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"regexp"
	"strings"
)

type filterGroupAst struct {
	expression           string
	normalizedExpression string
	nodes                *[]ast.Node
}

func (f *filterGroupAst) Eval(post *Post) (bool, error) {
	elems, err := evalFunction(post, f.nodes)
	if err != nil {
		return false, err
	}
	return evalElements(elems)
}

func (f *filterGroupAst) String() string {
	return strings.Join(stringsFromNodes(f.nodes), "\n")
}

func stringsFromNodes(nodes *[]ast.Node) []string {
	ret := []string{}
	for _, node := range *nodes {
		switch n := node.(type) {
		case *ast.BinaryExpr:
			ret = append(ret, fmt.Sprintf("BinaryExpr:%s", n.Op.String()))
		case *ast.BasicLit:
			ret = append(ret, fmt.Sprintf("BasicLit:%s %s", n.Kind, n.Value))
		case *ast.CallExpr:
			ret = append(ret, fmt.Sprintf("CallExpr:%s", n.Fun))
		}
	}
	return ret
}

var newlineRegexp = regexp.MustCompile(`\r?\n`)

// Parseは論理式を解析し、FilterGroupを生成する。
func Parse(expression string) (FilterGroup, error) {
	norm := newlineRegexp.ReplaceAllString(expression, "")
	b, err := format.Source([]byte(norm))
	if err != nil {
		return nil, err
	}
	nodes, err := parseExpression(string(b))
	if err != nil {
		return nil, err
	}
	return &filterGroupAst{
		expression:           expression,
		normalizedExpression: norm,
		nodes:                nodes,
	}, nil
}
