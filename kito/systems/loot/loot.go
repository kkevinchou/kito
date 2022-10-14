package loot

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/mechanics/items"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/kito/utils/entityutils"
)

type World interface {
	QueryEntity(componentFlags int) []entities.Entity
	RegisterEntities([]entities.Entity)
}

type LootSystem struct {
	*base.BaseSystem
	world   World
	modPool *items.ModPool
}

func NewLootSystem(world World) *LootSystem {
	modPool := items.NewModPool()

	for i := 0; i < 100; i++ {
		modPool.AddMod(&items.Mod{ID: i, AffixType: items.AffixTypePrefix})
	}
	for i := 1000; i < 1100; i++ {
		modPool.AddMod(&items.Mod{ID: i, AffixType: items.AffixTypeSuffix})
	}

	return &LootSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
		modPool:    modPool,
	}
}

func (s *LootSystem) Update(delta time.Duration) {
	lootEntities := s.world.QueryEntity(components.ComponentFlagLootDropper)

	for _, entity := range lootEntities {
		cc := entity.GetComponentContainer()
		ldComponent := cc.LootDropperComponent
		healthComponent := cc.HealthComponent
		if ldComponent == nil || healthComponent == nil {
			continue
		}

		if healthComponent.Value > 0 {
			continue
		}

		rarity := items.SelectRarity(ldComponent.Rarities, ldComponent.RarityWeights)
		mods := s.modPool.ChooseMods(rarity)
		_ = mods

		lootbox := entityutils.Spawn(types.EntityTypeLootbox, cc.TransformComponent.Position.Add(mgl64.Vec3{0, 25, 0}), cc.TransformComponent.Orientation)
		s.world.RegisterEntities([]entities.Entity{lootbox})
	}
}

func (s *LootSystem) Name() string {
	return "LootSystem"
}
