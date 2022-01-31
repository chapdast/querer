package querer

import "fmt"

type Option func(q *querer) error

// SelectFields add fields name to be selected from query if empty, all fields will get selected
func SelectFields(fields ...interface{}) Option {
	var selection []string
	for _, s := range fields {
		switch data := s.(type) {
		case string:
			selection = append(selection, data)
		case StringFunc:
			selection = append(selection, data())
		}
	}

	if len(fields) == 0 {
		selection = append(selection, "*")
	}
	return query_select(selection)
}

func Select(fields []string) Option {
	return query_select(fields)
}
func query_select(f []string) Option {
	return func(q *querer) error {
		q.action = ActionSelect
		for _, field := range f {
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

type Joiner int

const (
	AND Joiner = iota
	OR
)

type Conditional struct {
	Field     string
	Operation OperatorType
	Value     interface{}
	ANDOR Joiner
	Negate    bool
}

func (c Conditional) Operator() OperatorType {
	if c.Negate {
		return c.Operation.Flip()
	}
	return c.Operation
}

func Where(condition Conditional) Option {
	return func(q *querer) error {
		if q.conditions == nil {
			q.conditions = make(map[string][]Conditional)
		}
		q.conditions[condition.Field] = append(q.conditions[condition.Field], condition)

		q.data = append(q.data, condition.Value)
		return nil
	}
}
func Limit(limit int) Option {
	return func(q *querer) error {
		q.limit = limit
		//q.data = append(q.data, limit)
		return nil
	}
}
func Offset(offset int) Option {
	return func(q *querer) error {
		q.offset = offset
		//q.data = append(q.data, offset)
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
