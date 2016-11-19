package food

type ItemComponent struct {
	owned   bool
	ownedBy int
}

func (i *ItemComponent) OwnedBy() int {
	return i.ownedBy
}

func (i *ItemComponent) Owned() bool {
	return i.owned
}
