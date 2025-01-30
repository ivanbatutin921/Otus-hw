package hw09structvalidator

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Request struct {
		ID   int `validate:"min:5"`
		Apps []App
	}
)

func TestValidateUnsupportedType(t *testing.T) {
	expectedErr := ErrUnsupportedInputType
	err := Validate(42)
	require.Equal(t, expectedErr, err, "errors should be equal")
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			"validate user ok",
			User{"012345678901234567890123456789012345", "user", 25, "example@otus.ru", "admin", []string{"01234567890"}, []byte("json")}, //nolint
			ValidationErrors{},
		},
		{
			"validate user fail",
			User{"id", "user", 10, "example", "user", []string{"01234567890", "2"}, []byte("json")},
			ValidationErrors{
				{"ID", ErrStringLen},
				{"Age", ErrIntMin},
				{"Email", ErrStringRegexp},
				{"Role", ErrStringIn},
				{"Phones", ErrStringLen},
			},
		},
		{
			"validate app ok",
			App{"01234"},
			ValidationErrors{},
		},
		{
			"validate app fail",
			App{"0"},
			ValidationErrors{
				{"Version", ErrStringLen},
			},
		},
		{
			"validate token ok",
			Token{[]byte("header"), []byte("payload"), []byte("signature")},
			ValidationErrors{},
		},
		{
			"validate response ok",
			Response{200, "322"},
			ValidationErrors{},
		},
		{
			"validate response fail",
			Response{155, "322"},
			ValidationErrors{
				{"Code", ErrIntIn},
			},
		},
		{
			"validate request ok",
			Request{12345, []App{{"0.0.1"}}},
			ValidationErrors{},
		},
		{
			"validate request fail",
			Request{12345, []App{{"0.0.1"}, {"0.1"}}},
			ValidationErrors{
				{"Apps", ValidationErrors{
					{"Version", ErrStringLen},
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.ElementsMatch(t, tt.expectedErr, err, "errors should be equal")
		})
	}
}
