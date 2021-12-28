package querer


type OrderBy interface {
	Field() []byte
	Descending() bool
}



type orderBy struct {
	field string
	desc  bool
}

func (o orderBy) Field() []byte {
	return []byte(o.field)
}
func (o orderBy) Descending() bool {
	return o.desc
}
