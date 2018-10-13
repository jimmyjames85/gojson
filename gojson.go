package gojson

import "fmt"
import "bytes"

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

func parseString() {}

func parseNumber() {

}

func parseObject() {}

func parseArray() {}

// NOTE if empty string satisfies a json type then an error should not
// be returned

func ParseExp(b []byte) ([]byte, int) {

	// exp
	//     ""
	//     'E' sign digits
	//     'e' sign digits

	if len(b) == 0 {
		return nil, 0
	}

	if b[0] != 'e' && b[0] != 'E' {
		return nil, 0
	}

	c := 1 // we've [c]onsumed 'e'
	sign, consumed := ParseSign(b[c:])

	if consumed == 0 {
		return nil, 0
	}
	c += consumed

	digits, consumed, err := ParseDigits(b[c:])
	if err != nil {
		return nil, 0
	}
	c += consumed

	fmt.Printf("%c %s\n", sign, string(digits))

	return b[0:c], c

}

//  ParseSign returns '+' or '-' and 1, or return zero byte and 0
func ParseSign(b []byte) (byte, int) {
	// sign
	//     ""
	//     '+'
	//     '-'

	if len(b) == 0 {
		return 0, 0
	}

	switch b[0] {
	case '+', '-':
		return b[0], 1
	default:
		return 0, 0
	}
}

//  ParseDigit returns b if the ascii value of b represents the a
//  digit from 0 to 9, otherwise returns an error
func ParseDigit(c byte) (byte, error) {
	d := int(c - 48)
	if d >= 0 && d <= 9 {
		return c, nil
	}
	return 0, fmt.Errorf("not a digit")
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

//  ParseDigits returns b[0:x] such that every ascii value from b[0]
//  to b[x] represents a digit from 0 to 9, along with the length of b[0:x]
func ParseDigits(b []byte) ([]byte, int, error) {

	// todo this is wrong - should through error e.g. on non-digit "dsa"

	// digits
	//     digit
	//     digit digits

	for i, c := range b {
		_, err := ParseDigit(c)
		if err != nil {
			return b[0:i], i, nil
		}
	}
	return nil, 0, fmt.Errorf("expecting something but what?")
}
func ParseInt(b []byte) ([]byte, int, error) {
	// int
	//     digit
	//     onenine digits
	//     '-' digit
	//     '-' onenine digits

	return nil, 0, fmt.Errorf("expecting something but what?")
}

func ParseHex(b []byte) (byte, error) {
	// hex
	//     digit
	//     'A' . 'F'
	//     'a' . 'f'

	if len(b) == 0 {
		return 0, fmt.Errorf("expecting something but what?")
	}

	d, err := ParseDigit(b[0])
	if err == nil {
		return d, nil
	}

	// todo this can be faster
	if bytes.Contains([]byte("ABCDEFabcdef"), []byte{b[0]}) {
		return b[0], nil
	}

	return 0, fmt.Errorf("expecting something but what?")
}

func _TestParseDigits() {
	digits := []byte("2312012sada")

	ds, length, err := ParseDigits(digits)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d characters parsed\ndigits are: %s\n", length, string(ds))
}

func _TestParseSign() {
	byts := []byte("-2312012sada")

	sign, length := ParseSign(byts)
	fmt.Printf("%d characters parsed\nsign is: %s\n", length, string(sign))
}

func _TestParseExp() {
	byts := []byte("e-2312012sada")

	sign, length := ParseExp(byts)
	fmt.Printf("%d characters parsed\nexp is: %s\n", length, string(sign))
}

func _TestParseFrac() {
	byts := []byte(".234324asdf-2312012sada")

	frac, length := ParseFrac(byts)
	fmt.Printf("%d characters parsed\nfrac is: %s\n", length, string(frac))
}

func _TestParseHex() {
	byts := []byte("F2312012sada")
	hex, err := ParseHex(byts)
	fmt.Printf("err: %s\nhex is: %s\n", err, string(hex))
}

func Run(payload []byte) {
	_TestParseDigits()

	return

	for _, c := range payload {
		d, err := ParseDigit(c)

		fmt.Printf("%c", c)
		if err != nil {
			continue
		}
		fmt.Printf("[%d]", d)
	}
}
