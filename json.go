package json

import "errors"

func Load(in []byte) (Node, error) {
	lexed, err := lex(in)
	if err != nil {
		return nil, err
	}
	parsed, err := parse(lexed)
	if err != nil {
		return nil, err
	}
	if lexed.peek(0) != noMoreTokens {
		return nil, errors.New("extra data")
	}
	return parsed, nil
}
