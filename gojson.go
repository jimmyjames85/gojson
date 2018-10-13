package gojson

import (
	"fmt"
	"unicode/utf8"
)

// https://www.json.org/

type JsonType int

const (
	String JsonType = iota
	Number
	Object
	Array
)

type node struct {
	parent *node
}

func ParseString(b []byte) ([]byte, int, error) {

	// string
	//     '"' characters '"'

	if len(b) < 2 {
		// need at least two double quotes
		return nil, 0, fmt.Errorf("nothing to parse")
	}

	if b[0] != '"' {
		return nil, 0, fmt.Errorf("invalid char: expecting quote")
	}

	c := 1 // we've consumed the first double quote

	_, consumed := ParseCharacters(b[c:])
	c += consumed

	if len(b[c:]) == 0 {
		return nil, 0, fmt.Errorf("EOF")
	}

	if b[c:][0] != '"' {
		return nil, 0, fmt.Errorf("invalid char: expecting quote")
	}

	c += 1 // consume final quote

	return b[:c], c, nil
}

func ParseNumber(b []byte) ([]byte, int, error) {
	// number
	//     int frac exp

	_, consumed, err := ParseInt(b)
	if err != nil {
		return nil, 0, err
	}

	c := consumed
	_, consumed = ParseFrac(b[c:])
	c += consumed

	_, consumed = ParseExp(b[c:])
	c += consumed

	return b[:c], c, nil
}

func ParseInt(b []byte) ([]byte, int, error) {
	// int
	//     digit
	//     onenine digits
	//     '-' digit
	//     '-' onenine digits

	if len(b) == 0 {
		return nil, 0, fmt.Errorf("nothing to parse")
	}

	if b[0] == '0' {
		return b[:1], 1, nil // digit
	}

	if IsOneNine(b[0]) {
		_, consumed, err := ParseDigits(b)
		if err != nil {
			return b[:1], 1, nil // digit
		}
		return b[:consumed], consumed, nil // onenine digits
	}

	if b[0] != '-' {
		return nil, 0, fmt.Errorf("unexpected char")
	}

	if len(b) > 1 && b[1] == '-' {
		return nil, 0, fmt.Errorf("unexpected char")
	}

	_, consumed, err := ParseInt(b[1:])
	if err != nil {
		return nil, 0, err
	}

	consumed += 1 // the negative

	return b[0:consumed], consumed, nil
}

func ParseExp(b []byte) ([]byte, int) {

	// exp
	//     ""
	//     'E' sign digits
	//     'e' sign digits

	if len(b) == 0 ||
		(b[0] != 'e' && b[0] != 'E') {
		return nil, 0
	}

	c := 1 // we've [c]onsumed 'e'
	_, consumed := ParseSign(b[c:])

	if consumed == 0 {
		return nil, 0
	}
	c += consumed

	_, consumed, err := ParseDigits(b[c:])
	if err != nil {
		return nil, 0
	}
	c += consumed

	return b[0:c], c
}

// todo rename parse functions to consume and create interface{} Consumer... easier to test

func ParseSign(b []byte) ([]byte, int) {
	// sign
	//     ""
	//     '+'
	//     '-'

	if len(b) == 0 ||
		(b[0] != '+' && b[0] != '-') {
		return nil, 0
	}

	return b[:1], 1
}

func ParseFrac(b []byte) ([]byte, int) {
	// frac
	//     ""
	//     '.' digits

	if len(b) == 0 || b[0] != '.' {
		return nil, 0
	}
	c := 1 // we've consumed the '.'

	_, consumed, err := ParseDigits(b[c:]) // TODO: think: don't return []byte, just return how much ParseDigit consumed
	if err != nil {
		return nil, 0
	}

	c += consumed
	return b[0:c], c
}

func IsDigit(c byte) bool {
	// digit
	//     '0'
	//     onenine
	return c == '0' || IsOneNine(c)
}

func IsOneNine(c byte) bool {
	// onenine
	//     '1' . '9'
	return c >= '1' && c <= '9'
}
func ParseCharacters(b []byte) ([]byte, int) {
	// characters
	//     ""
	//     character characters

	if len(b) == 0 {
		return nil, 0
	}

	_, c, err := ParseCharacter(b)
	if err != nil {
		return nil, 0
	}

	_, consumed := ParseCharacters(b[c:])
	for consumed > 0 {
		c += consumed
		_, consumed = ParseCharacters(b[c:])
	}

	return b[:c], c
}

func ParseCharacter(b []byte) ([]byte, int, error) {
	// character
	//     '0020' . '10ffff' - '"' - '\'
	//     '\' escape

	if len(b) == 0 {
		return nil, 0, fmt.Errorf("nothing to parse")
	}

	if b[0] == '\\' { // single backslash character
		_, consumed, err := ParseEscape(b[1:])
		if err != nil {
			return nil, 0, fmt.Errorf("invalid character")
		}
		consumed += 1 // we consumed the backslash
		return b[:consumed], consumed, nil
	}

	if b[0] == '"' {
		return nil, 0, fmt.Errorf("invalid character")
	}

	// 0x10ffff overflows the length of a byte
	// So we need to extract the first rune from b
	// then verify we can verify we are within the range specified above

	r, size := utf8.DecodeRune(b)
	// r contains the first rune of the string
	// size is the size of the rune in bytes

	if r == utf8.RuneError {
		return nil, 0, fmt.Errorf("invalid char: Rune Error")
	}

	if 0x0020 <= r && r <= 0x10ffff {
		return b[:size], size, nil
	}

	return nil, 0, fmt.Errorf("invalid char")
}

func ParseEscape(b []byte) ([]byte, int, error) {
	// escape
	//     '"'
	//     '\'
	//     '/'
	//     'b'
	//     'n'
	//     'r'
	//     't'
	//     'u' hex hex hex hex

	if len(b) == 0 {
		return nil, 0, fmt.Errorf("nothing to parse")
	}

SWITCH:
	switch b[0] {
	// '\\' is single backslash character
	case '"', '\\', '/', 'b', 'n', 'r', 't':
		return b[:1], 1, nil
	case 'u':
		if len(b) < 5 {
			break // not enough hex to consume
		}
		for i := 1; i < 5; i++ {
			if !IsHex(b[i]) {
				break SWITCH
			}
		}
		return b[:5], 5, nil
	}

	return nil, 0, fmt.Errorf("Invalid escape")
}
func ParseWhitespace(b []byte) ([]byte, int) {
	// ws
	//     ""
	//     '0009' ws
	//     '000a' ws
	//     '000d' ws
	//     '0020' ws

	var i int
	var w byte

FOR:
	for i, w = range b {
		switch w {
		case 0x0009, 0x000a, 0x000d, 0x0020:
			continue
		default:
			break FOR
		}
	}

	return b[0:i], i
}

//  ParseDigits returns b[0:x] such that every ascii value from b[0]
//  to b[x] represents a digit from 0 to 9, along with the length of b[0:x]
func ParseDigits(b []byte) ([]byte, int, error) {

	// digits
	//     digit
	//     digit digits

	if len(b) == 0 {
		return nil, 0, fmt.Errorf("nothing to parse")
	}

	// check first digit
	if !IsDigit(b[0]) {
		return nil, 0, fmt.Errorf("first byte must be a digit")
	}

	// consume as many digits as possible
	for i, d := range b {
		if IsDigit(d) {
			continue
		}
		return b[0:i], i, nil
	}

	// we've consumed everything we return the whole slice back
	return b, len(b), nil
}

func IsHex(b byte) bool {
	// hex
	//     digit
	//     'A' . 'F'
	//     'a' . 'f'

	return IsDigit(b) ||
		('A' <= b && b <= 'F') ||
		('a' <= b && b <= 'f')

}

func Run() {

	return
	// var payload []byte
	// for _, c := range payload {
	// 	d, err := ParseDigit(c)

	// 	fmt.Printf("%c", c)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	fmt.Printf("[%d]", d)
	// }
}
