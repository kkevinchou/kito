package base

type BaseSystem struct {
}

func (b *BaseSystem) UpdateOnCommandFrame() bool {
	return true
}
