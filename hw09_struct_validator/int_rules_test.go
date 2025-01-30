package hw09structvalidator

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseIntRule(t *testing.T) {
	tests := []struct {
		name          string
		rule          string
		wantErr       bool
		returnedError error
	}{
		{"min ok", "min:25", false, nil},
		{"max ok", "max:36", false, nil},
		{"in ok", "in:15", false, nil},
		{"without ':'", "min,25", true, ErrInvalidTagSyntax},
		{"invalid rule", "invalid:rule", true, UnsupportedTagRuleError{"invalid"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIntRule(tt.rule)

			if tt.wantErr {
				require.ErrorIs(t, err, tt.returnedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
			}
		})
	}
}

func TestIntMinRule(t *testing.T) {
	tests := []struct {
		name           string
		controlValue   string
		checkValueOk   int64
		checkValueFail int64
		wantErr        bool
		returnedError  error
	}{
		{"ok", "2", 2, 1, false, nil},
		{"invalid arg", "short", 0, 0, true, strconv.ErrSyntax},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir := intRule{}
			got, err := ir.getMinRule(tt.controlValue)

			if tt.wantErr {
				require.ErrorIs(t, err, tt.returnedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NoError(t, got(tt.checkValueOk))
				require.ErrorIs(t, got(tt.checkValueFail), ErrIntMin)
			}
		})
	}
}

func TestIntMaxRule(t *testing.T) {
	tests := []struct {
		name           string
		controlValue   string
		checkValueOk   int64
		checkValueFail int64
		wantErr        bool
		returnedError  error
	}{
		{"ok", "2", 2, 3, false, nil},
		{"invalid arg", "short", 0, 0, true, strconv.ErrSyntax},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir := intRule{}
			got, err := ir.getMaxRule(tt.controlValue)

			if tt.wantErr {
				require.ErrorIs(t, err, tt.returnedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NoError(t, got(tt.checkValueOk))
				require.ErrorIs(t, got(tt.checkValueFail), ErrIntMax)
			}
		})
	}
}

func TestIntInRule(t *testing.T) {
	tests := []struct {
		name           string
		controlValue   string
		checkValueOk   int64
		checkValueFail int64
		wantErr        bool
		returnedError  error
	}{
		{"ok", "10,12", 10, 11, false, nil},
		{"empty arg", "", 0, 0, true, strconv.ErrSyntax},
		{"invalid arg", ",10", 0, 0, true, strconv.ErrSyntax},
		{"invalid arg", ",a", 0, 0, true, strconv.ErrSyntax},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir := intRule{}
			got, err := ir.getInRule(tt.controlValue)

			if tt.wantErr {
				require.ErrorIs(t, err, tt.returnedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				require.NoError(t, got(tt.checkValueOk))
				require.ErrorIs(t, got(tt.checkValueFail), ErrIntIn)
			}
		})
	}
}
