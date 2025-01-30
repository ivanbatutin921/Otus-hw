package hw09structvalidator

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrStringLen    = errors.New("string not equal target len")
	ErrStringRegexp = errors.New("string doesn't match regexp")
	ErrStringIn     = errors.New("string is not in the values list")
)

type (
	StringRule       = func(string) error
	StringRuleGetter = func(string) (StringRule, error)
)

type stringRule struct{}

func ParseStringRules(tag string) ([]StringRule, error) {
	rules := strings.Split(tag, tagDivider)
	result := make([]StringRule, 0)
	for _, rule := range rules {
		checkFunc, err := ParseStringRule(rule)
		if err != nil {
			return nil, err
		}
		result = append(result, checkFunc)
	}
	return result, nil
}

func ParseStringRule(rule string) (StringRule, error) {
	funcName, controlValue, found := strings.Cut(rule, ":")
	if !found {
		return nil, ErrInvalidTagSyntax
	}

	sr := stringRule{}
	ruleGetter, found := sr.getRuleGetter(funcName)
	if !found {
		return nil, UnsupportedTagRuleError{funcName}
	}

	checkFunc, err := ruleGetter(controlValue)
	if err != nil {
		return nil, fmt.Errorf("erorr while getting rule: %w", err)
	}
	return checkFunc, nil
}

func (sr stringRule) getRuleGetter(funcName string) (StringRuleGetter, bool) {
	switch funcName {
	case "len":
		return sr.getLenRule, true
	case "regexp":
		return sr.getRegexpRule, true
	case "in":
		return sr.getInRule, true
	default:
		return nil, false
	}
}

func (sr stringRule) getLenRule(controlValue string) (StringRule, error) {
	controlLen, err := strconv.Atoi(controlValue)
	if err != nil {
		return nil, err
	}

	return func(value string) error {
		if len(value) == controlLen {
			return nil
		}
		return ErrStringLen
	}, nil
}

func (sr stringRule) getRegexpRule(controlValue string) (StringRule, error) {
	re, err := regexp.Compile(controlValue)
	if err != nil {
		return nil, err
	}

	return func(value string) error {
		if re.MatchString(value) {
			return nil
		}
		return ErrStringRegexp
	}, nil
}

func stringSetFromSlice(sl []string) map[string]struct{} {
	result := make(map[string]struct{}, len(sl))
	for _, value := range sl {
		result[value] = struct{}{}
	}
	return result
}

func (sr stringRule) getInRule(controlValue string) (StringRule, error) {
	controlValues := stringSetFromSlice(strings.Split(controlValue, ","))

	return func(value string) error {
		if _, ok := controlValues[value]; !ok {
			return ErrStringIn
		}
		return nil
	}, nil
}
