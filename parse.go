package json

import (
	"errors"
	"fmt"
)

// parse a single value
func parse(t *tokens) (Node, error) {
	currentToken := t.peek(0)
	if currentToken == leftSqBracket {
		return parseArray(t)
	} else if currentToken == leftCurlBracket {
		return parseObject(t)
	} else {
		return TokenAsNode(t.consume()), nil
	}
}

func parseArray(t *tokens) (Node, error) {

	values := &Array{Value: []Node{}}

	if t.peek(0) != leftSqBracket {
		return nil, errors.New("arrays must begin with [")
	}
	t.consume()

	for t.peek(0) != noMoreTokens {

		// TODO: we probably don't need a comma check in two places

		if t.peek(0) == rightSqBracket {
			t.consume() // consume right square bracket
			return values, nil
		}

		if x := t.peek(-1); !(x == comma || x == leftSqBracket) {
			// this call to t.peek with -1 is only safe to do since we know we have at least 1 item in the token array
			// on account of the initial check
			return nil, errors.New("unexpected token")
		}

		parsedValue, err := parse(t)
		if err != nil {
			return nil, err
		}

		values.Value = append(values.Value, parsedValue)

		if t.peek(0) == comma {
			t.consume()
		}

	}

	return nil, errors.New("expected end of array")
}

func parseObject(t *tokens) (Node, error) {
	values := &Object{
		Value: map[string]Node{},
	}

	if t.peek(0) != leftCurlBracket {
		return nil, errors.New("objects must begin with }")
	}
	t.consume()

	for t.peek(0) != noMoreTokens {

		if t.peek(0) == rightCurlBracket {
			t.consume()
			return values, nil
		}

		fmt.Println(t.peek(-1), t.peek(0))
		if x := t.peek(-1); !(x == comma || x == leftCurlBracket) {
			// this call to t.peek with -1 is only safe to do since we know we have at least 1 item in the token array
			// on account of the initial check
			return nil, errors.New("unexpected token")
		}

		key := t.peek(0)
		sep := t.peek(1)
		val := t.peek(2)

		var (
			ok        bool
			stringKey string
		)
		if stringKey, ok = key.(string); !ok {
			return nil, errors.New("object keys must be strings")
		}

		if sep != colon {
			return nil, errors.New("expected colon")
		}

		if valRune, ok := val.(rune); ok && isStructural(valRune) && valRune != leftSqBracket && valRune != leftCurlBracket {
			fmt.Println(string(valRune))
			return nil, errors.New("expected value")
		}

		for i := 0; i < 2; i += 1 {
			t.consume()
		}

		parsedValue, err := parse(t)
		if err != nil {
			return nil, err
		}

		values.Value[stringKey] = parsedValue
	}

	return nil, errors.New("expected end of object")
}
