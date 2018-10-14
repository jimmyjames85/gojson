package gojson

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode/utf8"
)

// https://www.json.org/

type JsonType int

var (
	ErrNothingToParse            = fmt.Errorf("nothing to parse")
	ErrInvalidCharacter          = fmt.Errorf("invalid character")
	ErrInvalidDigit              = fmt.Errorf("invalid digit")
	ErrInvalidNull               = fmt.Errorf("invalid null")
	ErrInvalidBoolean            = fmt.Errorf("invalid boolean")
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

type Value []byte

// type value needs to indicate object array string number true, false, or null
func ParseValue(b []byte) (Value, int, error) {
	// value
	//     object
	//     array
	//     string
	//     number
	//     "true"
	//     "false"
	//     "null"

	// TODO add unit test for this...

	if len(b) == 0 {
		return nil, 0, ErrNothingToParse
	}

	// TODO is this preformant? we attempt each type and rescan on
	// failures, see if we can guess which Parse func to call

	// TODO order of Parse func matters, e.g. I assume null is not
	// used frequently so we call it last

	_, c, err := ParseObject(b)
	if err == nil {
		return b[:c], c, nil
	}

	_, c, err = ParseArray(b)
	if err == nil {
		return b[:c], c, nil
	}

	_, c, err = ParseString(b)
	if err == nil {
		return b[:c], c, nil
	}

	_, c, err = ParseNumber(b)
	if err == nil {
		return b[:c], c, nil
	}

	_, c, err = ParseBoolean(b)
	if err == nil {
		return b[:c], c, nil
	}

	_, c, err = ParseNull(b)
	if err == nil {
		return b[:c], c, nil
	}

	return nil, 0, ErrUnsupported
}

// if a value is null it doesn't give us any information about the type e.g. it could be an array or an object
//
// null satisfies both object and array, hence it might be difficult for us to glean the intended underlying type
type Null []byte

func ParseNull(b []byte) (Null, int, error) {
	// null is not a json.org 'named' value (such as string or number), but satisfies the "null" value
	//
	// null
	//     "null"

	if len(b) >= len(NullValue) && bytes.Equal(b[0:len(NullValue)], NullValue) {
		// "null"
		return b[:len(NullValue)], len(NullValue), nil
	}
	return nil, 0, ErrInvalidNull
}

type Boolean []byte

func ParseBoolean(b []byte) (Boolean, int, error) {
	// boolean is not a json.org 'named' value (such as string or number), but satisfies the "false" and/or "true" value
	//
	// boolean
	//     "true"
	//     "false"

	if len(b) >= len(TrueValue) && bytes.Equal(b[0:len(TrueValue)], TrueValue) {
		// "true"
		return b[:len(TrueValue)], len(TrueValue), nil
	} else if len(b) >= len(FalseValue) && bytes.Equal(b[0:len(FalseValue)], FalseValue) {
		// "false"
		return b[:len(FalseValue)], len(FalseValue), nil
	}

	return nil, 0, ErrInvalidBoolean
}

type Object Value // todo map[string]Element

func ParseObject(b []byte) (Object, int, error) {
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

// type members   ... this should be map[string] Value
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

	//  while there's moreToConsume && we the expected delimeter ',' is there...
	for len(b[c:]) > 0 && b[c:][0] == ',' {
		c++ // consume the ','
		_, consumed, err = ParseMember(b[c:])
		if err != nil {
			c-- // unconsume the last ','
			break
		}
		c += consumed
	}
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

type Elements []byte // todo this should eventually be []Value

func ParseElements(b []byte) (Elements, int, error) {
	// elements
	//     element
	//     element ',' elements

	_, consumed, err := ParseElement(b)
	if err != nil {
		// must parse at least one element
		return nil, 0, err
	}

	c := consumed
	// noMoreToConsume || next unconsumed byte is not ','
	if len(b[c:]) == 0 || b[c:][0] != ',' {
		return b[:c], c, nil
	}
	c++ // consume the ','

	// TODO what is the recusion limit in go? and will it limit
	// how many elements in a json 'list' we can support
	_, consumed, err = ParseElements(b[c:])
	if err != nil {
		c-- // unconsume the ','
		return b[:c], c, nil
	}

	c += consumed

	return b[:c], c, nil
}

type Element []byte

func ParseElement(b []byte) (Element, int, error) {
	// element
	//     ws value ws

	_, c := ParseWhitespace(b)

	_, consumed, err := ParseValue(b[c:])
	if err != nil {
		return nil, 0, err
	}
	c += consumed

	_, consumed = ParseWhitespace(b[c:])
	c += consumed

	return b[:c], c, nil
}

type String []byte

func ParseString(b []byte) (String, int, error) {
	// string
	//     '"' characters '"'

	if len(b) < 2 {
		// need at least two double quotes
		return nil, 0, ErrNothingToParse
	}

	if b[0] != '"' {
		return nil, 0, ErrInvalidStringOpen
	}

	c := 1 // we've consumed the first double quote

	_, consumed := ParseCharacters(b[c:])
	c += consumed

	// noMoreToConsume
	if len(b[c:]) == 0 {
		return nil, 0, ErrNothingToParse
	}

	// next unconsumed byte is not '"'
	if b[c:][0] != '"' {
		return nil, 0, ErrInvalidStringClose
	}

	c += 1 // consume final quote

	return b[:c], c, nil
}

type Number []byte

func (n Number) Int() (int64, error) {
	d := bytes.IndexByte(n, '.')
	if d == -1 {
		return strconv.ParseInt(string(n), 10, 64)
	}
	return 0, ErrParseInteger // TODO is this even being used?
}

func (n Number) Float64() float64 {
	ret, err := strconv.ParseFloat(string(n), 64)
	if err != nil {
		// then we did not parse digits correctly
		panic(err) // todo how do return message that makes sense, when we ParsedDigits there was an error that you ignored
	}
	return ret
}

func ParseNumber(b []byte) (Number, int, error) {
	// number
	//     int frac exp
	//
	// Note int is required, but frac and exp can be empty strings

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

type Int []byte

func (i Int) Int() int64 {
	ret, err := strconv.ParseInt(string(i), 10, 64)
	if err != nil {
		// then we did not parse digits correctly
		panic(err) // or better safe then sorry?... return 0
		// todo how do return message that makes sense, when you ParsedDigits you got an error that you ignored
	}
	return ret
}

func ParseInt(b []byte) (Int, int, error) {
	// int
	//     digit
	//     onenine digits
	//     '-' digit
	//     '-' onenine digits

	// note: int cannot have a '+' sign

	if len(b) == 0 {
		return nil, 0, ErrNothingToParse
	}

	if b[0] == '0' {
		return b[:1], 1, nil // digit
	}

	var c int
	if IsOneNine(b[0]) {
		_, c, err := ParseDigits(b)
		if err != nil {
			return b[:1], 1, nil // digit
		}
		//todo use c not consumed
		return b[:c], c, nil // onenine digits
	}

	if b[0] != '-' {
		return nil, 0, ErrUnexpectedChar
	}

	if len(b) > 1 && b[1] == '-' {
		return nil, 0, ErrUnexpectedChar
	}

	_, c, err := ParseInt(b[1:])
	if err != nil {
		return nil, 0, err
	}
	c += 1 // the negative

	return b[:c], c, nil
}

type Exp []byte

// returns 1 or -1 indicating the sign of the exponent
func (e Exp) Sign() int {
	sign, consumed := ParseSign(e)
	if consumed == 0 {
		return 1
	}
	if sign[0] == '-' {
		return -1
	}
	return 1
}

func ParseExp(b []byte) (Exp, int) {
	// exp
	//     ""
	//     'E' sign digits
	//     'e' sign digits

	if len(b) == 0 ||
		(b[0] != 'e' && b[0] != 'E') {
		return nil, 0 // empty string satisfies exp
	}

	c := 1 // we've [c]onsumed 'e'
	_, consumed := ParseSign(b[c:])
	c += consumed

	_, consumed, err := ParseDigits(b[c:])
	if err != nil {
		return nil, 0 // must have valid digits
	}
	c += consumed

	return b[0:c], c
}

// todo rename parse functions to consume and create interface{} Consumer... easier to test
type Sign []byte

func (s Sign) Positive() bool { return len(s) != 1 || s[0] == '+' }

func ParseSign(b []byte) (Sign, int) {
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

type Frac []byte

func (f Frac) Float64() float64 {
	// ParseFrac doesn't return error but can return b=empty string
	// so this should return 0 for the fraction part
	if len(f) == 0 {
		return float64(0)
	}

	ret, err := strconv.ParseFloat(string(f), 64)
	if err != nil {
		// then we did not parse float correctly
		panic(err)
	}

	return ret
}

func ParseFrac(b []byte) (Frac, int) {
	// frac
	//     ""
	//     '.' digits

	if len(b) == 0 || b[0] != '.' {
		return nil, 0
	}
	c := 1 // we've consumed the '.'

	_, consumed, err := ParseDigits(b[c:])
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

type Escape []byte

func ParseEscape(b []byte) (Escape, int, error) {
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
		return nil, 0, ErrNothingToParse
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

	return nil, 0, ErrInvalidEscape
}

type Whitespace []byte

func ParseWhitespace(b []byte) (Whitespace, int) {
	// ws
	//     ""
	//     '0009' ws
	//     '000a' ws
	//     '000d' ws
	//     '0020' ws

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

	return b[:n], n
}

// [wip] Digits satisfies the json.org language spec for both digit and digits

type Digits []byte

func (d Digits) Int() uint64 {
	// todo is there bytes.parseuint ?
	ret, err := strconv.ParseUint(string(d), 10, 64)
	if err != nil {
		// then we did not parse digits correctly
		panic(err) // or better safe then sorry?... return 0
		// todo how do return message that makes sense, when you ParsedDigits you got an error that you ignored
	}
	return ret
}

//  ParseDigits returns b[0:x] such that every ascii value from b[0]
//  to b[x] represents a digit from 0 to 9, along with the length of b[0:x]
func ParseDigits(b []byte) (Digits, int, error) {

	// digits
	//     digit
	//     digit digits

	if len(b) == 0 {
		return nil, 0, ErrNothingToParse
	}

	// check first digit
	if !IsDigit(b[0]) {
		return nil, 0, ErrInvalidDigit
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
