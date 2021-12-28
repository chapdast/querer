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
			s = append(s, []byte(fmt.Sprintf("$%d", q.queryPosition)))
			q.queryPosition += 1
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
			q.queryPosition += 1
		}
		if err := binary.Write(q.buf, binary.LittleEndian, bytes.Join(s, []byte(","))); err != nil {
			return "", err
		}
	}

	//add conditions
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
			switch oprator {
			case OprEqual:
				s = append(s, []byte(fmt.Sprintf("%s=$%d", k, q.queryPosition)))
			case OprNotEqual:
				s = append(s, []byte(fmt.Sprintf("%s!=$%d", k, q.queryPosition)))
			case OprGreater:
				s = append(s, []byte(fmt.Sprintf("$%d<%s", q.queryPosition, k)))
			case OprGreaterOrEqual:
				s = append(s, []byte(fmt.Sprintf("$%d<=%s", q.queryPosition, k)))
			case OprLower:
				s = append(s, []byte(fmt.Sprintf("%s<$%d", k, q.queryPosition)))
			case OprLowerOrEqual:
				s = append(s, []byte(fmt.Sprintf("%s<=$%d", k, q.queryPosition)))
			case OprInArray:
				s = append(s, []byte(fmt.Sprintf("%s = ANY($%d)", k, q.queryPosition)))
			case OprArrayOverlap:
				s = append(s, []byte(fmt.Sprintf("%s && $%d", k, q.queryPosition)))
			case OprSubstring:
				s = append(s, []byte(fmt.Sprintf("%s like '%%'||$%d||'%%'", k, q.queryPosition)))
			}
			q.queryPosition += 1
		}
		if err := binary.Write(q.buf, binary.LittleEndian,
			bytes.Join(s, []byte(" AND "))); err != nil {
			return "", err
		}
	}

	if err := binary.Write(q.buf, binary.LittleEndian,
		[]byte(";")); err != nil {
		return "", err
	}

	return q.buf.String(), nil

}
