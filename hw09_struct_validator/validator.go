package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	sb := new(strings.Builder)
	for _, e := range v {
		sb.WriteString(e.Field + ": " + e.Err.Error() + "\n")
	}
	return sb.String()
}

func Validate(v interface{}) error {
	// Шаг 1: Получаем тип и значение структуры
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return errors.New("not a struct")
	}

	var validationErrors ValidationErrors

	// Шаг 2: Проходим по всем полям структуры
	t := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := t.Field(i)
		fieldValue := value.Field(i)

		// Шаг 3: Считываем тег validate
		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		// Шаг 4: Распарсим тэг и выполним проверки
		rules := strings.Split(validateTag, ",") // Разделяем по запятой
		for _, rule := range rules {
			// Вызываем проверку для каждого правила
			err := applyRule(field.Name, fieldValue, rule)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{
					Field: field.Name,
					Err:   err,
				})
			}
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func applyRule(fieldName string, fieldValue reflect.Value, rule string) error {
	parts := strings.Split(rule, ":")
	switch parts[0] {
	case "min":
		minValue, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid min rule: %w", err)
		}
		if fieldValue.Kind() == reflect.Int {
			if fieldValue.Int() < int64(minValue) {
				return fmt.Errorf("field %q must be at least %d", fieldName, minValue)
			}
		} else if fieldValue.Kind() == reflect.String {
			if fieldValue.Len() < minValue {
				return fmt.Errorf("field %q must be at least %d characters long", fieldName, minValue)
			}
		}
	case "max":
		maxValue, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid max rule: %w", err)
		}
		if fieldValue.Kind() == reflect.Int {
			if fieldValue.Int() > int64(maxValue) {
				return fmt.Errorf("field %q must be at most %d", fieldName, maxValue)
			}
		} else if fieldValue.Kind() == reflect.String {
			if fieldValue.Len() > maxValue {
				return fmt.Errorf("field %q must be at most %d characters long", fieldName, maxValue)
			}
		}
	case "in":
		options := strings.Split(parts[1], ",")
		valid := false
		for _, option := range options {
			if fieldValue.String() == option {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("field %q must be one of [%s]", fieldName, strings.Join(options, ","))
		}
	}
	
	return nil
}
