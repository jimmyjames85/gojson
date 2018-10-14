package gojson

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode/utf8"
)

// https://www.json.org/

type JsonType int

const (
	TypeString JsonType = iota
	TypeNumber
	TypeObject
	TypeFrac
	TypeArray
	TypeElement
	TypeElements
	TypeExp
	TypeSign
	TypeDigits
	TypeInt
	TypeWhitespace
	TypeEscape
)

var (
	ErrNothingToParse            = fmt.Errorf("nothing to parse")
	ErrInvalidCharacter          = fmt.Errorf("invalid character")
	ErrInvalidDigit              = fmt.Errorf("invalid digit")
	ErrInvalidEscape             = fmt.Errorf("invalid escape")
	ErrInvalidCharacterRuneError = fmt.Errorf("invalid character: rune error")
	ErrInvalidObjectOpen         = fmt.Errorf("invalid object: expecting '{'")
	ErrInvalidObjectClose        = fmt.Errorf("invalid object: expecting '}'")
	ErrInvalidMemberMissingSep   = fmt.Errorf("invalid member: expecting ':'")
	ErrInvalidArrayOpen          = fmt.Errorf("invalid array: expecting '['")
	ErrInvalidArrayClose         = fmt.Errorf("invalid array: expecting ']'")
	ErrInvalidStringOpen         = fmt.Errorf(`invalid string: missing beginnig '"'`)
	ErrInvalidStringClose        = fmt.Errorf(`invalid string: missing ending '"'`)
	ErrUnexpectedChar            = fmt.Errorf("unexpected char")
	ErrParseInteger              = fmt.Errorf("parse error: not an integer")

	ErrUnsupported = fmt.Errorf("unsupported: should we panic") // todo

	TrueValue  = []byte(`true`)
	FalseValue = []byte(`false`)
	NullValue  = []byte(`null`)
)

type Node struct {
	Type     JsonType
	Parent   *Node
	Children []*Node

	b []byte // underlying data
	// Value  fmt.Stringer
}

func ParseJSON(b []byte) ([]byte, int, error) {
	// TODO add unit test for this...

	// json
	//     element

	_, c, err := ParseElement(b)
	if err != nil {
		return nil, 0, err
	}

	return b[:c], c, nil

}

func ParseValue(b []byte) ([]byte, int, error) {
	// TODO add unit test for this...

	if len(b) == 0 {
		return nil, 0, ErrNothingToParse
	}

	if bytes.HasPrefix(b, TrueValue) {
		// "true"
		return b[:len(TrueValue)], len(TrueValue), nil
	} else if bytes.HasPrefix(b, FalseValue) {
		// "false"
		return b[:len(FalseValue)], len(FalseValue), nil
	} else if bytes.HasPrefix(b, NullValue) {
		// "null"
		return b[:len(NullValue)], len(NullValue), nil
	}

	_, c, err := ParseNumber(b)
	if err == nil {
		return b[:c], c, nil
	}

	_, c, err = ParseString(b)
	if err == nil {
		return b[:c], c, nil
	}

	_, c, err = ParseArray(b)
	if err == nil {
		return b[:c], c, nil
	}

	// TODO this isn't preformant we attempt each type and rescan on failures
	_, c, err = ParseObject(b)
	if err == nil {
		return b[:c], c, nil
	}

	return nil, 0, ErrUnsupported
}

func ParseObject(b []byte) ([]byte, int, error) {
	// object
	//     '{' ws '}'
	//     '{' members '}'

	// copy/pasta from ParseArray
	if len(b) == 0 || b[0] != '{' {
		return nil, 0, ErrInvalidObjectOpen
	}
	c := 1 // consume the '{'

	// attempt to consume whitspace a check for closing '}'
	_, consumed := ParseWhitespace(b[c:])
	c += consumed // consume whitespace
	if len(b[c:]) > 0 && b[c:][0] == '}' {
		c++ // consume the '}'
		// we have empty object
		return b[:c], c, nil
	}
	c -= consumed // unconsume whitespace and let it be part of elements

	_, consumed, err := ParseMembers(b[c:])
	if err != nil {
		return nil, 0, err
	}
	c += consumed

	if len(b[c:]) > 0 && b[c:][0] == '}' {
		c++ // consume the '}'
		// we have non-empty object
		return b[:c], c, nil
	}

	return nil, 0, ErrInvalidObjectClose
}

func ParseMembers(b []byte) ([]byte, int, error) {
	// members
	//     member
	//     member ',' members

	// this is copy/pasta from ParseElements

	_, consumed, err := ParseMember(b)
	if err != nil {
		// must parse at least one member
		return nil, 0, err
	}

	c := consumed
	// noMoreToConsume || next unconsumed byte is not ','
	if len(b[c:]) == 0 || b[c:][0] != ',' {
		return b[:c], c, nil
	}
	c++ // consume the ','

	// TODO what is the recusion limit in go? and will it limit
	// how many members in a json 'list' we can support
	_, consumed, err = ParseMembers(b[c:])
	if err != nil {
		c-- // don't unconsume the ',' just the first member
		return b[:c], c, nil
	}

	c += consumed
	return b[:c], c, nil
}

func ParseMember(b []byte) ([]byte, int, error) {
	// member
	//     ws string ws ':' element

	_, consumed := ParseWhitespace(b)
	c := consumed

	_, consumed, err := ParseString(b[c:])
	if err != nil {
		return nil, 0, err
	}
	c += consumed

	_, consumed = ParseWhitespace(b[c:])
	c += consumed

	// noMoreToConsume || next unconsumed byte is not ':'  TODO make this if stmt a func... if !nextCharIs(b[c:], ':')
	if len(b[c:]) == 0 || b[c:][0] != ':' {
		return nil, 0, ErrInvalidMemberMissingSep
	}
	c += 1 // consume the ':'

	_, consumed, err = ParseElement(b[c:])
	if err != nil {
		return nil, 0, err
	}
	c += consumed

	return b[:c], c, nil
}

func ParseArray(b []byte) ([]byte, int, error) {
	// array
	//     '[' ws ']'
	//     '[' elements ']'

	if len(b) == 0 || b[0] != '[' {
		return nil, 0, ErrInvalidArrayOpen
	}
	c := 1 // consume the '['

	// attempt to consume whitspace a check for closing ']'
	_, consumed := ParseWhitespace(b[c:])
	c += consumed // consume whitespace
	if len(b[c:]) > 0 && b[c:][0] == ']' {
		c++ // consume the ']'
		// we have empty array
		return b[:c], c, nil
	}
	c -= consumed // unconsume whitespace and let it be part of elements

	_, consumed, err := ParseElements(b[c:])
	if err != nil {
		return nil, 0, err
	}
	c += consumed

	if len(b[c:]) > 0 && b[c:][0] == ']' {
		c++ // consume the ']'
		// we have non-empty array
		return b[:c], c, nil
	}

	return nil, 0, ErrInvalidArrayClose
}

type Elements struct{ Node }

func ParseElements(b []byte) (*Elements, int, error) {
	// elements
	//     element
	//     element ',' elements

	ret := &Elements{Node: Node{
		Type:   TypeElements,
		Parent: nil, //todo
	}}

	_, consumed, err := ParseElement(b)
	if err != nil {
		// must parse at least one element
		return ret, 0, err
	}

	c := consumed
	// noMoreToConsume || next unconsumed byte is not ','
	if len(b[c:]) == 0 || b[c:][0] != ',' {
		ret.b = b[:c]
		return ret, c, nil
	}
	c++ // consume the ','

	// TODO what is the recusion limit in go? and will it limit
	// how many elements in a json 'list' we can support
	_, consumed, err = ParseElements(b[c:])
	if err != nil {
		c-- // unconsume the ','
		ret.b = b[:c]
		return ret, c, nil
	}

	c += consumed

	ret.b = b[:c]
	return ret, c, nil
}

type Element struct{ Node }

func ParseElement(b []byte) (*Element, int, error) {
	// element
	//     ws value ws

	ret := &Element{Node: Node{
		Type:   TypeElement,
		Parent: nil, // todo
	}}

	_, c := ParseWhitespace(b)

	_, consumed, err := ParseValue(b[c:])
	if err != nil {
		return ret, 0, err
	}
	c += consumed

	_, consumed = ParseWhitespace(b[c:])
	c += consumed

	ret.b = b[:c]
	return ret, c, nil
}

type String struct{ Node }

func ParseString(b []byte) (*String, int, error) {
	// string
	//     '"' characters '"'

	ret := &String{Node: Node{
		Type:   TypeString,
		Parent: nil, // todo
	}}

	if len(b) < 2 {
		// need at least two double quotes
		return ret, 0, ErrNothingToParse
	}

	if b[0] != '"' {
		return ret, 0, ErrInvalidStringOpen
	}

	c := 1 // we've consumed the first double quote

	_, consumed := ParseCharacters(b[c:])
	c += consumed

	// noMoreToConsume
	if len(b[c:]) == 0 {
		return ret, 0, ErrNothingToParse
	}

	// next unconsumed byte is not '"'
	if b[c:][0] != '"' {
		return ret, 0, ErrInvalidStringClose
	}

	c += 1 // consume final quote

	ret.b = b[:c]
	return ret, c, nil
}

type Number struct{ Node }

func (n *Number) Int() (int64, error) {
	d := bytes.IndexByte(n.b, '.')
	if d == -1 {
		return strconv.ParseInt(string(n.b), 10, 64)
	}
	return 0, ErrParseInteger // TODO is this even being used?
}

func (n *Number) Float64() float64 {
	ret, err := strconv.ParseFloat(string(n.b), 64)
	if err != nil {
		// then we did not parse digits correctly
		panic(err) // todo how do return message that makes sense, when you ParsedDigits you got an error that you ignored
	}
	return ret
}

func ParseNumber(b []byte) (*Number, int, error) {
	// number
	//     int frac exp
	//
	// Note int is required, but frac and exp can be empty strings

	ret := &Number{Node: Node{
		Type:   TypeNumber,
		Parent: nil, // todo
	}}
	_, consumed, err := ParseInt(b)
	if err != nil {
		return ret, 0, err
	}

	c := consumed
	_, consumed = ParseFrac(b[c:])
	c += consumed

	_, consumed = ParseExp(b[c:])
	c += consumed

	ret.b = b[:c]

	return ret, c, nil
}

type Int struct{ Node }

func (i *Int) Int() int64 {
	ret, err := strconv.ParseInt(string(i.b), 10, 64)
	if err != nil {
		// then we did not parse digits correctly
		panic(err) // or better safe then sorry?... return 0
		// todo how do return message that makes sense, when you ParsedDigits you got an error that you ignored
	}
	return ret
}

func ParseInt(b []byte) (*Int, int, error) {
	// int
	//     digit
	//     onenine digits
	//     '-' digit
	//     '-' onenine digits

	ret := &Int{Node: Node{
		Type:   TypeInt,
		Parent: nil, // todo

	}}

	if len(b) == 0 {
		return ret, 0, ErrNothingToParse
	}

	if b[0] == '0' {
		ret.b = b[:1]
		return ret, 1, nil // digit
	}

	if IsOneNine(b[0]) {
		_, consumed, err := ParseDigits(b)
		if err != nil {
			ret.b = b[:1]
			return ret, 1, nil // digit
		}
		ret.b = b[:consumed]      //todo use c not consumed
		return ret, consumed, nil // onenine digits
	}

	if b[0] != '-' {
		return ret, 0, ErrUnexpectedChar
	}

	if len(b) > 1 && b[1] == '-' {
		return ret, 0, ErrUnexpectedChar
	}

	_, consumed, err := ParseInt(b[1:])
	if err != nil {
		return ret, 0, err
	}

	consumed += 1 // the negative

	ret.b = b[:consumed]
	return ret, consumed, nil
}

type Exp struct {
	Node
	sign *Sign
	// TODO Digits
}

// todo this is sign of the exponent not the number
func (e *Exp) Positive() bool { return e.sign == nil || e.sign.Positive() }

func ParseExp(b []byte) (*Exp, int) {

	// exp
	//     ""
	//     'E' sign digits
	//     'e' sign digits

	ret := &Exp{
		Node: Node{
			Type:   TypeExp,
			Parent: nil, // todo

		},
	}

	if len(b) == 0 ||
		(b[0] != 'e' && b[0] != 'E') {
		return ret, 0
	}

	c := 1 // we've [c]onsumed 'e'
	sign, consumed := ParseSign(b[c:])
	ret.sign = sign // todo why can't i do ret.Sign, consumed := ParseSign(b[c:])

	if consumed == 0 {
		return ret, 0
	}
	c += consumed

	_, consumed, err := ParseDigits(b[c:])
	if err != nil {
		return ret, 0
	}
	c += consumed

	ret.b = b[0:c]
	return ret, c
}

// todo rename parse functions to consume and create interface{} Consumer... easier to test
type Sign struct{ Node }

func (s *Sign) Positive() bool { return len(s.b) != 1 || s.b[0] == '+' }

func ParseSign(b []byte) (*Sign, int) {
	// sign
	//     ""
	//     '+'
	//     '-'

	ret := &Sign{Node: Node{
		Type:   TypeSign,
		Parent: nil, // todo
	}}

	if len(b) == 0 ||
		(b[0] != '+' && b[0] != '-') {
		return ret, 0
	}

	ret.b = b[:1]

	return ret, 1
}

type Frac struct{ Node }

func (f *Frac) Float64() float64 {
	// ParseFrac doesn't return error but can return b=empty string
	// so this should return 0 for the fraction part
	if len(f.b) == 0 {
		return float64(0)
	}

	ret, err := strconv.ParseFloat(string(f.b), 64)
	if err != nil {
		// then we did not parse float correctly
		panic(err)
	}

	return ret
}

func ParseFrac(b []byte) (*Frac, int) {
	// frac
	//     ""
	//     '.' digits

	ret := &Frac{Node: Node{
		Type:   TypeFrac,
		Parent: nil, // todo
	}}

	if len(b) == 0 || b[0] != '.' {
		return ret, 0
	}
	c := 1 // we've consumed the '.'

	_, consumed, err := ParseDigits(b[c:]) // TODO: think: don't return []byte, just return how much ParseDigit consumed
	if err != nil {
		return ret, 0
	}

	c += consumed
	ret.b = b[0:c]
	return ret, c
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

// type Character struct{ Node }

// I've made sure that all Parse functions that return *Something doesn't return nil
// This is a contract that should be documented somewhere. But I wonder, should we
// return non-pointer... need to see impact on heap vs stack memory allocation, garbage collection, speed
func ParseCharacter(b []byte) ([]byte, int, error) {
	// character
	//     '0020' . '10ffff' - '"' - '\'
	//     '\' escape

	// most time is spent inside this function so we should avoid mallocs

	if len(b) == 0 {
		return nil, 0, ErrNothingToParse
	}

	if b[0] == '\\' { // single backslash character
		_, consumed, err := ParseEscape(b[1:])
		if err != nil {
			return nil, 0, ErrInvalidCharacter
		}
		consumed += 1 // we consumed the backslash
		return b[:consumed], consumed, nil
	}

	if b[0] == '"' {
		return nil, 0, ErrInvalidCharacter
	}

	// 0x10ffff overflows the length of a byte
	// So we need to extract the first rune from b
	// then verify we can verify we are within the range specified above

	r, size := utf8.DecodeRune(b)
	// r contains the first rune of the string
	// size is the size of the rune in bytes

	if r == utf8.RuneError {
		return nil, 0, ErrInvalidCharacterRuneError
	}

	if 0x0020 <= r && r <= 0x10ffff {
		return b[:size], size, nil
	}

	return nil, 0, ErrInvalidCharacter
}

type Escape struct{ Node }

func ParseEscape(b []byte) (*Escape, int, error) {
	// escape
	//     '"'
	//     '\'
	//     '/'
	//     'b'
	//     'n'
	//     'r'
	//     't'
	//     'u' hex hex hex hex

	ret := &Escape{Node: Node{
		Type:   TypeEscape,
		Parent: nil, // todo
	}}

	if len(b) == 0 {
		return ret, 0, ErrNothingToParse
	}

SWITCH:
	switch b[0] {
	// '\\' is single backslash character
	case '"', '\\', '/', 'b', 'n', 'r', 't':
		ret.b = b[:1]
		return ret, 1, nil
	case 'u':
		if len(b) < 5 {
			break // not enough hex to consume
		}
		for i := 1; i < 5; i++ {
			if !IsHex(b[i]) {
				break SWITCH
			}
		}
		ret.b = b[:5]
		return ret, 5, nil
	}

	return ret, 0, ErrInvalidEscape
}

type Whitespace struct{ Node }

func ParseWhitespace(b []byte) (*Whitespace, int) {
	// ws
	//     ""
	//     '0009' ws
	//     '000a' ws
	//     '000d' ws
	//     '0020' ws

	ret := &Whitespace{Node: Node{
		Type:   TypeWhitespace,
		Parent: nil, // todo
	}}

	var n int

FOR:
	for _, ws := range b {
		switch ws {
		case 0x0009, 0x000a, 0x000d, 0x0020:
			n++
			continue
		default:
			break FOR
		}
	}

	// TODO for range loop doesn't incrememnt on the last loop?
	// for i, ws := range b{} if len(b)== 10 then i should be 10?

	ret.b = b[:n]
	return ret, n
}

// [wip] Digits satisfies the json.org language spec for both digit and digits

type Digits struct{ Node }

func (d *Digits) Int() uint64 {
	ret, err := strconv.ParseUint(string(d.b), 10, 64)
	if err != nil {
		// then we did not parse digits correctly
		panic(err) // or better safe then sorry?... return 0
		// todo how do return message that makes sense, when you ParsedDigits you got an error that you ignored
	}
	return ret
}

//  ParseDigits returns b[0:x] such that every ascii value from b[0]
//  to b[x] represents a digit from 0 to 9, along with the length of b[0:x]
func ParseDigits(b []byte) (*Digits, int, error) {

	// digits
	//     digit
	//     digit digits

	ret := &Digits{Node: Node{
		Type:   TypeDigits,
		Parent: nil, // todo
	}}

	if len(b) == 0 {
		return ret, 0, ErrNothingToParse
	}

	// check first digit
	if !IsDigit(b[0]) {
		return ret, 0, ErrInvalidDigit
	}

	// consume as many digits as possible
	for i, d := range b {
		if IsDigit(d) {
			continue
		}
		ret.b = b[0:i]
		return ret, i, nil
	}

	// we've consumed everything we return the whole slice back
	ret.b = b
	return ret, len(b), nil
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
