package json

func isWhitespace(r rune) bool {
	return r == tabulation ||
		r == carriageReturn ||
		r == lineFeed ||
		r == space
}

func isStructural(r rune) bool {
	return r == leftSqBracket ||
		r == rightSqBracket ||
		r == leftCurlBracket ||
		r == rightCurlBracket ||
		r == colon ||
		r == comma
}

func isNumeric(r rune) bool {
	return '0' <= r && r <= '9'
}

func peekSequence(peek func(int) rune, from, length int) string {
	var buf []rune
	for i := 0; i  < length; i += 1 {
		x := peek(from + i)
		if x != 0 {
			buf = append(buf, x)
		}
	}
	return string(buf)
}

func consumeSequence(consume func() rune, length int) string {
	var buf []rune
	for i := 0; i < length; i += 1 {
		x := consume()
		if x != 0 {
			buf = append(buf, x)
		}
	}
	return string(buf)
}