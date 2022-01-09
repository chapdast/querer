package querer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type Builder interface {
	Build(opts ...Option) (string, error)
}

func (q *querer) Build(opts ...Option) (string, error) {
	defer q.reset()
	if err := q.loadOpts(opts); err != nil {
		return "", err
	}
	// set query type
	switch q.action {
	case ActionSelect:
		if err := binary.Write(q.buf, binary.LittleEndian, []byte("SELECT")); err != nil {
			return "", err
		}
	case ActionInsert:
		if err := binary.Write(q.buf, binary.LittleEndian, []byte("INSERT INTO")); err != nil {
			return "", err
		}
		// add table name
		if err := binary.Write(q.buf, binary.LittleEndian, []byte(" "+q.tableName)); err != nil {
			return "", err
		}
	case ActionUpdate:
		if err := binary.Write(q.buf, binary.LittleEndian, []byte("UPDATE")); err != nil {
			return "", err
		}
		// add table name
		if err := binary.Write(q.buf, binary.LittleEndian, []byte(" "+q.tableName)); err != nil {
			return "", err
		}
	case ActionDelete:
		if err := binary.Write(q.buf, binary.LittleEndian, []byte("DELETE FROM")); err != nil {
			return "", err
		}
		// add table name
		if err := binary.Write(q.buf, binary.LittleEndian, []byte(" "+q.tableName)); err != nil {
			return "", err
		}
	}

	// add fields by action
	switch q.action {
	case ActionSelect:
		if err := binary.Write(q.buf, binary.LittleEndian,
			[]byte(" "+strings.Join(q.fields, ", "))); err != nil {
			return "", err
		}
		// add table name
		if err := binary.Write(q.buf, binary.LittleEndian, []byte(" FROM "+q.tableName)); err != nil {
			return "", err
		}
	case ActionInsert:
		if err := binary.Write(q.buf, binary.LittleEndian,
			[]byte(" ("+strings.Join(q.fields, ", ")+") VALUES (")); err != nil {
			return "", err
		}

		var s [][]byte
		for range q.fields {
			s = append(s, q.PositionalArg())
		}
		if err := binary.Write(q.buf, binary.LittleEndian,
			bytes.Join(s, []byte(", "))); err != nil {
			return "", err
		}
		if err := binary.Write(q.buf, binary.LittleEndian,
			[]byte(" )")); err != nil {
			return "", err
		}
	case ActionUpdate:
		if err := binary.Write(q.buf, binary.LittleEndian,
			[]byte(" SET")); err != nil {
			return "", err
		}
		var s [][]byte
		for _, field := range q.fields {
			s = append(s, []byte(fmt.Sprintf(" %s=$%d", field, q.queryPosition)))

		}
		if err := binary.Write(q.buf, binary.LittleEndian, bytes.Join(s, []byte(","))); err != nil {
			return "", err
		}
	}

	//add conditions
	if len(q.conditions) > 0 {
		switch q.action {

		case ActionSelect:
			fallthrough
		case ActionUpdate:
			fallthrough
		case ActionDelete:
			if err := binary.Write(q.buf, binary.LittleEndian,
				[]byte(" WHERE ")); err != nil {
				return "", err
			}
			var s [][]byte
			for k, oprator := range q.conditions {
				s = append(s, q.PositionalFieldArg(k, oprator))
			}
			if err := binary.Write(q.buf, binary.LittleEndian,
				bytes.Join(s, []byte(" AND "))); err != nil {
				return "", err
			}
		}
	}

	if q.offset != 0 {
		if err := binary.Write(q.buf, binary.LittleEndian,
			[]byte(fmt.Sprintf(" OFFSET %s", q.PositionalArg()))); err != nil {
			return "", err
		}
	}
	if q.limit != 0 {
		if err := binary.Write(q.buf, binary.LittleEndian,
			[]byte(fmt.Sprintf(" LIMIT %s", q.PositionalArg()))); err != nil {
			return "", err
		}
	}

	if err := binary.Write(q.buf, binary.LittleEndian,
		[]byte(";")); err != nil {
		return "", err
	}

	return q.buf.String(), nil

}

func (q *querer) getPos() int {
	q.queryPosition++
	return q.queryPosition - 1
}
func (q *querer) PositionalArg() []byte {
	return []byte(fmt.Sprintf("$%d", q.getPos()))
}
func (q *querer) PositionalFieldArg(field string, operator OperatorType) []byte {
	var format string
	switch operator {
	case OprEqual:
		format = "%s=$%d"
	case OprNotEqual:
		format = "%s!=$%d"
	case OprGreater:
		format = "$%d<%s"
	case OprGreaterOrEqual:
		format="%s>=$%d"
	case OprLower:
		format = "%s<$%d"
	case OprLowerOrEqual:
		format = "%s<=$%d"
	case OprInArray:
		format = "%s = ANY($%d)"
	case OprArrayOverlap:
		format = "%s && $%d"
	case OprSubstring:
		format = "%s like '%%'||$%d||'%%'"
	}
	return []byte(fmt.Sprintf(format, field, q.getPos()))
}

