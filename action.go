package querer


type Action int

const (
	ActionNONE Action = iota
	ActionSelect
	ActionUpdate
	ActionDelete
	ActionInsert
)


