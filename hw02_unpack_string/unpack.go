package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var b strings.Builder
	var count int
	var err error

	if str == "" {
		return "", nil
	}

	for i := 0; i < len(str); i++ {
		if str[i] >= '0' && str[i] <= '9' {
			return "", ErrInvalidString
		}

		if i+1 < len(str) && str[i+1] >= '0' && str[i+1] <= '9' {
			count, err = strconv.Atoi(string(str[i+1]))
			if err != nil {
				return "", err
			}
			b.WriteString(strings.Repeat(string(str[i]), count))
			i++
		} else {
			b.WriteString(string(str[i]))
		}
	}

	return b.String(), nil
}
