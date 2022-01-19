package querer

import "fmt"

type StringFunc func() string

func Coalesce(field string, defValue interface{}, cast ...string) StringFunc {
	return func() string {
		field := fmt.Sprintf("COALESCE(%s, '%s')", field, defValue)
		if len(cast) == 1 {
			field = fmt.Sprintf("%s::%s", field, cast[0])
		}
		return field
	}
}

func Extract(field string, what string, cast ...string) StringFunc {
	return func() string {
		field := fmt.Sprintf("EXTRACT(%s from %s)", what, field)
		if len(cast) == 1 {
			field = fmt.Sprintf("%s::%s", field, cast[0])
		}
		return field
	}
}
