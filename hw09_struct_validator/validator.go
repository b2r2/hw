package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

var (
	ErrNotStruct              = errors.New("received not struct")
	ErrValidationTag          = errors.New("incorrect validation error")
	ErrNotSupportedType       = errors.New("type is not supported")
	ErrNotSupportedValidation = errors.New("tag is not supported tag name")
	ErrValidation             = errors.New("field failed validation")
)

const (
	validate = "validate"
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errs := strings.Builder{}
	for _, err := range v {
		errs.WriteString(err.Field)
		errs.WriteString(":")
		errs.WriteString(err.Err.Error())
		errs.WriteString("\n")
	}
	return errs.String()
}

func Validate(iv interface{}) error {
	if iv == nil {
		return nil
	}
	v := reflect.Indirect(reflect.ValueOf(iv))
	if v.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	var errs ValidationErrors
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get(validate)
		if tag == "" {
			continue
		}

		typeFiled := v.Type().Field(i)
		containers, err := handleTag(tag)
		if err != nil {
			errs = append(errs, ValidationError{typeFiled.Name, err})
			continue
		}
		if len(containers) == 0 {
			continue
		}
		err = handleValues(containers, v.Field(i))
		if err != nil {
			errs = append(errs, ValidationError{typeFiled.Name, err})
		}
	}
	return errs
}

type rule struct {
	name, value string
}

func handleTag(tag string) ([]rule, error) {
	var rules []rule //nolint:prealloc
	for _, r := range strings.Split(tag, "|") {
		v := strings.Split(r, ":")
		if len(v) != 2 {
			return nil, ErrValidationTag
		}
		rules = append(rules, rule{v[0], v[1]})
	}
	return rules, nil
}

func handleValues(containers []rule, value reflect.Value) error {
	switch value.Kind() {
	case reflect.Int:
		return handleInt(int(value.Int()), containers)
	case reflect.String:
		return handleString(value.String(), containers)
	case reflect.Slice:
		switch v := value.Interface().(type) {
		case []string:
			var values []string
			values = append(values, v...)
			for i := range values {
				if err := handleString(values[i], containers); err != nil {
					return err
				}
			}

		case []int:
			var values []int
			values = append(values, v...)
			for i := range values {
				if err := handleInt(values[i], containers); err != nil {
					return err
				}
			}
		}
	default:
		return ErrNotSupportedType
	}
	return nil
}

func handleInt(i int, rules []rule) error {
	for _, r := range rules {
		switch r.name {
		case "in":
			s := strings.Split(r.value, ",")
			min, err := strconv.Atoi(s[0])
			if err != nil {
				return err
			}
			max, err := strconv.Atoi(s[2])
			if err != nil {
				return err
			}
			if min > i || i > max {
				return ErrValidation
			}

		case "min":
			min, err := strconv.Atoi(r.value)
			if err != nil {
				return err
			}
			if i < min {
				return ErrValidation
			}
		case "max":
			max, err := strconv.Atoi(r.value)
			if err != nil {
				return err
			}
			if i > max {
				return ErrValidation
			}
		default:
			return ErrNotSupportedValidation
		}
	}
	return nil
}

func handleString(str string, rules []rule) error {
	for _, r := range rules {
		switch r.name {
		case "in":
			ss := strings.Split(r.value, ",")
			for _, s := range ss {
				if s == str {
					return nil
				}
			}
		case "len":
			v, err := strconv.Atoi(r.value)
			if err != nil {
				return err
			}
			if v != len(str) {
				return ErrValidation
			}

		case "regexp":
			s := strings.ReplaceAll(r.value, "\\\\", `\`)
			r := regexp.MustCompile(s)
			if !r.MatchString(str) {
				return ErrValidation
			}
		default:
			return ErrNotSupportedValidation
		}
	}
	return nil
}
