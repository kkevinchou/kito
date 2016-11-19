package id

var counter int

type IdComponent struct {
	id int
}

func NewIdComponent() *IdComponent {
	component := IdComponent{id: counter}
	counter += 1
	return &component
}

func (i *IdComponent) ID() int {
	return i.id
}
