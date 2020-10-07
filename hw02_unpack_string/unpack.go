package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}

	if unicode.IsDigit(rune(str[0])) {
		return "", ErrInvalidString
	}

	for curIndex, symbol := range str {
		if !unicode.IsDigit(symbol) {
			continue
		}

		nextIndex := curIndex + 1
		lastIndex := len(str) - 1

		if lastIndex != curIndex && unicode.IsDigit(rune(str[nextIndex])) {
			return "", ErrInvalidString
		}

		digit, err := strconv.Atoi(string(symbol))
		if err != nil {
			return "", err
		}

		builder := strings.Builder{}
		builder.Write([]byte(str[:curIndex-1]))

		for j := 0; j < digit; j++ {
			builder.Write([]byte{str[curIndex-1]})
		}

		builder.Write([]byte(str[curIndex+1:]))

		str, err = Unpack(builder.String())
		if err != nil {
			return "", err
		}

		break
	}

	return str, nil
}
