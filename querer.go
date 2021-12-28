package querer

import (
	"bytes"
)

type Querer interface {
	Builder
}

func New(TableName string, fields []string) Querer {
	return &querer{
		tableName:     TableName,
		fields:        fields,
		buf:           &bytes.Buffer{},
		queryPosition: 1,
	}
}

type querer struct {
	tableName string
	fields    []string // update or insert
	// these will reset to zero value after build
	action        Action
	conditions    map[string]OperatorType
	limit         int
	offset        int
	order         OrderBy
	buf           *bytes.Buffer
	queryPosition int
}

func (q *querer) reset() {
	q.buf.Reset()
	q.action = ActionNONE
	q.limit = 0
	q.offset = 0
	q.order = nil
	q.conditions = nil
	q.queryPosition = 1
}

func (q *querer) loadOpts(opts []Option) error {
	for _, opt := range opts {
		if opt != nil {
			if err := opt(q); err != nil {
				return err
			}
		}
	}
	return nil
}

func (q *querer) notInFields(s string) bool {
	for _, f := range q.fields {
		if s == f {
			return false
		}
	}
	return true
}
