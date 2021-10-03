package types

type EntityType int

const (
	EntityTypeBob EntityType = iota
	EntityTypeBlock
	EntityTypeRigidBody
	EntityTypeCamera
	EntityTypeScene
	EntityTypeStaticSlime
)
