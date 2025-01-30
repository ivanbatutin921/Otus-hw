package hw09structvalidator

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseStringRule(t *testing.T) {
	tests := []struct {
		name          string
		rule          string
		wantErr       bool
		returnedError error
	}{
		{"len ok", "len:32", false, nil},
		{"regexp ok", "regexp:\\d", false, nil},
		{"in ok", "in:a", false, nil},
		{"without ':'", "len32", true, ErrInvalidTagSyntax},
		{"invalid rule", "invalid:rule", true, UnsupportedTagRuleError{"invalid"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStringRule(tt.rule)

			if tt.wantErr {
				require.ErrorIs(t, err, tt.returnedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
			}
		})
	}
}

func TestStringLenRule(t *testing.T) {
	tests := []struct {
		name           string
		controlValue   string
		checkValueOk   string
		checkValueFail string
		wantErr        bool
		returnedError  error
	}{
		{"ok", "2", "aa", "aaa", false, nil},
		{"invalid arg", "short", "", "", true, strconv.ErrSyntax},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := stringRule{}
			got, err := sr.getLenRule(tt.controlValue)

			if tt.wantErr {
				require.ErrorIs(t, err, tt.returnedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NoError(t, got(tt.checkValueOk))
				require.ErrorIs(t, got(tt.checkValueFail), ErrStringLen)
			}
		})
	}
}

func TestStringRegexpRule(t *testing.T) {
	tests := []struct {
		name           string
		controlValue   string
		checkValueOk   string
		checkValueFail string
		wantErr        bool
		returnedError  string
	}{
		{"ok", "\\d+", "1234", "aaa", false, ""},
		{"invalid arg", "?d", "", "", true, "error parsing regexp: missing argument to repetition operator: `?`"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := stringRule{}
			got, err := sr.getRegexpRule(tt.controlValue)

			if tt.wantErr {
				require.EqualError(t, err, tt.returnedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NoError(t, got(tt.checkValueOk))
				require.ErrorIs(t, got(tt.checkValueFail), ErrStringRegexp)
			}
		})
	}
}

func TestStringInRule(t *testing.T) {
	tests := []struct {
		name           string
		controlValue   string
		checkValueOk   string
		checkValueFail string
		wantErr        bool
		returnedError  error
	}{
		{"ok", "a,b", "a", "c", false, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := stringRule{}
			got, err := sr.getInRule(tt.controlValue)

			if tt.wantErr {
				require.ErrorIs(t, err, tt.returnedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NoError(t, got(tt.checkValueOk))
				require.ErrorIs(t, got(tt.checkValueFail), ErrStringIn)
			}
		})
	}
}
