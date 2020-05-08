package main

import (
	"fmt"
	"reflect"
	"strings"
)

var Filters = map[string]interface{}{
	// テスト用の関数
	"hasTag": func(post *Post, arg1 string, arg2 string) bool {
		value, exists := post.Tags[arg1]
		if !exists {
			return false
		}
		return value == arg2
	},
	"includeWordInName": func(post *Post, arg1 string) bool {
		return strings.Contains(post.Name, arg1)
	},
	"includeWordInText": func(post *Post, arg1 string) bool {
		return strings.Contains(post.Text, arg1)
	},
}

func validateFilter(v interface{}) error {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Func {
		return fmt.Errorf("ErrInvalidFilter: Not function")
	}
	if t.NumIn() < 1 {
		return fmt.Errorf("ErrInvalidFilter: Number of arguments must be larger than 1")
	}
	if t.In(0).String() != "*Post" {
		return fmt.Errorf("ErrInvalidFilter: Arg0 is not *Post")
	}
	if t.NumOut() != 1 {
		return fmt.Errorf("ErrInvalidFilter: Number of return values must be 1")
	}
	if t.Out(0).Kind() != reflect.Bool {
		return fmt.Errorf("ErrInvalidFilter: Return value must be bool")
	}
	return nil
}

func ValidateFilters() error {
	for _, v := range Filters {
		if err := validateFilter(v); err != nil {
			return err
		}
	}
	return nil
}
