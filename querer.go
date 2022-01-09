package querer

import (
	"bytes"
)

type Querer interface {
	Builder
}

func New(TableName string) Querer {
	return &querer{
		tableName:     TableName,
		buf:           &bytes.Buffer{},
		queryPosition: 1,
	}
}

type querer struct {
	tableName string
	// these will reset to zero value after build
	action        Action
	fields        []string // update or insert
	conditions    map[string]OperatorType
	limit         int
	offset        int
	order         OrderBy
	buf           *bytes.Buffer
	queryPosition int
	data          []interface{}
}

func (q *querer) reset() {
	q.buf.Reset()
	q.action = ActionNONE
	q.limit = 0
	q.offset = 0
	q.order = nil
	q.conditions = nil
	q.queryPosition = 1
	q.data = nil
	q.fields = nil
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
