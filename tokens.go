package json

import "fmt"

type token interface{}

func TokenAsNode(t token) Node {

	switch t := t.(type) {
	case rune:
		panic("rune token cannot be converted into Node")
	case string:
		return &String{Value: t}
	case int64:
		return &Integer{Value: t}
	case float64:
		return &Float{Value: t}
	case bool:
		return &Bool{Value: t}
	case nmt:
		return &Null{}
	default:
		panic(fmt.Sprintf("unknown type %T", t))
	}
}

type tokens struct {
	tokens []token
	index  int
}

type nmt struct{}

var noMoreTokens = nmt{}

func (t *tokens) peek(offset int) token {
	if t.index+offset >= len(t.tokens) {
		return noMoreTokens
	}
	return t.tokens[t.index+offset]
}

func (t *tokens) consume() token {
	if t.index >= len(t.tokens) {
		return noMoreTokens
	}
	t.index += 1
	return t.tokens[t.index-1]
}
