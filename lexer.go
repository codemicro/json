package json

import (
	"errors"
	"fmt"
	"strconv"
)

func lex(input []byte) (*tokens, error) {

	inputLength := len(input)
	var index int

	peek := func(offset int) rune {
		if index+offset >= inputLength {
			return 0
		}
		return rune(input[index+offset])
	}

	consume := func() rune {
		if index >= inputLength {
			return 0
		}
		index += 1
		return rune(input[index-1])
	}

	var (
		tokenBuf []token
		err      error
	)

	for index < inputLength {

		if isString, st, err := lexString(peek, consume); err != nil {
			return nil, err
		} else if isString {
			tokenBuf = append(tokenBuf, st)
		} else if isNumber, nt, err := lexNumber(peek, consume); err != nil {
			return nil, err
		} else if isNumber {
			tokenBuf = append(tokenBuf, nt)
		} else if isBool, bt := lexBool(peek, consume); isBool {
			tokenBuf = append(tokenBuf, bt)
		} else if isNull, nt := lexNull(peek, consume); isNull {
			tokenBuf = append(tokenBuf, nt)
		} else if isWhitespace(peek(0)) {
			consume()
		} else if isStructural(peek(0)) {
			tokenBuf = append(tokenBuf, token(consume()))
		} else {
			return nil, fmt.Errorf("unknown character: %d", peek(0))
		}

	}

	return &tokens{
		tokens: tokenBuf,
	}, err
}

func lexString(peek func(int) rune, consume func() rune) (bool, token, error) {

	// TODO: interpret escape sequences

	if peek(0) != quotationMark {
		return false, nil, nil
	}
	consume() // first quote

	var buf []rune

	for peek(0) != 0 { // run until end of input

		if peek(0) == quotationMark {
			consume()
			return true, string(buf), nil
		}

		if x := peek(0); isWhitespace(x) && x != space {
			return false, nil, errors.New("whitespace excluding U+0020 disallowed in string")
		}

		buf = append(buf, consume())
	}

	return false, nil, errors.New("string literals must end with quotation mark (U+0022)")
}

func lexNumber(peek func(int) rune, consume func() rune) (bool, token, error) {

	var buf []rune
	var numPoints int

	// TODO: parsing for e number notation
	for isNumeric(peek(0)) || peek(0) == '.' {

		if peek(0) == '.' {
			numPoints += 1
		}

		buf = append(buf, consume())
	}

	if len(buf) == 0 {
		return false, nil, nil
	}

	switch numPoints {
	case 0:
		asInteger, err := strconv.ParseInt(string(buf), 10, 64)
		if err != nil {
			return false, nil, err
		}
		return true, asInteger, nil
	case 1:
		asFloat, err := strconv.ParseFloat(string(buf), 64)
		if err != nil {
			return false, nil, err
		}
		return true, asFloat, nil
	default:
		return false, nil, errors.New("numbers cannot have more than one point")
	}
}

func lexBool(peek func(int) rune, consume func() rune) (bool, token) {

	if peekSequence(peek, 0, len(trueLiteral)) == trueLiteral {
		_ = consumeSequence(consume, len(trueLiteral))
		return true, true
	} else if peekSequence(peek, 0, len(falseLiteral)) == falseLiteral {
		_ = consumeSequence(consume, len(falseLiteral))
		return true, false
	}

	return false, nil
}

func lexNull(peek func(int) rune, consume func() rune) (bool, token) {

	if peekSequence(peek, 0, len(nullLiteral)) == nullLiteral {
		consumeSequence(consume, len(nullLiteral))
		return true, nil
	}

	return false, nil
}
