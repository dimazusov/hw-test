package main //nolint:golint,stylecheck

import (
	"errors"
	"regexp"
	"strconv"
)

var ErrInvalidString = errors.New("invalid string")

func isValid(str string) bool {
	validRules := []bool{
		regexp.MustCompile(`^[0-9]`).MatchString(str),
		regexp.MustCompile(`[0-9]{2}`).MatchString(str),
	}

	for _, rule := range validRules {
		if rule {
			return false
		}
	}

	return true
}

func isNeedUnpack(re *regexp.Regexp, str string) bool {
	return re.MatchString(str)
}

func unpackPart(re *regexp.Regexp, str string) (string, error) {
	part := re.FindString(str)

	digit := part[len(part)-1:]
	count, err := strconv.ParseInt(digit, 10, 64)
	if err != nil {
		return "", err
	}

	unpackedPart := ""
	unpackingStr := part[:len(part)-1]
	for i := 0; i < int(count); i++ {
		unpackedPart += unpackingStr
	}

	findIndexes := re.FindIndex([]byte(str))
	if len(findIndexes) > 0 {
		str = str[:findIndexes[0]] + unpackedPart + str[findIndexes[1]:]
	}

	if isNeedUnpack(re, str) {
		if str, err = unpackPart(re, str); err != nil {
			return "", err
		}
	}

	return str, nil
}

func Unpack(str string) (result string, err error) {
	if !isValid(str) {
		return "", ErrInvalidString
	}

	unpackPatterns := []string{
		"[\\\\]{1}[a-z]{1}[0-9]{1}",
		"[a-z]{1}[0-9]{1}",
	}

	for _, patern := range unpackPatterns {
		var re = regexp.MustCompile(patern)

		if isNeedUnpack(re, str) {
			if str, err = unpackPart(re, str); err != nil {
				return "", ErrInvalidString
			}
		}
	}

	return str, nil
}
