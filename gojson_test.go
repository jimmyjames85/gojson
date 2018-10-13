package gojson

import "testing"

type testCase struct {
	name  string
	input []byte

	// expected
	expected    []byte
	expectedLen int
	wantErr     bool
}

func TestParseDigits(t *testing.T) {
	tests := []testCase{
		{
			name:        "leading zeros",
			input:       []byte("001234567asdag"),
			expected:    []byte("001234567"),
			expectedLen: 9,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual, actualLen, err := ParseDigits(tc.input)

			if tc.wantErr && err == nil {
				t.Errorf("test: %s: expecting error but got <nil>", tc.name)
			} else if !tc.wantErr && err != nil {
				t.Errorf("test: %s: unexpected error: %s", tc.name, err.Error())
			}

			if tc.expectedLen != actualLen {
				t.Errorf("test: %s: unexpected length: wanted %d got %d", tc.name, tc.expectedLen, actualLen)
			}

			if string(tc.expected) != string(actual) {
				t.Errorf("test: %s: unexpected return: wanted %s got %s", tc.name, string(tc.expected), string(actual))
			}
		})
	}
}
