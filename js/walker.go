package js

import (
	"fmt"
	"io"

	"github.com/jloup/scraper/js/nodedata"
	"github.com/jloup/scraper/node"
	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/parser"
)

const (
	cMAX_DEPTH = 50
)

func ScrapJS(scrapers []*node.ScraperNode, r io.Reader) ([]map[string]interface{}, error) {
	var err error
	var scraper *node.ScraperNode
	items := make([]map[string]interface{}, 0, 0)

	for _, s := range scrapers {
		s.Init()
	}

	w := jsWalker{nil, func(n nodedata.NodeData, depth int) error {
		for _, scraper = range scrapers {
			//ss := strings.Repeat("  ", depth)
			//fmt.Printf("%s%v (%v)\n", ss, n, depth)
			err = scraper.ProcessNode(&n, &items, depth)
			if err != nil {
				return err
			}
		}

		return nil
	}}

	program, err := parser.ParseFile(nil, "", r, 0)
	if err != nil {
		return nil, err
	}

	for _, declaration := range program.DeclarationList {
		t, ok := declaration.(*ast.FunctionDeclaration)
		if ok {
			w.walkAst(t.Function, 0)
		}
	}

	for _, el := range program.Body {
		w.walkAst(el, 0)
	}

	for _, scraper = range scrapers {
		scraper.End(&items)
	}

	return items, w.err
}

type jsWalker struct {
	err error
	fn  func(nodedata.NodeData, int) error
}

func (j *jsWalker) walkerCallback(n nodedata.NodeData, depth int) error {
	if j.err != nil {
		return j.err
	}

	j.err = j.fn(n, depth)
	return j.err
}

func (j *jsWalker) walkAst(n ast.Node, depth int) int {
	if j.err != nil {
		return depth
	}

	if depth > cMAX_DEPTH {
		j.err = fmt.Errorf("max depth reached")
		return depth
	}

	switch nn := n.(type) {
	case *ast.ArrayLiteral:
		j.walkerCallback(nodedata.NodeData{nodedata.ArrayLiteral, "", ""}, depth)
		for _, element := range nn.Value {
			j.walkAst(element, depth+1)
		}

	case *ast.AssignExpression:
		j.walkerCallback(nodedata.NodeData{nodedata.AssignExpression, "", ""}, depth)
		d := j.walkAst(nn.Left, depth)
		j.walkAst(nn.Right, d+1)

	case *ast.Identifier:
		j.walkerCallback(nodedata.NodeData{nodedata.Identifier, nn.Name, ""}, depth)

	case *ast.ObjectLiteral:
		j.walkerCallback(nodedata.NodeData{nodedata.ObjectLiteral, "", ""}, depth)
		for _, key := range nn.Value {
			j.walkerCallback(nodedata.NodeData{nodedata.Property, key.Key, ""}, depth+1)
			j.walkAst(key.Value, depth+2)
		}

	case *ast.VariableExpression:
		j.walkerCallback(nodedata.NodeData{nodedata.VariableExpression, nn.Name, ""}, depth)
		if nn.Initializer != nil {
			j.walkAst(nn.Initializer, depth+1)
		}

	case *ast.ReturnStatement:
		j.walkerCallback(nodedata.NodeData{nodedata.ReturnStatement, "", ""}, depth)
		j.walkAst(nn.Argument, depth+1)

	case *ast.DotExpression:
		j.walkAst(nn.Left, depth)
		depth = j.walkAst(&nn.Identifier, depth+1)

	case *ast.FunctionLiteral:
		var name string
		if nn.Name != nil {
			name = nn.Name.Name
		}

		j.walkerCallback(nodedata.NodeData{nodedata.FunctionLiteral, name, ""}, depth)

		if nn.ParameterList != nil {
			for _, element := range nn.ParameterList.List {
				j.walkAst(element, depth+1)
			}
		}

		j.walkAst(nn.Body, depth+1)

	case *ast.NullLiteral:
		j.walkerCallback(nodedata.NodeData{nodedata.NullLiteral, "", nn.Literal}, depth)

	case *ast.BooleanLiteral:
		j.walkerCallback(nodedata.NodeData{nodedata.BooleanLiteral, "", nn.Literal}, depth)

	case *ast.NumberLiteral:
		j.walkerCallback(nodedata.NodeData{nodedata.NumberLiteral, "", nn.Literal}, depth)

	case *ast.StringLiteral:
		// remove leading and trailing "
		j.walkerCallback(nodedata.NodeData{nodedata.StringLiteral, "", nn.Literal[1 : len(nn.Literal)-1]}, depth)

	case *ast.ThisExpression:
		j.walkerCallback(nodedata.NodeData{nodedata.ThisExpression, "this", ""}, depth)

	case *ast.RegExpLiteral:
		j.walkerCallback(nodedata.NodeData{nodedata.RegExpLiteral, "", nn.Literal}, depth)

	case *ast.BracketExpression:
		j.walkAst(nn.Left, depth)
		j.walkAst(nn.Member, depth+1)

	case *ast.BinaryExpression:
		j.walkAst(nn.Left, depth)
		j.walkAst(nn.Right, depth)

	case *ast.CallExpression:
		d := j.walkAst(nn.Callee, depth)
		for _, arg := range nn.ArgumentList {
			j.walkAst(arg, d+1)
		}

	case *ast.BlockStatement:
		for _, element := range nn.List {
			j.walkAst(element, depth+1)
		}

	case *ast.SequenceExpression:
		for _, element := range nn.Sequence {
			j.walkAst(element, depth+1)
		}

	case *ast.ExpressionStatement:
		depth = j.walkAst(nn.Expression, depth)

	case *ast.IfStatement:
		j.walkAst(nn.Test, depth)
		j.walkAst(nn.Consequent, depth+1)
		j.walkAst(nn.Alternate, depth+1)

	case *ast.ConditionalExpression:
		j.walkAst(nn.Test, depth)
		j.walkAst(nn.Consequent, depth+1)
		j.walkAst(nn.Alternate, depth+1)

	case *ast.VariableStatement:
		for _, element := range nn.List {
			j.walkAst(element, depth)
		}
	case *ast.UnaryExpression:
		j.walkerCallback(nodedata.NodeData{nodedata.UnaryExpression, nn.Operator.String(), ""}, depth)
		j.walkAst(nn.Operand, depth+1)

	case *ast.ForStatement:
		j.walkerCallback(nodedata.NodeData{nodedata.ForStatement, "", ""}, depth)
		j.walkAst(nn.Initializer, depth+1)
		j.walkAst(nn.Update, depth+1)
		j.walkAst(nn.Test, depth+1)
		j.walkAst(nn.Body, depth+1)

	case *ast.NewExpression:
		j.walkerCallback(nodedata.NodeData{nodedata.NewExpression, "", ""}, depth)
		d := j.walkAst(nn.Callee, depth+1)
		for _, arg := range nn.ArgumentList {
			j.walkAst(arg, d+1)
		}

	case *ast.WhileStatement:
		j.walkerCallback(nodedata.NodeData{nodedata.WhileStatement, "", ""}, depth)
		j.walkAst(nn.Test, depth)
		j.walkAst(nn.Body, depth+1)

	case *ast.EmptyExpression:
	case *ast.EmptyStatement:
	case nil:

	// will recur infinitely
	case *ast.BadExpression:
		j.err = fmt.Errorf("BadExpression not handled")
	case *ast.BadStatement:
		j.err = fmt.Errorf("BadStatement not handled")
	case *ast.BranchStatement:
		j.err = fmt.Errorf("BranchStatement not handled")
	case *ast.CaseStatement:
		j.err = fmt.Errorf("CaseStatement not handled")
	case *ast.CatchStatement:
		j.err = fmt.Errorf("CatchStatement not handled")
	case *ast.DebuggerStatement:
		j.err = fmt.Errorf("DebuggerStatement not handled")
	case *ast.DoWhileStatement:
		j.err = fmt.Errorf("DoWhileStatement not handled")
	case *ast.ForInStatement:
		j.err = fmt.Errorf("ForInStatement not handled")
	case *ast.LabelledStatement:
		j.err = fmt.Errorf("LabelledStatement not handled")
	case *ast.WithStatement:
		j.err = fmt.Errorf("WithStatement not handled")
	case *ast.SwitchStatement:
		j.err = fmt.Errorf("SwitchStatement not handled")
	case *ast.ThrowStatement:
		j.err = fmt.Errorf("ThrowStatement not handled")
	case *ast.TryStatement:
		j.err = fmt.Errorf("TryStatement not handled")
	default:
		j.err = fmt.Errorf("not recognized node %v", n)
	}

	return depth
}
