package loot

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/events"
	"github.com/kkevinchou/kito/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/kito/mechanics/items"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/kito/utils/entityutils"
)

type World interface {
	QueryEntity(componentFlags int) []entities.Entity
	RegisterEntities([]entities.Entity)
	GetEntityByID(id int) entities.Entity
	CommandFrame() int
	GetEventBroker() eventbroker.EventBroker
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

	// drop loot from entities who have reached 0 health
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

	// add loot to a player's inventory
	inventoryEntities := s.world.QueryEntity(components.ComponentFlagInventory)
	for _, entity := range inventoryEntities {
		cc := entity.GetComponentContainer()
		if cc.ColliderComponent == nil || len(cc.ColliderComponent.Contacts) == 0 {
			continue
		}

		for e2ID := range cc.ColliderComponent.Contacts {
			cEntity := s.world.GetEntityByID(e2ID)
			if cEntity.Type() == types.EntityTypeLootbox {
				cc.InventoryComponent.Items = append(cc.InventoryComponent.Items, items.Item{ID: 69, Type: items.ItemTypeCoin})
				event := &events.UnregisterEntityEvent{
					GlobalCommandFrame: s.world.CommandFrame(),
					EntityID:           cEntity.GetID(),
				}
				s.world.GetEventBroker().Broadcast(event)
			}
		}
	}
}

func (s *LootSystem) Name() string {
	return "LootSystem"
}
