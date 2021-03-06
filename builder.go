package querer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type Builder interface {
	Build(opts ...Option) (string, []interface{}, error)
}

func (q *querer) Build(opts ...Option) (string, []interface{}, error) {

	defer q.reset()
	if err := q.loadOpts(opts); err != nil {
		return "", nil, err
	}
	// set query type
	switch q.action {
	case ActionSelect:
		if err := binary.Write(q.buf, binary.LittleEndian, []byte("SELECT")); err != nil {
			return "", nil, err
		}
	case ActionInsert:
		if err := binary.Write(q.buf, binary.LittleEndian, []byte("INSERT INTO")); err != nil {
			return "", nil, err
		}
		// add table name
		if err := binary.Write(q.buf, binary.LittleEndian, []byte(" "+q.tableName)); err != nil {
			return "", nil, err
		}
	case ActionUpdate:
		if err := binary.Write(q.buf, binary.LittleEndian, []byte("UPDATE")); err != nil {
			return "", nil, err
		}
		// add table name
		if err := binary.Write(q.buf, binary.LittleEndian, []byte(" "+q.tableName)); err != nil {
			return "", nil, err
		}
	case ActionDelete:
		if err := binary.Write(q.buf, binary.LittleEndian, []byte("DELETE FROM")); err != nil {
			return "", nil, err
		}
		// add table name
		if err := binary.Write(q.buf, binary.LittleEndian, []byte(" "+q.tableName)); err != nil {
			return "", nil, err
		}
	}

	// add fields by action
	switch q.action {
	case ActionSelect:
		if err := binary.Write(q.buf, binary.LittleEndian,
			[]byte(" "+strings.Join(q.fields, ", "))); err != nil {
			return "", nil, err
		}
		// add table name
		if err := binary.Write(q.buf, binary.LittleEndian, []byte(" FROM "+q.tableName)); err != nil {
			return "", nil, err
		}
	case ActionInsert:
		if err := binary.Write(q.buf, binary.LittleEndian,
			[]byte(" ("+strings.Join(q.fields, ", ")+") VALUES (")); err != nil {
			return "", nil, err
		}

		var s [][]byte
		for range q.fields {
			s = append(s, q.PositionalArg())
		}
		if err := binary.Write(q.buf, binary.LittleEndian,
			bytes.Join(s, []byte(", "))); err != nil {
			return "", nil, err
		}
		if err := binary.Write(q.buf, binary.LittleEndian,
			[]byte(" )")); err != nil {
			return "", nil, err
		}
	case ActionUpdate:
		if err := binary.Write(q.buf, binary.LittleEndian,
			[]byte(" SET ")); err != nil {
			return "", nil, err
		}
		var s [][]byte
		for _, field := range q.fields {
			s = append(s, q.PositionalFieldArg(field, OprEqual))
		}
		if err := binary.Write(q.buf, binary.LittleEndian, bytes.Join(s, []byte(", "))); err != nil {
			return "", nil, err
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
				return "", nil, err
			}
			var s [][]byte
			var join []Joiner
			for k, condition := range q.conditions {
				for _, cond := range condition {
					s = append(s, q.PositionalFieldArg(k, cond.Operator()))
					join = append(join, cond.ANDOR)
				}

			}
			var conditions [][]byte
			if len(s)>1{
				for i, c:= range s {
					conditions = append(conditions, c)
					if i != len(join)-1 {
						if join[i] == AND {
							conditions = append(conditions, []byte("AND"))
						} else if join[i] == OR {
							conditions = append(conditions, []byte("OR"))
						}
					}
				}
			}else{
				conditions = s
			}

 			if err := binary.Write(q.buf, binary.LittleEndian,
				bytes.Join(conditions, []byte(" "))); err != nil {
				return "", nil, err
			}
		}
	}

	if q.offset != 0 {
		if err := binary.Write(q.buf, binary.LittleEndian,
			[]byte(fmt.Sprintf(" OFFSET %s", q.PositionalArg()))); err != nil {
			return "", nil, err
		}
		q.data = append(q.data, q.offset)
	}
	if q.limit != 0 {
		if err := binary.Write(q.buf, binary.LittleEndian,
			[]byte(fmt.Sprintf(" LIMIT %s", q.PositionalArg()))); err != nil {
			return "", nil, err
		}
		q.data = append(q.data, q.limit)
	}

	if err := binary.Write(q.buf, binary.LittleEndian,
		[]byte(";")); err != nil {
		return "", nil, err
	}

	return q.buf.String(), q.data, nil

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
		format = "%s>$%d"
	case OprGreaterOrEqual:
		format = "%s>=$%d"
	case OprLower:
		format = "%s<$%d"
	case OprLowerOrEqual:
		format = "%s<=$%d"
	case OprInArray:
		format = "%s = ANY($%d)"
	case OprNotInArray:
		format = "%s != ANY($%d)"
	case OprArrayOverlap:
		format = "%s && $%d"
	case OprNotArrayOverlap:
		format = " NOT %s && $%d"
	case OprSubstring:
		format = "%s LIKE '%%'||$%d||'%%'"
	case OprNotSubstring:
		format = "%s NOT LIKE '%%'||$%d||'%%'"
	}
	return []byte(fmt.Sprintf(format, field, q.getPos()))
}
