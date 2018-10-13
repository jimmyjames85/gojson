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
			name:    "invalid escape sequence",
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
		// {
		// 	name:    "greater than 0x10ffff",
		// 	input:   []byte(string([]rune{0x110000})),
		// },
		// {
		// 	name:     "some utf8 char",
		// 	input:    []byte("世 wat is this"),
		// 	expected: []byte("世"),
		// },
		// {
		// 	name:     "escape sequence",
		// 	input:    []byte(`\r asdfa`),
		// 	expected: []byte(`\r`),
		// },
		// {
		// 	name:     "escape sequence forward slash",
		// 	input:    []byte(`\/ asdfa`),
		// 	expected: []byte(`\/`),
		// },
		// {
		// 	name:    "invalid escape sequence",
		// 	input:   []byte(`\5 asdfa`),
		// 	wantErr: true,
		// },
		// {
		// 	name:    "invalid escape sequence",
		// 	input:   []byte(`\`),
		// 	wantErr: true,
		// },
		// {
		// 	name:    "nil input",
		// 	input:   nil,
		// 	wantErr: true,
		// },
		// {
		// 	name:    "emtpy slice",
		// 	input:   []byte{},
		// 	wantErr: true,
		// },
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
