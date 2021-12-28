package querer

type OperatorType int

const (
	OprEqual OperatorType = iota
	OprNotEqual
	OprGreater
	OprGreaterOrEqual
	OprLower
	OprLowerOrEqual
	OprInArray
	OprArrayOverlap
	OprSubstring
)

type Condition interface {
	Type() OperatorType
}

type condition struct {
	condType  OperatorType
	condValue []byte
}

func (c condition) Type() OperatorType {
	return c.condType
}
