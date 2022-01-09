package querer

import "fmt"

type Option func(q *querer) error

func Select(fields []string) Option {
	return func(q *querer) error {
		q.action = ActionSelect
		for _, field := range fields {
			q.fields = append(q.fields, field)
		}
		return nil
	}
}
func Insert(fields []string, values []interface{}) Option {
	return func(q *querer) error {
		if len(fields) != len(values) {
			return fmt.Errorf("fields not match values")
		}
		q.action = ActionInsert
		for index, field := range fields {
			q.fields = append(q.fields, field)
			q.data = append(q.data, values[index])
		}
		return nil
	}
}
func Update(fields []string, values []interface{}) Option {
	return func(q *querer) error {
		if len(fields) != len(values) {
			return fmt.Errorf("fields not match values")
		}
		q.action = ActionUpdate
		for index, field := range fields {
			q.fields = append(q.fields, field)
			q.data = append(q.data, values[index])
		}
		return nil
	}
}
func Delete() Option {
	return func(q *querer) error {
		q.action = ActionDelete
		return nil
	}
}

type Conditional struct {
	Field     string
	Operation OperatorType
	Value     interface{}
}

func Where(condition Conditional) Option {
	return func(q *querer) error {
		if q.conditions == nil {
			q.conditions = make(map[string]OperatorType)
		}
		q.conditions[condition.Field] = condition.Operation
		q.data = append(q.data, condition.Value)
		return nil
	}
}
func Limit(limit int) Option {
	return func(q *querer) error {
		q.limit = limit
		q.data = append(q.data, limit)
		return nil
	}
}
func Offset(offset int) Option {
	return func(q *querer) error {
		q.offset = offset
		q.data = append(q.data, offset)
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
