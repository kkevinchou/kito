package entityutils

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/knetwork"
	"github.com/kkevinchou/kito/kito/types"
)

func Spawn(entityID int, entityType types.EntityType, position mgl64.Vec3, orientation mgl64.Quat) *entities.EntityImpl {
	var newEntity *entities.EntityImpl

	if types.EntityType(entityType) == types.EntityTypeBob {
		newEntity = entities.NewBob()
	} else if types.EntityType(entityType) == types.EntityTypeScene {
		newEntity = entities.NewScene()
	} else if types.EntityType(entityType) == types.EntityTypeStaticSlime {
		newEntity = entities.NewSlime()
	} else if types.EntityType(entityType) == types.EntityTypeDynamicRigidBody {
		newEntity = entities.NewDynamicRigidBody()
	} else if types.EntityType(entityType) == types.EntityTypeStaticRigidBody {
		newEntity = entities.NewStaticRigidBody()
	} else if types.EntityType(entityType) == types.EntityTypeProjectile {
		newEntity = entities.NewProjectile(position)
	} else if types.EntityType(entityType) == types.EntityTypeEnemy {
		newEntity = entities.NewEnemy()
	} else {
		fmt.Printf("unrecognized entity with type %v to spawn\n", entityType)
		return nil
	}

	newEntity.ID = entityID
	cc := newEntity.GetComponentContainer()
	cc.TransformComponent.Position = position
	cc.TransformComponent.Orientation = orientation
	return newEntity
}

func ConstructEntitySnapshot(entity entities.Entity) knetwork.EntitySnapshot {
	cc := entity.GetComponentContainer()
	transformComponent := cc.TransformComponent
	tpcComponent := cc.ThirdPersonControllerComponent

	snapshot := knetwork.EntitySnapshot{
		ID:          entity.GetID(),
		Type:        int(entity.Type()),
		Position:    transformComponent.Position,
		Orientation: transformComponent.Orientation,
	}

	physicsComponent := cc.PhysicsComponent
	if physicsComponent != nil {
		snapshot.Velocity = physicsComponent.Velocity
		snapshot.Impulses = physicsComponent.Impulses
	} else if tpcComponent != nil {
		snapshot.Velocity = tpcComponent.Velocity
	}

	animationComponent := cc.AnimationComponent
	if animationComponent != nil {
		snapshot.Animation = animationComponent.Player.CurrentAnimation()
	}

	return snapshot
}
