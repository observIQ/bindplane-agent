package logcountprocessor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpressionMatch(t *testing.T) {
	var testCases = []struct {
		name        string
		expr        string
		env         map[string]interface{}
		expected    bool
		expectedErr error
	}{
		{
			name:     "simple true",
			expr:     "true",
			env:      map[string]interface{}{},
			expected: true,
		},
		{
			name:     "simple false",
			expr:     "false",
			env:      map[string]interface{}{},
			expected: false,
		},
		{
			name:     "true with env",
			expr:     `foo == "bar"`,
			env:      map[string]interface{}{"foo": "bar"},
			expected: true,
		},
		{
			name:     "false with env",
			expr:     `foo == "bar"`,
			env:      map[string]interface{}{"foo": "baz"},
			expected: false,
		},
		{
			name:        "invalid expression",
			expr:        `foo`,
			env:         map[string]interface{}{"foo": "bar"},
			expectedErr: errors.New("expression did not return a boolean"),
		},
		{
			name:        "invalid env",
			expr:        `foo + "bar"`,
			env:         map[string]interface{}{"foo": 1},
			expectedErr: errors.New("invalid operation: int + string"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := NewExpression(tc.expr)
			require.NoError(t, err)

			matches, err := expr.Match(tc.env)
			if tc.expectedErr != nil {
				require.Contains(t, err.Error(), tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, matches)
		})
	}
}

func TestNewExpression(t *testing.T) {
	var testCases = []struct {
		name        string
		expr        string
		expectedErr error
	}{
		{
			name: "simple true",
			expr: "true",
		},
		{
			name: "simple false",
			expr: "false",
		},
		{
			name:        "invalid expression",
			expr:        "",
			expectedErr: errors.New("unexpected token EOF"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewExpression(tc.expr)
			if tc.expectedErr != nil {
				require.Contains(t, err.Error(), tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
