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
	OprNotInArray
	OprArrayOverlap
	OprNotArrayOverlap
	OprSubstring
	OprNotSubstring
)

func (o OperatorType) Flip() OperatorType {
	switch o {
	case OprEqual:
		return OprNotEqual
	case OprNotEqual:
		return OprEqual
	case OprGreater:
		return OprLowerOrEqual
	case OprGreaterOrEqual:
		return OprLower
	case OprLower:
		return OprGreaterOrEqual
	case OprLowerOrEqual:
		return OprGreater
	case OprInArray:
		return OprNotInArray
	case OprNotInArray:
		return OprInArray
	case OprArrayOverlap:
		return OprNotArrayOverlap
	case OprNotArrayOverlap:
		return OprArrayOverlap
	case OprSubstring:
		return OprNotSubstring
	case OprNotSubstring:
		return OprSubstring
	}
	return o
}

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
