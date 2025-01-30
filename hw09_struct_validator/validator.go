package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrUnsupportedInputType = errors.New("unsupported input type")
	ErrInvalidTagSyntax     = fmt.Errorf("tag must contain %q", tagDefinder)
)

type UnsupportedTagRuleError struct {
	rule string
}

func (err UnsupportedTagRuleError) Error() string {
	return fmt.Sprintf("tag rule %q isn't supported", err.rule)
}

const (
	validationTag string = "validate"
	tagDivider    string = "|"
	tagDefinder   string = ":"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var result strings.Builder
	result.WriteString("errors while validation:\n")
	for _, err := range v {
		result.WriteString(fmt.Sprintf("%v\t%v\n", err.Field, err.Err))
	}
	return result.String()
}

func Validate(v interface{}) error {
	vValue := reflect.ValueOf(v)
	vType := vValue.Type()

	if vValue.Kind() != reflect.Struct {
		return ErrUnsupportedInputType
	}

	validator, err := ParseRules(vType)
	if err != nil {
		return err
	}

	return validator.Validate(vValue)
}

type Validator struct {
	errors           ValidationErrors
	structValidators map[string]Validator
	intRules         map[string][]IntRule
	stringRules      map[string][]StringRule
	isSlice          map[string]bool
}

func ParseRules(t reflect.Type) (Validator, error) {
	validator := Validator{
		make(ValidationErrors, 0),
		make(map[string]Validator, 0),
		make(map[string][]IntRule, 0),
		make(map[string][]StringRule, 0),
		make(map[string]bool, 0),
	}
	var err error

	for _, field := range reflect.VisibleFields(t) {
		fieldTag, found := field.Tag.Lookup(validationTag)

		kind := field.Type.Kind()
		validator.isSlice[field.Name] = false
		if kind == reflect.Slice {
			kind = field.Type.Elem().Kind()
			validator.isSlice[field.Name] = true
		}

		if !found && kind != reflect.Struct {
			continue
		}

		switch kind { //nolint
		case reflect.Struct:
			if validator.isSlice[field.Name] {
				validator.structValidators[field.Name], err = ParseRules(field.Type.Elem())
			} else {
				validator.structValidators[field.Name], err = ParseRules(field.Type)
			}
		case reflect.Int:
			validator.intRules[field.Name], err = ParseIntRules(fieldTag)
		case reflect.String:
			validator.stringRules[field.Name], err = ParseStringRules(fieldTag)
		default:
			return validator, fmt.Errorf("unsupported type %q", kind)
		}
		if err != nil {
			return validator, fmt.Errorf("field %q has invalid tag %w", field.Name, err)
		}
	}
	return validator, nil
}

func (v *Validator) Validate(value reflect.Value) error {
	v.errors = make(ValidationErrors, 0)
	v.validateStructFields(value)
	v.validateIntFields(value)
	v.validateStringFields(value)
	return v.errors
}

func (v *Validator) validateStructFields(value reflect.Value) {
	for fieldName, validator := range v.structValidators {
		field := value.FieldByName(fieldName)
		if v.isSlice[fieldName] {
			for elemIndex := 0; elemIndex < field.Len(); elemIndex++ {
				v.checkStructField(fieldName, field.Index(elemIndex), validator)
			}
		} else {
			v.checkStructField(fieldName, field, validator)
		}
	}
}

func (v *Validator) validateIntFields(value reflect.Value) {
	for fieldName, rule := range v.intRules {
		field := value.FieldByName(fieldName)
		if v.isSlice[fieldName] {
			for elemIndex := 0; elemIndex < field.Len(); elemIndex++ {
				v.checkIntField(fieldName, field.Index(elemIndex), rule)
			}
		} else {
			v.checkIntField(fieldName, field, rule)
		}
	}
}

func (v *Validator) validateStringFields(value reflect.Value) {
	for fieldName, rule := range v.stringRules {
		field := value.FieldByName(fieldName)
		if v.isSlice[fieldName] {
			for elemIndex := 0; elemIndex < field.Len(); elemIndex++ {
				v.checkStringField(fieldName, field.Index(elemIndex), rule)
			}
		} else {
			v.checkStringField(fieldName, field, rule)
		}
	}
}

func (v *Validator) checkStructField(fieldName string, field reflect.Value, validator Validator) {
	err := validator.Validate(field).(ValidationErrors) //nolint
	if len(err) > 0 {
		v.errors = append(v.errors, ValidationError{fieldName, err})
	}
}

func (v *Validator) checkIntField(fieldName string, field reflect.Value, rules []IntRule) {
	if !field.CanInt() {
		v.errors = append(v.errors, ValidationError{fieldName, fmt.Errorf("wrong type %T", field.Kind())})
		return
	}
	value := field.Int()
	for _, rule := range rules {
		err := rule(value)
		if err != nil {
			v.errors = append(v.errors, ValidationError{fieldName, err})
		}
	}
}

func (v *Validator) checkStringField(fieldName string, field reflect.Value, rules []StringRule) {
	value := field.String()
	for _, rule := range rules {
		err := rule(value)
		if err != nil {
			v.errors = append(v.errors, ValidationError{fieldName, err})
		}
	}
}
