package entitymanager

import (
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/entities"
)

type EntityManager struct {
	entities  map[int]entities.Entity
	entityIDs []int
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		entities: map[int]entities.Entity{},
	}

}

func (em *EntityManager) AddEntity(e entities.Entity) {
	em.entities[e.GetID()] = e
	em.entityIDs = append(em.entityIDs, e.GetID())

}

func (em *EntityManager) GetEntityByID(id int) entities.Entity {
	return em.entities[id]
}

// TODO: cache queries
func (em *EntityManager) Query(cs ...components.Component) []entities.Entity {
	var bitFlags int

	for _, c := range cs {
		bitFlags |= c.ComponentFlag()
	}

	var matches []entities.Entity
	for _, id := range em.entityIDs {
		e := em.entities[id]
		cc := e.GetComponentContainer()
		if cc.MatchBitFlags(bitFlags) {
			matches = append(matches, e)
		}
	}

	return matches
}
