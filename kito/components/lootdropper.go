package components

import "github.com/kkevinchou/kito/kito/mechanics/items"

type LootDropperComponent struct {
	Rarities              []items.Rarity
	RarityWeights         []int
	IncreasedItemRarity   float64
	IncreasedItemQuantity float64
}

func (c *LootDropperComponent) AddToComponentContainer(container *ComponentContainer) {
	container.LootDropperComponent = c
}

func (c *LootDropperComponent) ComponentFlag() int {
	return ComponentFlagLootDropper
}
