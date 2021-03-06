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
			return nil, fmt.Errorf("unknown character: %s", string(peek(0)))
		}

	}

	return &tokens{
		tokens: tokenBuf,
	}, err
}

func lexString(peek func(int) rune, consume func() rune) (bool, token, error) {

	if peek(0) != quotationMark {
		return false, nil, nil
	}
	consume() // first quote

	var buf []rune

	for peek(0) != 0 { // run until end of input

		if peek(0) == quotationMark {
			// end of string
			consume()
			return true, string(buf), nil
		}

		if x := peek(0); isWhitespace(x) && x != space {
			return false, nil, errors.New("whitespace excluding U+0020 disallowed in string")
		}

		// string escape sequences
		var char rune
		switch peekSequence(peek, 0, 2) {
		case `\"`:
			char = '"'
			consumeSequence(consume, 2)
		case `\\`:
			char = '\\'
			consumeSequence(consume, 2)
		case `\/`:
			char = '/'
			consumeSequence(consume, 2)
		case `\b`:
			char = '\b'
			consumeSequence(consume, 2)
		case `\f`:
			char = '\f'
			consumeSequence(consume, 2)
		case `\n`:
			char = '\n'
			consumeSequence(consume, 2)
		case `\r`:
			char = '\r'
			consumeSequence(consume, 2)
		case `\t`:
			char = '\t'
			consumeSequence(consume, 2)
		case `\u`:
			nextFour := peekSequence(peek, 2, 4)
			parsed64, err := strconv.ParseInt(nextFour, 16, 32)
			if err != nil {
				return false, nil, errors.New("invalid unicode escape sequence: " + err.Error())
			}
			char = rune(parsed64)
			consumeSequence(consume, 6)
		}

		if char == 0 {
			// if no escape sequence set
			char = consume()
		}

		buf = append(buf, char)
	}

	return false, nil, errors.New("string literals must end with quotation mark (U+0022)")
}

func lexNumber(peek func(int) rune, consume func() rune) (bool, token, error) {

	var buf []rune
	var (
		numPoints int
		numEs     int
	)

	for isNumeric(peek(0)) {

		if peek(0) == '.' {
			numPoints += 1
		} else if isE(peek(0)) {
			numEs += 1
		}

		buf = append(buf, consume())
	}

	if len(buf) == 0 {
		return false, nil, nil
	}

	if buf[0] == '+' {
		return false, nil, errors.New("unexpected +")
	}

	// check the first digit to see if it's a zero
	var firstDigit rune
	var nextDigit rune
	for _, char := range buf {
		if isDigit(char) || char == '.' {
			if firstDigit == 0 {
				firstDigit = char
			} else {
				nextDigit = char
				break
			}
		}
	}

	if firstDigit == '0' && !(nextDigit == 0 || nextDigit == '.') {
		return false, nil, errors.New("numbers cannot have zero as their first digit unless the value is zero or the number is a decimal")
	}

	if numPoints > 1 {
		return false, nil, errors.New("numbers cannot have more than one point")
	} else if numEs > 1 {
		return false, nil, errors.New("numbers cannot have more than one E")
	}

	if numEs == 0 && numPoints == 0 {
		asInteger, err := strconv.ParseInt(string(buf), 10, 64)
		if err != nil {
			return false, nil, err
		}
		return true, asInteger, nil
	} else {

		// reject sequences that end with a point
		if buf[len(buf)-1] == '.' {
			return false, nil, errors.New("numbers cannot end with a point")
		}

		// reject sequences that start with a point
		if buf[0] == '.' {
			return false, nil, errors.New("numbers cannot start with a point")
		}

		asFloat, err := strconv.ParseFloat(string(buf), 64)
		if err != nil {
			return false, nil, err
		}
		return true, asFloat, nil
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
