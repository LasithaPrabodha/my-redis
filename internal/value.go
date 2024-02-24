package internal

type Value struct {
	Typ   string
	Str   string
	Num   int
	Bulk  string
	Array []Value
}
