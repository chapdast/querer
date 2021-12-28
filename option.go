package querer

import "fmt"

type Option func(q *querer) error

func Select() Option {
	return func(q *querer) error {
		q.action = ActionSelect
		return nil
	}
}
func Insert() Option {
	return func(q *querer) error {
		q.action = ActionInsert
		return nil
	}
}
func Update() Option {
	return func(q *querer) error {
		q.action = ActionUpdate
		return nil
	}
}
func Delete() Option {
	return func(q *querer) error {
		q.action = ActionDelete
		return nil
	}
}
func Where(c map[string]OperatorType) Option {
	return func(q *querer) error {
		q.conditions = c
		return nil
	}
}
func Limit(limit int) Option {
	return func(q *querer) error {
		q.limit = limit
		return nil
	}
}
func Offset(offset int) Option {
	return func(q *querer) error {
		q.offset = offset
		return nil
	}
}

func WithOrder(field string, desc bool) Option {
	return func(q *querer) error {
		if len(field) == 0 {
			return fmt.Errorf("can not order by empty field name")
		}
		if q.notInFields(field) {
			return fmt.Errorf("unkown field to order with")
		}
		q.order = &orderBy{
			field: field,
			desc:  desc,
		}
		return nil
	}
}
