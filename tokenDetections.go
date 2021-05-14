package json

func isValue(t token) bool {
	// assuming all language tokens are defined as runes
	_, isRune := t.(rune)
	return !isRune
}
