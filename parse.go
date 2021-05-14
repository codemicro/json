package json

import "errors"

type nmt struct{}

var noMoreTokens = nmt{}

func parse(t *tokens) (interface{}, error) {
	if t.peek(0) == leftSqBracket {
		return parseArray(t)
	} else if t.peek(0) == leftCurlBracket {
		return parseObject(t)
	} else if len(t.tokens) > 1 {
		return nil, errors.New("extra data")
	} else {
		return t.consume(), nil
	}
}

func parseArray(t *tokens) (interface{}, error) {

	var values []interface{}

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

		if !isValue(t.peek(0)) {
			return nil, errors.New("expected value")
		}

		parsedValue, err := parse(t)
		if err != nil {
			return nil, err
		}

		values = append(values, parsedValue)

		if t.peek(0) == comma {
			t.consume()
		}

	}

	return nil, errors.New("expected end of array")
}

func parseObject(t *tokens) (interface{}, error) {

	values := make(map[string]interface{})

	if t.peek(0) != leftCurlBracket {
		return nil, errors.New("objects must begin with }")
	}
	t.consume()

	for t.peek(0) != noMoreTokens {

		if t.peek(0) == rightCurlBracket {
			t.consume()
			return values, nil
		}

		if x := t.peek(-1); !(x == comma || x == rightCurlBracket) {
			// this call to t.peek with -1 is only safe to do since we know we have at least 1 item in the token array
			// on account of the initial check
			return nil, errors.New("unexpected token")
		}

		key := t.peek(0)
		sep := t.peek(1)
		val := t.peek(2)

		var (
			ok bool
			stringKey string
		)
		if stringKey, ok = key.(string); !ok {
			return nil, errors.New("object keys must be strings")
		}

		if sep != colon {
			return nil, errors.New("expected colon")
		}

		if !isValue(val) {
			return nil, errors.New("expected value")
		}

		values[stringKey] = val

		for i := 0; i < 3; i += 1 {
			t.consume()
		}

	}

	return nil, errors.New("expected end of object")
}