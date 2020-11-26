package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"github.com/pkg/errors"

	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Validation struct {
	validators []Validator
	errors     []ValidationError
}

type ValidationError struct {
	Field string
	Err   error
}

var ErrNotStruct = errors.New("value is not a struct")
var ErrTypeNotSupport = errors.New("type do not support for validate")
var ErrValidatorIsNotAlowed = errors.New("validator is not allowed")
var ErrValidatorIsNotValid = errors.New("validator is not valid")
var ErrValidatorNameIsNotAlowed = errors.New("validator name is not allowed")
var ErrLengthNotEqual = errors.New("length not equal")
var ErrRegexStringNotFound = errors.New("string not found")
var ErrInStringNotFound = errors.New("string not found")
var ErrMax = errors.New("min value")
var ErrMin = errors.New("max value")
var ErrIn = errors.New("in value")

var ErrMessageValidatorValueIsNotAlowed = "validator value is not allowed"
var ErrMessageValidatorMustContaintsTwoValues = "validator value must containts two values"
var ErrMessageValidatorMustBeInteger = "validator value must be integer"
var ErrMessageValidatorMustHaveOneOrMoreSymbols = "validator value must have one or more symbols"
var ErrMessageValidatorErrror = "validator error"

const validateTag = "validate"

const minValidator = "min"
const maxValidator = "max"
const inValidator = "in"
const lenValidator = "len"
const regexpValidator = "regexp"

const stringDefaultValue = ""
const intDefaultValue = 0

var validateTypes = map[reflect.Kind][]string{
	reflect.Int:    {minValidator, maxValidator, inValidator},
	reflect.String: {lenValidator, regexpValidator, inValidator},
}

type Validator struct {
	Name  string
	Value string
	Kind  reflect.Kind
}

func (m Validator) Validate(value interface{}) error {
	switch m.Kind {
	case reflect.String:
		v := reflect.ValueOf(value).String()

		if v == stringDefaultValue {
			return nil
		}

		switch m.Name {
		case lenValidator:
			validatorValue, err := strconv.ParseInt(m.Value, 10, 64)
			if err != nil {
				return err
			}

			if len(v) != int(validatorValue) {
				return errors.Wrap(ErrLengthNotEqual, fmt.Sprintf("given length: %d, expected length: %d", len(v), validatorValue))
			}

			return nil
		case regexpValidator:
			r, err := regexp.Compile(m.Value)
			if err != nil {
				return errors.Wrap(err, ErrMessageValidatorErrror)
			}

			if !r.MatchString(v) {
				return ErrRegexStringNotFound
			}

			return nil
		case inValidator:
			for _, val := range strings.Split(m.Value, ",") {
				if val == v {
					return nil
				}
			}

			return ErrInStringNotFound
		}
		break
	case reflect.Int:
		v := reflect.ValueOf(value).Int()

		if v == intDefaultValue {
			return nil
		}

		switch m.Name {
		case minValidator:
			validatorValue, err := strconv.ParseInt(m.Value, 10, 64)
			if err != nil {
				return err
			}

			if v < validatorValue {
				return ErrMin
			}

			return nil
		case maxValidator:
			validatorValue, err := strconv.ParseInt(m.Value, 10, 64)
			if err != nil {
				return err
			}

			if v > validatorValue {
				return ErrMax
			}

			return nil
		case inValidator:
			res := strings.Split(m.Value, ",")

			from, err := strconv.ParseInt(res[0], 10, 64)
			if err != nil {
				return err
			}

			to, err := strconv.ParseInt(res[1], 10, 64)
			if err != nil {
				return err
			}

			if v < from || v > to {
				return ErrIn
			}

			return nil
		}
		break
	}

	return nil
}

func Validate(v interface{}) (errs []ValidationError) {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Struct {
		errs = append(errs, ValidationError{
			Field: "none",
			Err:   ErrNotStruct,
		})
		return errs
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.CanInterface() {
			continue
		}

		tag := val.Type().Field(i).Tag
		fieldName := val.Type().Name()

		validators, errsValidators := getValidatorsForField(field, tag, fieldName)
		if errsValidators != nil {
			errs = append(errs, errsValidators...)
			continue
		}

		fieldValues := getFieldValuesForValidation(field)

		errsFieldVals := validateFieldValues(fieldName, fieldValues, validators)
		if errsFieldVals != nil {
			errs = append(errs, errsFieldVals...)
			continue
		}
	}

	return errs
}

func getFieldValuesForValidation(field reflect.Value) []interface{} {
	fieldValues := []interface{}{}

	if field.Kind() == reflect.Slice {
		for i := 0; i < field.Len(); i++ {
			fieldValues = append(fieldValues, field.Index(i).Interface())
		}
	} else {
		fieldValues = append(fieldValues, field.Interface())
	}

	return fieldValues
}

func validateFieldValues(fieldName string, fieldValues []interface{}, validators []*Validator) (errs []ValidationError) {
	for _, v := range validators {
		for _, value := range fieldValues {
			err := v.Validate(value)
			if err != nil {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   err,
				})
			}
		}
	}
	return errs
}

func getValidatorsForField(field reflect.Value, tag reflect.StructTag, fieldName string) (validators []*Validator, errs []ValidationError) {
	fieldKind := field.Kind()

	tagValue := tag.Get(validateTag)
	if len(tagValue) == 0 {
		return nil, nil
	}

	if fieldKind == reflect.Slice {
		if field.Len() == 0 {
			return nil, nil
		}
		fieldKind = field.Index(0).Kind()
	}

	_, ok := validateTypes[fieldKind]
	if !ok {
		errs = append(errs, ValidationError{
			Field: fieldName,
			Err:   ErrTypeNotSupport,
		})
		return nil, errs
	}

	validators, validatorErrors := CreateValidatorsFromString(tagValue, fieldKind)
	if validatorErrors != nil {
		for _, err := range validatorErrors {
			errs = append(errs, ValidationError{
				Field: fieldName,
				Err:   err,
			})
		}
		return nil, errs
	}

	return validators, errs
}

func isValidatorsAllow(allowValidators, validators []string) error {
	for _, validator := range validators {
		isAllow := false
		for _, allowValidator := range allowValidators {
			if allowValidator == validator {
				isAllow = true
				break
			}
		}

		if !isAllow {
			return ErrValidatorIsNotAlowed
		}
	}

	return nil
}

func CreateValidatorsFromString(str string, kind reflect.Kind) (validators []*Validator, errs []error) {
	validatorsStrings := strings.Split(str, "|")
	for _, validatorString := range validatorsStrings {
		validator, err := CreateValidatorFromString(validatorString, kind)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		validators = append(validators, validator)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return validators, nil
}

func CreateValidatorFromString(str string, kind reflect.Kind) (*Validator, error) {
	values := strings.Split(str, ":")
	if len(values) != 2 {
		return nil, ErrValidatorIsNotValid
	}

	if err := ValidateValidatorName(values[0], kind); err != nil {
		return nil, err
	}

	if err := ValidateValidatorValue(values[0], values[1], kind); err != nil {
		return nil, err
	}

	return &Validator{
		Name:  values[0],
		Value: values[1],
		Kind:  kind,
	}, nil
}

func ValidateValidatorName(validatorName string, kind reflect.Kind) error {
	allowTypes, _ := validateTypes[kind]

	isAllow := false
	for _, allowName := range allowTypes {
		if allowName == validatorName {
			isAllow = true
			break
		}
	}

	if !isAllow {
		return ErrValidatorNameIsNotAlowed
	}

	return nil
}

func ValidateValidatorValue(validatorName, validatorValue string, kind reflect.Kind) error {
	switch validatorName {
	case minValidator:
		_, err := strconv.ParseInt(validatorValue, 0, 64)
		if err != nil {
			return errors.Wrap(err, ErrMessageValidatorValueIsNotAlowed)
		}
		break
	case maxValidator:
		_, err := strconv.ParseInt(validatorValue, 0, 64)
		if err != nil {
			return errors.Wrap(err, ErrMessageValidatorValueIsNotAlowed)
		}
		break
	case inValidator:
		if kind == reflect.String {
			res := strings.Split(validatorValue, ",")
			for _, v := range res {
				if v == "" {
					return errors.Wrap(ErrValidatorIsNotValid, ErrMessageValidatorMustHaveOneOrMoreSymbols)
				}
			}
		} else {
			res := strings.Split(validatorValue, ",")
			if len(res) != 2 {
				return errors.Wrap(ErrValidatorIsNotValid, ErrMessageValidatorMustContaintsTwoValues)
			}

			_, err := strconv.ParseInt(res[0], 0, 64)
			if err != nil {
				return errors.Wrap(err, ErrMessageValidatorMustBeInteger)
			}

			_, err = strconv.ParseInt(res[1], 0, 64)
			if err != nil {
				return errors.Wrap(err, ErrMessageValidatorMustBeInteger)
			}
		}
		break
	case lenValidator:
		_, err := strconv.ParseInt(validatorValue, 0, 64)
		if err != nil {
			return errors.Wrap(err, ErrMessageValidatorValueIsNotAlowed)
		}
		break
	case regexpValidator:
		break
	default:
		return ErrValidatorNameIsNotAlowed
	}

	return nil
}
