package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: &User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "John",
				Age:    30,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil, // Данные валидны
		},
		{
			in: &User{
				ID:     "short-id",
				Name:   "John",
				Age:    17, // Меньше минимального возраста
				Email:  "invalid-email",
				Role:   "admin",           // Неверный Role
				Phones: []string{"12345"}, // Длина не равна 11
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("field %q must be exactly 36 characters long", "ID")},
				{Field: "Age", Err: fmt.Errorf("field %q must be at least 18", "Age")},
				{Field: "Email", Err: fmt.Errorf("field %q does not match the pattern", "Email")},
				{Field: "Role", Err: fmt.Errorf("field %q must be one of [admin,stuff]", "Role")},
				{Field: "Phones", Err: fmt.Errorf("field %q must be exactly 11 characters long", "Phones")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			// вызов Validate
			err := Validate(tt.in)

			// проверяем, если ошибок не ожидалось
			if tt.expectedErr == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
				return
			}

			// Проверяем, что ошибка является ValidationErrors
			if tt.expectedErr != nil {
				validationErrors, ok := err.(ValidationErrors)
				if !ok {
					t.Errorf("expected ValidationErrors, got %T", err)
					return
				}

				// проверяем количество ошибок
				expectedErrors := tt.expectedErr.(ValidationErrors)
				if len(validationErrors) != len(expectedErrors) {
					t.Errorf("unexpected number of validation errors: got %d, want %d", len(validationErrors), len(expectedErrors))
					return
				}

				// проверяем каждую ошибку
				for j, ve := range validationErrors {
					expected := expectedErrors[j]
					if ve.Field != expected.Field {
						t.Errorf("unexpected validation error: got %v, want %v", ve, expected)
					}
					if ve.Err.Error() != expected.Err.Error() {
						t.Errorf("unexpected error message: got %v, want %v", ve.Err.Error(), expected.Err.Error())
					}
				}
			}
		})
	}
}

