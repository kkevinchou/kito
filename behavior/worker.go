package behavior

type WorkerI interface {
	AddItem(interface{})
	DropItem(interface{})
	HasItem(interface{})
}
