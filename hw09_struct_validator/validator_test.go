package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
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
	ExtraID struct {
		ID []int `validate:"min:1|max:11"`
	}
	NotSupported struct {
		Users map[UserRole]User `validate:"len:5"`
	}
)

func TestValidate(t *testing.T) {
	TestPositiveValidate(t)
	TestNegativeValidate(t)
	TestOthers(t)
}

func TestPositiveValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			Response{200, "valid response test"},
			nil,
		},
		{
			Response{404, "valid response test"},
			nil,
		},
		{
			Response{500, "valid response test"},
			nil,
		},
		{
			App{"1.0.0"},
			nil,
		},
		{
			User{
				"123456789012345678901234567890123456",
				"hw09",
				20,
				"home@work.com",
				UserRole("admin"),
				[]string{"12345678901", "12345678901"},
				nil,
			},
			nil,
		},
		{
			Token{
				[]byte("header"),
				[]byte("payload"),
				[]byte("signature"),
			},
			nil,
		},
		{
			ExtraID{
				[]int{1, 11},
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Nil(t, err)
		})
	}
}

func TestNegativeValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			Response{1, "invalid response test"},
			ErrValidation,
		},
		{
			App{"1.0.0.0.0"},
			ErrValidation,
		},
		{
			User{
				"",
				"",
				0,
				"home@work.com",
				UserRole("admin"),
				[]string{"12345678901", "12345678901"},
				nil,
			},
			ErrValidation,
		},
		{
			NotSupported{},
			ErrNotSupportedType,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			errs := ValidationErrors{}

			err := Validate(tt.in)
			if errors.As(err, &errs) {
				for _, e := range errs {
					require.Error(t, e.Err)
					require.Equal(t, tt.expectedErr, e.Err)
				}
			}
		})
	}
}

func TestOthers(t *testing.T) {
	require.Equal(t, ErrNotStruct, Validate(5))
	require.Nil(t, Validate(nil))
}
