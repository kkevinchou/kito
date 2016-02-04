package food

type ItemComponent struct {
}

func (f *ItemComponent) OwnedBy() int {
	return 0
}

func (f *ItemComponent) Owned() bool {
	return true
}
