package stail

type STail interface {
	Tail(filepath string, tailLine int, call func(content string)) (err error)
	Total(filepath string, call func(content string)) (err error)
}
