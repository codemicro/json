package json

type NodeType uint8

type Node interface {
	Type() NodeType
}

const (
	TypeObject = iota
	TypeArray
	TypeBool
	TypeInteger
	TypeFloat
	TypeNull
	TypeString
)

type Object struct {
	Value map[string]Node
}

func (*Object) Type() NodeType {
	return TypeObject
}

type Array struct {
	Value []Node
}

func (*Array) Type() NodeType {
	return TypeArray
}

type Bool struct {
	Value bool
}

func (*Bool) Type() NodeType {
	return TypeBool
}

type Integer struct {
	Value int64
}

func (*Integer) Type() NodeType {
	return TypeInteger
}

type Float struct {
	Value float64
}

func (*Float) Type() NodeType {
	return TypeFloat
}

type Null struct {}

func (*Null) Type() NodeType {
	return TypeNull
}

type String struct {
	Value string
}

func (*String) Type() NodeType {
	return TypeString
}