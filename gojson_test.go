package gojson

import (
	"testing"

	"github.com/jimmyjames85/gojson"
)

type testCase struct {
	name  string
	input []byte

	// expected
	expected []byte
	wantErr  bool
}

func TestParseSign(t *testing.T) {
	tests := []testCase{
		{
			name:     "no sign",
			input:    []byte("343"),
			expected: []byte(""),
		},
		{
			name:     "sign in middle",
			input:    []byte("  +001234567asdag"),
			expected: []byte(""),
		},
		{
			name:     "positive",
			input:    []byte("+343"),
			expected: []byte("+"),
		},
		{
			name:     "negative",
			input:    []byte("-343"),
			expected: []byte("-"),
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "emtpy slice",
			input:    []byte{},
			expected: []byte{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, actualLen := gojson.ParseSign(tc.input)

			if len(tc.expected) != actualLen {
				t.Errorf("unexpected length: wanted %d got %d", len(tc.expected), actualLen)
			}

			// byte.Compare
			if string(tc.expected) != string(actual) {
				t.Errorf("unexpected return: wanted %q got %q", string(tc.expected), string(actual))
			}
		})
	}
}

func TestParseWhitespace(t *testing.T) {
	tests := []testCase{
		{
			name:     "leading spaces",
			input:    []byte("      001234567asdag"),
			expected: []byte("      "),
		},
		{
			name:     "all whitespace",
			input:    []byte("      "),
			expected: []byte("      "),
		},
		{
			name:     "no whitespace",
			input:    []byte("001234567asdag"),
			expected: []byte(""),
		},
		{
			name:     "tab newline carriage return",
			input:    []byte("\t\r\n    \n\rasdfas"),
			expected: []byte("\t\r\n    \n\r"),
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "emtpy slice",
			input:    []byte{},
			expected: []byte{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			//t.Parallel()
			actual, actualLen := gojson.ParseWhitespace(tc.input)

			if len(tc.expected) != actualLen {
				t.Errorf("unexpected length: wanted %d got %d", len(tc.expected), actualLen)
			}

			// byte.Compare
			if string(tc.expected) != string(actual) {
				t.Errorf("unexpected return: wanted %q got %q", string(tc.expected), string(actual))
			}
		})
	}
}

func TestParseDigits(t *testing.T) {
	tests := []testCase{
		{
			name:     "leading zeros",
			input:    []byte("001234567asdag"),
			expected: []byte("001234567"),
		},
		{
			name:    "must start with digit",
			input:   []byte("asf234"),
			wantErr: true,
		},
		{
			name:     "consume everything",
			input:    []byte("123456"),
			expected: []byte("123456"),
		},
		{
			name:     "single digit",
			input:    []byte("1 "),
			expected: []byte("1"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			//t.Parallel()
			actual, actualLen, err := gojson.ParseDigits(tc.input)

			if tc.wantErr && err == nil {
				t.Errorf("expecting error but got <nil>")
			} else if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			if len(tc.expected) != actualLen {
				t.Errorf("unexpected length: wanted %d got %d", len(tc.expected), actualLen)
			}

			// byte.Compare
			if string(tc.expected) != string(actual) {
				t.Errorf("unexpected return: wanted %q got %q", string(tc.expected), string(actual))
			}
		})
	}
}

func TestParseExp(t *testing.T) {
	tests := []testCase{
		{
			name:     "consume nothing",
			input:    []byte("asdf"),
			expected: []byte(""),
		},
		{
			name:     "e with no sign",
			input:    []byte("e2432"),
			expected: []byte(""),
		},
		{
			name:     "e sign no digits",
			input:    []byte("e+abcd"),
			expected: []byte(""),
		},
		{
			name:     "happy path +",
			input:    []byte("e+2342foobar"),
			expected: []byte("e+2342"),
		},
		{
			name:     "captial E happy path +",
			input:    []byte("E+2342foobar"),
			expected: []byte("E+2342"),
		},
		{
			name:     "happy path -",
			input:    []byte("e-71 54"),
			expected: []byte("e-71"),
		},
		{
			name:     "captial E happy path -",
			input:    []byte("E-96 foobar"),
			expected: []byte("E-96"),
		},

		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "emtpy slice",
			input:    []byte{},
			expected: []byte{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			//t.Parallel()
			actual, actualLen := gojson.ParseExp(tc.input)

			if len(tc.expected) != actualLen {
				t.Errorf("unexpected length: wanted %d got %d", len(tc.expected), actualLen)
			}

			// byte.Compare
			if string(tc.expected) != string(actual) {
				t.Errorf("unexpected return: wanted %q got %q", string(tc.expected), string(actual))
			}
		})
	}
}

func TestParseFrac(t *testing.T) {
	tests := []testCase{
		{
			name:     "consume nothing",
			input:    []byte("asdf"),
			expected: []byte(""),
		},
		{
			name:     "just a . should consume nothing",
			input:    []byte(". asdf"),
			expected: []byte(""),
		},
		{
			name:     "begin with just digits should consume nothing",
			input:    []byte("234.324"),
			expected: []byte(""),
		},
		{
			name:     "happy path",
			input:    []byte(".14159 asdf"),
			expected: []byte(".14159"),
		},
		{
			name:     "multiple ....",
			input:    []byte("...."),
			expected: []byte(""),
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "emtpy slice",
			input:    []byte{},
			expected: []byte{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			actual, actualLen := gojson.ParseFrac(tc.input)

			if len(tc.expected) != actualLen {
				t.Errorf("unexpected length: wanted %d got %d", len(tc.expected), actualLen)
			}

			// byte.Compare
			if string(tc.expected) != string(actual) {
				t.Errorf("unexpected return: wanted %q got %q", string(tc.expected), string(actual))
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []testCase{
		{
			name:     "one digit",
			input:    []byte("0asdfas"),
			expected: []byte("0"),
		},
		{
			name:     "multiple digits",
			input:    []byte("2432323 asdjais"),
			expected: []byte("2432323"),
		},
		{
			name:     "can only consume one zero",
			input:    []byte("000002432323"),
			expected: []byte("0"),
		},
		{
			name:     "negative zero is ok",
			input:    []byte("-0 asdjais"),
			expected: []byte("-0"),
		},

		{
			name:     "negative multiple digits",
			input:    []byte("-123456 asdjais"),
			expected: []byte("-123456"),
		},
		{
			name:     "can only consume one negative zero",
			input:    []byte("-000002432323"),
			expected: []byte("-0"),
		},
		{
			name:    "cant have multiple negative signs",
			input:   []byte("----123456 asdjais"),
			wantErr: true,
		},
		{
			name:     "negative single digit",
			input:    []byte("-5 "),
			expected: []byte("-5"),
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "emtpy slice",
			input:   []byte{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			actual, actualLen, err := gojson.ParseInt(tc.input)

			if tc.wantErr && err == nil {
				t.Errorf("expecting error but got <nil>")
			} else if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			if len(tc.expected) != actualLen {
				t.Errorf("unexpected length: wanted %d got %d", len(tc.expected), actualLen)
			}

			// byte.Compare
			if string(tc.expected) != string(actual) {
				t.Errorf("unexpected return: wanted %q got %q", string(tc.expected), string(actual))
			}
		})
	}
}

func TestParseNumber(t *testing.T) {
	tests := []testCase{
		{
			name:     "one digit",
			input:    []byte("0asdfas"),
			expected: []byte("0"),
		},
		{
			name:     "multiple digits",
			input:    []byte("2432323 asdjais"),
			expected: []byte("2432323"),
		},
		{
			name:     "can only consume one zero",
			input:    []byte("000002432323"),
			expected: []byte("0"),
		},
		{
			name:     "negative zero is ok",
			input:    []byte("-0 asdjais"),
			expected: []byte("-0"),
		},

		{
			name:     "negative multiple digits",
			input:    []byte("-123456 asdjais"),
			expected: []byte("-123456"),
		},
		{
			name:     "can only consume one negative zero",
			input:    []byte("-000002432323"),
			expected: []byte("-0"),
		},
		{
			name:    "cant have multiple negative signs",
			input:   []byte("----123456 asdjais"),
			wantErr: true,
		},
		{
			name:     "negative single digit",
			input:    []byte("-5 "),
			expected: []byte("-5"),
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "emtpy slice",
			input:   []byte{},
			wantErr: true,
		},
		{
			name:     "decimal bad exp",
			input:    []byte("-0.5e232 jim"),
			expected: []byte("-0.5"),
		},
		{
			name:     "decimal with exp +",
			input:    []byte("-0.5e+232 jim"),
			expected: []byte("-0.5e+232"),
		},
		{
			name:     "decimal with EXP -",
			input:    []byte("-20.500E-232 jim"),
			expected: []byte("-20.500E-232"),
		},
		{
			name:     "dont consume bad exp",
			input:    []byte("3.14159E -6 jim"),
			expected: []byte("3.14159"),
		},
		{
			name:     "exp with no frac",
			input:    []byte("-5e+232 foo"),
			expected: []byte("-5e+232"),
		},
		{
			name:     "frac with no exp",
			input:    []byte("-5.03 bar"),
			expected: []byte("-5.03"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			actual, actualLen, err := gojson.ParseNumber(tc.input)

			if tc.wantErr && err == nil {
				t.Errorf("expecting error but got <nil>")
			} else if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			if len(tc.expected) != actualLen {
				t.Errorf("unexpected length: wanted %d got %d", len(tc.expected), actualLen)
			}

			// byte.Compare
			if string(tc.expected) != string(actual) {
				t.Errorf("unexpected return: wanted %q got %q", string(tc.expected), string(actual))
			}
		})
	}
}

func TestParseEscape(t *testing.T) {
	tests := []testCase{
		{
			name:     "double quote",
			input:    []byte(`"asdfad`),
			expected: []byte(`"`),
		},
		{
			name:     "backslash",
			input:    []byte(`\ forever`),
			expected: []byte(`\`),
		},
		{
			name:     "forwardslash",
			input:    []byte(`/ never`),
			expected: []byte(`/`),
		},

		{
			name:     "the letter b backspace?",
			input:    []byte(`b maybe`),
			expected: []byte(`b`),
		},
		{
			name:     "n for newline",
			input:    []byte(`nthis would be on a newline`),
			expected: []byte(`n`),
		},

		{
			name:     "r for caraige return",
			input:    []byte(`r whatabout linefeed? is that just \n`),
			expected: []byte(`r`),
		},
		{
			name:     "t is for tab",
			input:    []byte(`talso is whitespace`),
			expected: []byte(`t`),
		},
		{
			name:    "anything else should error",
			input:   []byte("i am not an escape"),
			wantErr: true,
		},
		{
			name:     "unicode escape",
			input:    []byte("uf0A02"),
			expected: []byte("uf0A0"),
		},
		{
			name:    "bad unicode escape",
			input:   []byte("uG0a0 is bad"),
			wantErr: true,
		},
		{
			name:    "not enough for unicode escape",
			input:   []byte("uFA0 missing just one more"),
			wantErr: true,
		},
		{
			name:    "no more after unicode escape",
			input:   []byte("u where is the rest"),
			wantErr: true,
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "emtpy slice",
			input:   []byte{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			actual, actualLen, err := gojson.ParseEscape(tc.input)

			if tc.wantErr && err == nil {
				t.Errorf("expecting error but got <nil>")
			} else if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			if len(tc.expected) != actualLen {
				t.Errorf("unexpected length: wanted %d got %d", len(tc.expected), actualLen)
			}

			// byte.Compare
			if string(tc.expected) != string(actual) {
				t.Errorf("unexpected return: wanted %q got %q", string(tc.expected), string(actual))
			}
		})
	}
}

func TestParseCharacter(t *testing.T) {
	tests := []testCase{
		{
			name:    "less than 0x0020",
			input:   []byte{0x19},
			wantErr: true,
		},
		{
			name:     "0x10ffff is okay",
			input:    []byte(string([]rune{0x10ffff})),
			expected: []byte(string([]rune{0x10ffff})),
		},
		{
			name:    "greater than 0x10ffff",
			input:   []byte(string([]rune{0x110000})),
			wantErr: true,
		},
		{
			name:     "some utf8 char",
			input:    []byte("世 wat is this"),
			expected: []byte("世"),
		},
		{
			name:     "escape sequence",
			input:    []byte(`\r asdfa`),
			expected: []byte(`\r`),
		},
		{
			name:     "escape sequence forward slash",
			input:    []byte(`\/ asdfa`),
			expected: []byte(`\/`),
		},
		{
			name:    "invalid escape sequence",
			input:   []byte(`\5 asdfa`),
			wantErr: true,
		},
		{
			name:    "invalid escape sequence stops",
			input:   []byte(`\`),
			wantErr: true,
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "emtpy slice",
			input:   []byte{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			actual, actualLen, err := gojson.ParseCharacter(tc.input)

			if tc.wantErr && err == nil {
				t.Errorf("expecting error but got <nil>")
			} else if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			if len(tc.expected) != actualLen {
				t.Errorf("unexpected length: wanted %d got %d", len(tc.expected), actualLen)
			}

			// byte.Compare
			if string(tc.expected) != string(actual) {
				t.Errorf("unexpected return: wanted %q got %q", string(tc.expected), string(actual))
			}
		})
	}
}

func TestParseCharacters(t *testing.T) {
	tests := []testCase{
		{
			name:     "less than 0x0020",
			input:    []byte{0x19},
			expected: []byte{},
		},
		{
			name:     "0x10ffff is okay but shouldnt get 0x19",
			input:    []byte(string([]rune{0x60, 0x10ffff, 0x60, 0x19})),
			expected: []byte(string([]rune{0x60, 0x10ffff, 0x60})),
		},
		{
			name:     "stop at anything greater than 0x10ffff",
			input:    []byte(string([]rune{0x40, 0x34, 0x110000})),
			expected: []byte(string([]rune{0x40, 0x34})),
		},
		{
			name:     "some utf8 chars",
			input:    []byte("Cos'è ユニコードとはか？ ಯುನಿಕೋಡ್ ಎಂದರೇನು 유니코드에 대해 Што е युनिकोड के हो? гэж юу вэ qu'es aquò  يونی‌کُد چيست؟ யூனிக்கோடு என்றால் என்ன యూనీకోడ్ అంటే ఏమిటి"),
			expected: []byte("Cos'è ユニコードとはか？ ಯುನಿಕೋಡ್ ಎಂದರೇನು 유니코드에 대해 Што е युनिकोड के हो? гэж юу вэ qu'es aquò  يونی‌کُد چيست؟ யூனிக்கோடு என்றால் என்ன యూనీకోడ్ అంటే ఏమిటి"),
		},
		{
			name:     "strings with escape sequence",
			input:    []byte(`this is a line and \r\nthis should be on a newline \t with a tab`),
			expected: []byte(`this is a line and \r\nthis should be on a newline \t with a tab`),
		},
		{
			name:     "invalid escape sequence",
			input:    []byte(`this should stop here\5 asdfa`),
			expected: []byte(`this should stop here`),
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "emtpy slice",
			input:    []byte{},
			expected: []byte{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			actual, actualLen := gojson.ParseCharacters(tc.input)

			if len(tc.expected) != actualLen {
				t.Errorf("unexpected length: wanted %d got %d", len(tc.expected), actualLen)
			}

			// byte.Compare
			if string(tc.expected) != string(actual) {
				t.Errorf("unexpected return: wanted %q got %q", string(tc.expected), string(actual))
			}
		})
	}
}

func TestParseString(t *testing.T) {
	tests := []testCase{
		{
			name:    "no quotes",
			input:   []byte(`a asksja asdf`),
			wantErr: true,
		},
		{
			name:    "only one quote",
			input:   []byte(`"oh we started of soo good`),
			wantErr: true,
		},
		{
			name:    "two quotes but invalid escape sequence",
			input:   []byte(`"oh we started of \u00Ga but messed up in the middle"`),
			wantErr: true,
		},
		{
			name:     "happy path with quotes",
			input:    []byte(`"the cat in the hat" is really fat`),
			expected: []byte(`"the cat in the hat"`),
		},
		{
			name:     "some utf8 char",
			input:    []byte(`"世 wat is this"  世`),
			expected: []byte(`"世 wat is this"`),
		},

		{
			name:     "some more utf8 chars",
			input:    []byte(`"Cos'è ユニコードとはか？ ಯುನಿಕೋಡ್" what does it say ಎಂದರೇನು 유니코드에 대해 Што е युनिकोड के हो? гэж юу вэ qu'es aquò  يونی‌کُد چيست؟ யூனிக்கோடு என்றால் என்ன యూనీకోడ్ అంటే ఏమిటి`),
			expected: []byte(`"Cos'è ユニコードとはか？ ಯುನಿಕೋಡ್"`),
		},
		{
			name:    "out of range char is bad",
			input:   []byte{0x22, 0x19, 0x22},
			wantErr: true,
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "emtpy slice",
			input:   []byte{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			actual, actualLen, err := gojson.ParseString(tc.input)

			if tc.wantErr && err == nil {
				t.Errorf("expecting error but got <nil>")
			} else if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			if len(tc.expected) != actualLen {
				t.Errorf("unexpected length: wanted %d got %d", len(tc.expected), actualLen)
			}

			// byte.Compare
			if string(tc.expected) != string(actual) {
				t.Errorf("unexpected return: wanted %q got %q", string(tc.expected), string(actual))
			}
		})
	}
}

func TestParseElement(t *testing.T) {
	tests := []testCase{
		{
			name:     "null",
			input:    []byte("  null \n \r \t   ,STOP sadfasf "),
			expected: []byte("  null \n \r \t   "),
		},
		{
			name:     "true",
			input:    []byte("\n\n\n\r\t  \t true \n NOMORE"),
			expected: []byte("\n\n\n\r\t  \t true \n "),
		},
		{
			name:     "false",
			input:    []byte(" false  \t , what: "),
			expected: []byte(" false  \t "),
		},
		{
			name:    "TRUE cant be in caps",
			input:   []byte("  \t TRUE: \n what: "),
			wantErr: true,
		},
		{
			name:    "at least must have an element: all white space is bad",
			input:   []byte("     \t \n"),
			wantErr: true,
		},
		{
			name:     "a string",
			input:    []byte(" \t \r \n \"string goes here with an escaped \\n newline : ) \" \r \n \t    <- more whitespace should be consumed as well"),
			expected: []byte(" \t \r \n \"string goes here with an escaped \\n newline : ) \" \r \n \t    "),
		},
		{
			name:    "a bad string: cannot have an unescapped newline in a string",
			input:   []byte(" \t \r \n \"string goes here \n\""),
			wantErr: true,
		},
		// todo not all values are tested e.g. Number
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			actual, actualLen, err := gojson.ParseElement(tc.input)

			if tc.wantErr && err == nil {
				t.Errorf("expecting error but got <nil>")
			} else if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			if len(tc.expected) != actualLen {
				t.Errorf("unexpected length: wanted %d got %d", len(tc.expected), actualLen)
			}

			// byte.Compare
			if string(tc.expected) != string(actual) {
				t.Errorf("unexpected return: wanted %q got %q", string(tc.expected), string(actual))
			}
		})
	}
}

func TestParseElements(t *testing.T) {
	tests := []testCase{
		{
			name:     "signle null",
			input:    []byte("  null \n \r \t   ,STOP sadfasf "),
			expected: []byte("  null \n \r \t   "),
		},
		{
			name:     "two elements",
			input:    []byte("\n\n\n\r\t  \t true \n, false NOMORE"),
			expected: []byte("\n\n\n\r\t  \t true \n, false "),
		},
		{
			name:     "dont consume last comma , ",
			input:    []byte(" 234, 2342 , 222 , STOP"),
			expected: []byte(" 234, 2342 , 222 "),
		},
		{
			name:    "can't start with a comma",
			input:   []byte(", 34,  \t TRUE: \n what: "),
			wantErr: true,
		},
		{
			name:    "at least must have an element: all white space is bad",
			input:   []byte("     \t \n"),
			wantErr: true,
		},
		{
			name:     "consume as many proper elements as possible: not the bad string",
			input:    []byte("12, 1233 , 12424,  \t \r \n \"string goes here \n\", 234, 22 "),
			expected: []byte("12, 1233 , 12424"),
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "emtpy slice",
			input:   []byte{},
			wantErr: true,
		},
		// todo not all values are tested e.g. Number
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			actual, actualLen, err := gojson.ParseElements(tc.input)

			if tc.wantErr && err == nil {
				t.Errorf("expecting error but got <nil>")
			} else if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			if len(tc.expected) != actualLen {
				t.Errorf("unexpected length: wanted %d got %d", len(tc.expected), actualLen)
			}

			// byte.Compare
			if string(tc.expected) != string(actual) {
				t.Errorf("unexpected return: wanted %q got %q", string(tc.expected), string(actual))
			}
		})
	}
}

func TestParseArray(t *testing.T) {
	tests := []testCase{
		{
			name:     "white space array",
			input:    []byte("[   \n \n \r \t ]  STOP here"),
			expected: []byte("[   \n \n \r \t ]"),
		},
		{
			name:     "empty array",
			input:    []byte("[]  STOP here"),
			expected: []byte("[]"),
		},
		{
			name:    "must start with [",
			input:   []byte("   []"),
			wantErr: true,
		},
		{
			name:     "non-empty array",
			input:    []byte("[ 234, 2342 , 222 , true, null ] no mo"),
			expected: []byte("[ 234, 2342 , 222 , true, null ]"),
		},
		{
			name:    "can't start with a comma",
			input:   []byte("[ , 234, 2342 , 222 , true, null ] no mo"),
			wantErr: true,
		},
		{
			name:    "invalid element",
			input:   []byte("[ 234, 2342 , 222 , TRUE, null ]"),
			wantErr: true,
		},
		{
			name:    "missing ] is bad",
			input:   []byte(`[ "jim", "was", true, "but forgot the ]"`),
			wantErr: true,
		},
		{
			name:     "parse ] outside of string",
			input:    []byte(`[ "jim", "was", true, "but forgot the ]"]`),
			expected: []byte(`[ "jim", "was", true, "but forgot the ]"]`),
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "emtpy slice",
			input:   []byte{},
			wantErr: true,
		},

		// todo not all values are tested e.g. Object
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			actual, actualLen, err := gojson.ParseArray(tc.input)

			if tc.wantErr && err == nil {
				t.Errorf("expecting error but got <nil>")
			} else if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			if len(tc.expected) != actualLen {
				t.Errorf("unexpected length: wanted %d got %d", len(tc.expected), actualLen)
			}

			// byte.Compare
			if string(tc.expected) != string(actual) {
				t.Errorf("unexpected return: wanted %q got %q", string(tc.expected), string(actual))
			}
		})
	}
}

func TestParseMember(t *testing.T) {
	tests := []testCase{
		{
			name:     "simple member",
			input:    []byte(`   "key": 2343    stop here`),
			expected: []byte(`   "key": 2343    `),
		},
		{
			name:     "whitespace is okay",
			input:    []byte("  \r \n \"thisisthekey\"  \t\n\r : \r\n null\r\t  STOP here"),
			expected: []byte("  \r \n \"thisisthekey\"  \t\n\r : \r\n null\r\t  "),
		},
		{
			name:    "must start with whitespace and/or a string",
			input:   []byte(`   null, "not":ok`),
			wantErr: true,
		},
		{
			name:     "value is an array is okay",
			input:    []byte(`"my_array": [ 234, 2342 , 222 , true, null ] no mo`),
			expected: []byte(`"my_array": [ 234, 2342 , 222 , true, null ] `),
		},
		{
			name:    "invalid string",
			input:   []byte(" \"we cant have \n\": 34"),
			wantErr: true,
		},
		{
			name:     "valid string",
			input:    []byte(" \"we can have ws\"   :   \"\\t\\n34\"  stop"),
			expected: []byte(" \"we can have ws\"   :   \"\\t\\n34\"  "),
		},
		{
			name:     "emtpy strings are okay",
			input:    []byte(` ""  :  ""  stop`),
			expected: []byte(` ""  :  ""  `),
		},
		{
			name:     "should consume the entire input",
			input:    []byte(` ""  :  ""   `),
			expected: []byte(` ""  :  ""   `),
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "emtpy slice",
			input:   []byte{},
			wantErr: true,
		},

		// todo not all values are tested e.g. Object
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			actual, actualLen, err := gojson.ParseMember(tc.input)

			if tc.wantErr && err == nil {
				t.Errorf("expecting error but got <nil>")
			} else if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			if len(tc.expected) != actualLen {
				t.Errorf("unexpected length: wanted %d got %d", len(tc.expected), actualLen)
			}

			// byte.Compare
			if string(tc.expected) != string(actual) {
				t.Errorf("unexpected return: wanted %q got %q", string(tc.expected), string(actual))
			}
		})
	}
}

func TestParseMembers(t *testing.T) {
	tests := []testCase{
		{
			name:     "simple member",
			input:    []byte(`   "key": 2343    stop here`),
			expected: []byte(`   "key": 2343    `),
		},
		{
			name:     "two members",
			input:    []byte(`"jim":null, "key":[3 ,2 ,2,343,"jim"] STOP`),
			expected: []byte(`"jim":null, "key":[3 ,2 ,2,343,"jim"] `),
		},
		{
			name:     "consume all",
			input:    []byte(`"jim":null, "key": "foo", "wow":   2234334      `),
			expected: []byte(`"jim":null, "key": "foo", "wow":   2234334      `),
		},
		{
			name:    "must start with whitespace and/or a string",
			input:   []byte(`   ,null, "not":ok`),
			wantErr: true,
		},
		{
			name:     "consume what we can: invalid string",
			input:    []byte("\"firstkey\":2, \"we cant have \n\": 34"),
			expected: []byte("\"firstkey\":2"),
		},
		{
			name:     "emtpy strings are okay",
			input:    []byte(` ""  :  ""  , "key": "stop"   `),
			expected: []byte(` ""  :  ""  , "key": "stop"   `),
		},
		{
			name:     "whitespace between memebres",
			input:    []byte("\"key1\" \n  :  \r\t\n 1 \n  , \n \"key2\":2   "),
			expected: []byte("\"key1\" \n  :  \r\t\n 1 \n  , \n \"key2\":2   "),
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "emtpy slice",
			input:   []byte{},
			wantErr: true,
		},
		// todo not all values are tested e.g. Object
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			actual, actualLen, err := gojson.ParseMembers(tc.input)

			if tc.wantErr && err == nil {
				t.Errorf("expecting error but got <nil>")
			} else if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			if len(tc.expected) != actualLen {
				t.Errorf("unexpected length: wanted %d got %d", len(tc.expected), actualLen)
			}

			// byte.Compare
			if string(tc.expected) != string(actual) {
				t.Errorf("unexpected return: wanted %q got %q", string(tc.expected), string(actual))
			}
		})
	}
}

func TestParseObject(t *testing.T) {
	tests := []testCase{
		{
			name:     "empty object",
			input:    []byte(`{}   `),
			expected: []byte(`{}`),
		},
		{
			name:     "whitespace empty object",
			input:    []byte("{   \r \t \n     }   "),
			expected: []byte("{   \r \t \n     }"),
		},
		{
			name:     "one member",
			input:    []byte(`{ "foo"   : 3422 }  STOP at curly boi`),
			expected: []byte(`{ "foo"   : 3422 }`),
		},
		{
			name:     "two members",
			input:    []byte(`{"jim":null, "key":[3 ,2 ,2,343,"jim"] } STOP`),
			expected: []byte(`{"jim":null, "key":[3 ,2 ,2,343,"jim"] }`),
		},
		{
			name:    "invalid member",
			input:   []byte(`{"jim":null, "key":[3 ,2 ,2,343,"jim"], INVALID }`),
			wantErr: true,
		},
		{
			name:     "consume all",
			input:    []byte(`{"jim":null, "key": "foo", "wow":   2234334      }`),
			expected: []byte(`{"jim":null, "key": "foo", "wow":   2234334      }`),
		},
		{
			name:    "must start with curly brace",
			input:   []byte(`   ,null, "not":ok`),
			wantErr: true,
		},
		{
			name:    "invalid member string should error",
			input:   []byte("{\"firstkey\":2, \"we cant have \n\": 34   }"),
			wantErr: true,
		},
		{
			name:     "emtpy strings are okay",
			input:    []byte(`{ ""  :  ""  , "key": "stop"   }   stop parsing this`),
			expected: []byte(`{ ""  :  ""  , "key": "stop"   }`),
		},
		{
			name:     "whitespace between memebres",
			input:    []byte("{\"key1\" \n  :  \r\t\n 1 \n  , \n \"key2\":2   }"),
			expected: []byte("{\"key1\" \n  :  \r\t\n 1 \n  , \n \"key2\":2   }"),
		},
		{
			name:     "member is an object",
			input:    []byte(`{"key" : [ {"sub":1}, 23, null], "obj": {"woah":23}}`),
			expected: []byte(`{"key" : [ {"sub":1}, 23, null], "obj": {"woah":23}}`),
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "emtpy slice",
			input:   []byte{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			actual, actualLen, err := gojson.ParseObject(tc.input)

			if tc.wantErr && err == nil {
				t.Errorf("expecting error but got <nil>")
			} else if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}

			if len(tc.expected) != actualLen {
				t.Errorf("unexpected length: wanted %d got %d", len(tc.expected), actualLen)
			}

			// byte.Compare
			if string(tc.expected) != string(actual) {
				t.Errorf("unexpected return: wanted %q got %q", string(tc.expected), string(actual))
			}
		})
	}
}
