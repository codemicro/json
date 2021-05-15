package json

func Load(in []byte) (interface{}, error) {
	lexed, err := lex(in)
	if err != nil {
		return nil, err
	}
	return parse(lexed)
}
