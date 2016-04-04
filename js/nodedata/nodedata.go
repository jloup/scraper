package nodedata

type AstNodeType uint8

const (
	UnknownAstNodeType AstNodeType = iota
	AllAstNodeType
	LeafLiteral
	ArrayLiteral
	AssignExpression
	BadExpression
	BinaryExpression
	BooleanLiteral
	BracketExpression
	CallExpression
	ConditionalExpression
	DotExpression
	EmptyExpression
	FunctionLiteral
	Identifier
	NewExpression
	NullLiteral
	NumberLiteral
	ObjectLiteral
	ParameterList
	Property
	RegExpLiteral
	SequenceExpression
	StringLiteral
	ThisExpression
	UnaryExpression
	VariableExpression
	BadStatement
	BlockStatement
	BranchStatement
	CaseStatement
	CatchStatement
	DebuggerStatement
	DoWhileStatement
	EmptyStatement
	ExpressionStatement
	ForInStatement
	ForStatement
	IfStatement
	LabelledStatement
	ReturnStatement
	SwitchStatement
	ThrowStatement
	TryStatement
	VariableStatement
	WhileStatement
	WithStatement
)

func IsLeafLiteral(t AstNodeType) bool {
	return t == StringLiteral || t == BooleanLiteral || t == NumberLiteral || t == RegExpLiteral
}

type NodeData struct {
	Type       AstNodeType
	Identifier string
	Content    string
}
