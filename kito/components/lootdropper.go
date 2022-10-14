package components

import "github.com/kkevinchou/kito/kito/mechanics/items"

type LootDropperComponent struct {
	Rarities      []items.Rarity
	RarityWeights []int
}

func (c *LootDropperComponent) AddToComponentContainer(container *ComponentContainer) {
	container.LootDropperComponent = c
}

func (c *LootDropperComponent) ComponentFlag() int {
	return ComponentFlagLootDropper
}

func DefaultLootDropper() *LootDropperComponent {
	return &LootDropperComponent{
		Rarities:      []items.Rarity{items.RarityRare},
		RarityWeights: []int{1},
	}
}
