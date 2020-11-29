package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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
		in   interface{}
		errs []ValidationError
	}{
		{
			in: User{
				ID:     "asdgasdgawegawegasdvdsfasfasfjwdssoi",
				Name:   "test",
				Age:    20,
				Email:  "vasiliy@mail.ru",
				Role:   "stuff",
				Phones: []string{"79524399025", "75924858843"},
				meta:   []byte("[]"),
			},
			errs: []ValidationError{},
		},
		{
			in: App{
				Version: "23141",
			},
			errs: []ValidationError{},
		},
		{
			in: Token{
				Header: []byte("fasdgasdg"),
			},
			errs: []ValidationError{},
		},
		{
			in: Response{
				Code: 200,
			},
			errs: []ValidationError{
				{
					Field: "Response",
					Err:   ErrValidatorIsNotValid,
				},
			},
		},
		{
			in: App{
				Version: "2314",
			},
			errs: []ValidationError{
				{
					Field: "App",
					Err: ErrLengthNotEqual,
				},
			},
		},
		{
			in: User{
				ID:     "asdgasdgawegawegasdvdsfasfasfjwdssoi",
				Name:   "test",
				Age:    20,
				Email:  "@tes@t@mail.ru",
				Role:   "stuffffffffffff",
				Phones: []string{"79524399025", "75924858843"},
				meta:   []byte("[]"),
			},
			errs: []ValidationError{
				{
					Field: "User",
					Err: ErrRegexNotFound,
				},
				{
					Field: "Role",
					Err: ErrInStringNotFound,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			errs := Validate(tt.in)
			require.Equal(t, len(errs), len(tt.errs))

			for i, _ := range errs {
				require.True(t, errors.Is(errs[i].Err, tt.errs[i].Err))
			}
		})
	}
}

func TestCreateValidatorFromString(t *testing.T) {
	validator, err := CreateValidatorFromString("min:10", reflect.Int)
	require.NotNil(t, validator)
	require.Nil(t, err)
}

func TestValidateValidatorValue(t *testing.T) {
	err := ValidateValidatorValue("in", "123,42", reflect.String)
	require.Nil(t, err)

	err = ValidateValidatorValue("max", "123", reflect.String)
	require.Nil(t, err)

	err = ValidateValidatorValue("max", "str123", reflect.String)
	require.NotNil(t, err)
}

func TestCreateValidatorsFromString(t *testing.T) {
	validators, errs := CreateValidatorsFromString(`min:12|max:20`, reflect.Int)
	require.NotNil(t, validators)
	require.Nil(t, errs)

	require.Equal(t, "min", validators[0].Name)
	require.Equal(t, "12", validators[0].Value)

	require.Equal(t, "max", validators[1].Name)
	require.Equal(t, "20", validators[1].Value)
}

func TestValidatorValidate(t *testing.T) {
	validators, errs := CreateValidatorsFromString(`min:12|max:20`, reflect.Int)
	require.Nil(t, errs)
	value := 12
	for _, validator := range validators {
		err := validator.Validate(value)
		require.Nil(t, err)
	}

	validators, errs = CreateValidatorsFromString(`min:12`, reflect.Int)
	require.Nil(t, errs)
	value = 11

	err := validators[0].Validate(value)
	require.NotNil(t, err)
}
