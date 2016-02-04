package item

type ItemI interface {
	OwnedBy() int
	Owned() bool
	Id() int
}

type Manager struct {
	items map[int]ItemI
}

func (i *Manager) Register(item ItemI) {
	i.items[item.Id()] = item
}

func NewManager() *Manager {
	return &Manager{
		items: map[int]ItemI{},
	}
}
