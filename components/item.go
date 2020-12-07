package components

import "github.com/kkevinchou/kito/interfaces"

type ItemComponent struct {
	owner interfaces.ItemReceiver
}

func (i *ItemComponent) OwnedBy() interfaces.ItemReceiver {
	return i.owner
}

func (i *ItemComponent) Owned() bool {
	return i.owner != nil
}

func (i *ItemComponent) SetOwner(owner interfaces.ItemReceiver) {
	i.owner = owner
}
