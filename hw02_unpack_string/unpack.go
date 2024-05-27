package hw02unpackstring

import (
	"errors"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}
	if str[0] < '0' || str[0] > '9' {
		return "", nil
	}
	for i := 0; i < len(str); i++ {
		
	}
	return "", nil
}
