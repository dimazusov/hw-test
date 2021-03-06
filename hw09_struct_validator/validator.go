package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type ValidationError struct {
	Field string
	Err   error
}

var ErrNotStruct = errors.New("value is not a struct")
var ErrTypeNotSupport = errors.New("type do not support for validate")
var ErrValidatorIsNotValid = errors.New("validator is not valid")
var ErrValidatorNameIsNotAlowed = errors.New("validator name is not allowed")
var ErrLengthNotEqual = errors.New("length not equal")
var ErrRegexNotFound = errors.New("not found by regex")
var ErrInStringNotFound = errors.New("value not found in allow")
var ErrMax = errors.New("min value")
var ErrMin = errors.New("max value")

var ErrMessageRegexValues = "value must be in %s"
var ErrMessageLengthNotEqual = "given length: %d, expected length: %d"
var ErrMessageMaxValue = "%d is more than %d"
var ErrMessageMinValue = "%d is less than %d"
var ErrMessageInDigits = "%d value is not found in validator"
var ErrMessageInStrings = "%s value is not found in validator"
var ErrMessageValidatorIsNotAlowed = "%s validator is not allowed"
var ErrMessageNotStruct = "%s is not struct"
var ErrMessageWrongValidatorValue = "%s wrong validator value"
var ErrMessageValidatorValueIsNotAlowed = "validator value is not allowed"
var ErrMessageValidatorMustContaintsTwoValues = "validator value must contains two values"
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
	switch m.Kind { //nolint:exhaustive
	case reflect.String:
		v := reflect.ValueOf(value).String()

		return m.validateStringValue(v)
	case reflect.Int:
		v := reflect.ValueOf(value).Int()

		return m.validateIntValue(v)
	}

	return nil
}

func (m Validator) validateStringValue(v string) error {
	if v == stringDefaultValue {
		return nil
	}

	switch m.Name {
	case lenValidator:
		validatorValue, err := strconv.ParseInt(m.Value, 10, 64)
		if err != nil {
			return errors.Wrap(err, ErrMessageWrongValidatorValue)
		}

		if len(v) != int(validatorValue) {
			return errors.Wrap(ErrLengthNotEqual, fmt.Sprintf(ErrMessageLengthNotEqual, len(v), validatorValue))
		}

		return nil
	case regexpValidator:
		r, err := regexp.Compile(m.Value)
		if err != nil {
			return errors.Wrap(err, ErrMessageValidatorErrror)
		}

		if !r.MatchString(v) {
			return errors.Wrap(ErrRegexNotFound, fmt.Sprintf(ErrMessageRegexValues, v, m.Value))
		}

		return nil
	case inValidator:
		for _, val := range strings.Split(m.Value, ",") {
			if val == v {
				return nil
			}
		}

		return errors.Wrap(ErrInStringNotFound, fmt.Sprintf(ErrMessageInStrings, m.Value))
	}

	return nil
}

func (m Validator) validateIntValue(v int64) error {
	if v == intDefaultValue {
		return nil
	}

	switch m.Name {
	case minValidator:
		validatorValue, err := strconv.ParseInt(m.Value, 10, 64)
		if err != nil {
			return errors.Wrap(err, ErrMessageWrongValidatorValue)
		}

		if v < validatorValue {
			return errors.Wrap(ErrMin, fmt.Sprintf(ErrMessageMinValue, v, validatorValue))
		}

		return nil
	case maxValidator:
		validatorValue, err := strconv.ParseInt(m.Value, 10, 64)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf(ErrMessageWrongValidatorValue, validatorValue))
		}

		if v > validatorValue {
			return errors.Wrap(ErrMax, fmt.Sprintf(ErrMessageMaxValue, v, validatorValue))
		}

		return nil
	case inValidator:
		res := strings.Split(m.Value, ",")

		from, err := strconv.ParseInt(res[0], 10, 64)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf(ErrMessageWrongValidatorValue, res[0]))
		}

		to, err := strconv.ParseInt(res[1], 10, 64)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf(ErrMessageWrongValidatorValue, res[1]))
		}

		if v < from || v > to {
			return errors.Wrap(ErrMin, fmt.Sprintf(ErrMessageInDigits, v))
		}

		return nil
	}

	return nil
}

func Validate(v interface{}) (errs []ValidationError) {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Struct {
		errs = append(errs, ValidationError{
			Field: "none",
			Err:   errors.Wrap(ErrNotStruct, fmt.Sprintf(ErrMessageNotStruct, val.Kind().String())),
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
		return nil, errors.Wrap(ErrValidatorIsNotValid, fmt.Sprintf(ErrMessageValidatorIsNotAlowed, str))
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
	allowTypes := validateTypes[kind]

	isAllow := false
	for _, allowName := range allowTypes {
		if allowName == validatorName {
			isAllow = true
			break
		}
	}

	if !isAllow {
		return errors.Wrap(ErrValidatorNameIsNotAlowed, fmt.Sprintf(ErrMessageValidatorIsNotAlowed, validatorName))
	}

	return nil
}

func ValidateValidatorValue(validatorName, validatorValue string, kind reflect.Kind) error {
	switch kind { //nolint:exhaustive
	case reflect.String:
		return ValidateStringValidator(validatorName, validatorValue)
	case reflect.Int:
		return ValidateIntValidator(validatorName, validatorValue)
	default:
		return errors.Wrap(ErrValidatorNameIsNotAlowed, fmt.Sprintf(ErrMessageValidatorIsNotAlowed, validatorName))
	}
}

func ValidateStringValidator(validatorName, validatorValue string) error {
	switch validatorName {
	case inValidator:
		res := strings.Split(validatorValue, ",")
		for _, v := range res {
			if v == "" {
				return errors.Wrap(ErrValidatorIsNotValid, ErrMessageValidatorMustHaveOneOrMoreSymbols)
			}
		}
		return nil
	case lenValidator:
		_, err := strconv.ParseInt(validatorValue, 0, 64)
		if err != nil {
			return errors.Wrap(err, ErrMessageValidatorValueIsNotAlowed)
		}
		return nil
	case regexpValidator:
		return nil
	}

	return errors.Wrap(ErrValidatorNameIsNotAlowed, fmt.Sprintf(ErrMessageValidatorIsNotAlowed, validatorName))
}

func ValidateIntValidator(validatorName, validatorValue string) error {
	switch validatorName {
	case minValidator:
		_, err := strconv.ParseInt(validatorValue, 0, 64)
		if err != nil {
			return errors.Wrap(err, ErrMessageValidatorValueIsNotAlowed)
		}
		return nil
	case maxValidator:
		_, err := strconv.ParseInt(validatorValue, 0, 64)
		if err != nil {
			return errors.Wrap(err, ErrMessageValidatorValueIsNotAlowed)
		}
		return nil
	case inValidator:
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
	return errors.Wrap(ErrValidatorNameIsNotAlowed, fmt.Sprintf(ErrMessageValidatorIsNotAlowed, validatorName))
}
