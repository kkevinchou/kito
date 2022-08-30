package loot

import (
	"time"

	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/mechanics/items"
	"github.com/kkevinchou/kito/kito/systems/base"
)

type World interface {
	// GetSingleton() *singleton.Singleton
	// GetEntityByID(id int) entities.Entity
	// GetPlayerEntity() entities.Entity
	// GetPlayer() *player.Player
	QueryEntity(componentFlags int) []entities.Entity
}

type LootSystem struct {
	*base.BaseSystem
	world   World
	modPool items.ModPool
}

func NewLootSystem(world World) *LootSystem {
	return &LootSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
		modPool:    *items.NewModPool(),
	}
}

func (s *LootSystem) Update(delta time.Duration) {
	entities := s.world.QueryEntity(components.ComponentFlagLootDropper)

	for _, entity := range entities {
		ldComponent := entity.GetComponentContainer().LootDropperComponent
		if !ldComponent.Drop {
			continue
		}

		rarity := items.SelectRarity(ldComponent.Rarities, ldComponent.RarityWeights)
		modCount := items.RarityToModCount(rarity)
		maxPrefix, maxSuffix := items.MaxCountsByRarity(rarity)
		s.modPool.ChooseMods(modCount, maxPrefix, maxSuffix)
	}
}

func (s *LootSystem) Name() string {
	return "LootSystem"
}
