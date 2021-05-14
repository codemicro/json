package json

type token interface{}

type tokens struct {
	tokens []token
	index int
}

func (t *tokens) peek(offset int) token {
	if t.index+offset >= len(t.tokens) {
		return 0
	}
	return t.tokens[t.index+offset]
}

func (t *tokens) consume() token {
	if t.index >= len(t.tokens) {
		return 0
	}
	t.index += 1
	return t.tokens[t.index-1]
}