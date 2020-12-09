package components

import "github.com/kkevinchou/kito/types"

type ItemComponent struct {
	owner types.ItemReceiver
}

func (i *ItemComponent) OwnedBy() types.ItemReceiver {
	return i.owner
}

func (i *ItemComponent) Owned() bool {
	return i.owner != nil
}

func (i *ItemComponent) SetOwner(owner types.ItemReceiver) {
	i.owner = owner
}
