package spatialpartition

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/lib/collision/collider"
)

type World interface {
	// GetSingleton() *singleton.Singleton
	// GetPlayerEntity() entities.Entity
	QueryEntity(componentFlags int) []entities.Entity
	// GetPlayer() *player.Player
	// GetEntityByID(id int) entities.Entity
}

type Partition struct {
	AABB     *collider.BoundingBox
	entities []entities.Entity
}

type SpatialPartition struct {
	Partitions         [][][]Partition
	PartitionDimension int
	PartitionCount     int
}

// NewSpatialPartition creates a spatial partition with the bottom at <0, 0, 0>
// the spatial partition spans the rectangular space for
// d = partitionDimension * partitionCount
// <-d, 0, -d> to <d, 2 * d, d>
func NewSpatialPartition(partitionDimension int, partitionCount int) *SpatialPartition {
	d := partitionDimension * partitionCount
	// partitions := make([][]Partition, 4*partitionDimension)
	partitions := make([][][]Partition, partitionCount)
	for i := 0; i < partitionCount; i++ {
		partitions[i] = make([][]Partition, partitionCount)
		for j := 0; j < partitionCount; j++ {
			partitions[i][j] = make([]Partition, partitionCount)
			for k := 0; k < partitionCount; k++ {
				partitions[i][j][k] = Partition{
					AABB: &collider.BoundingBox{
						MinVertex: mgl64.Vec3{float64(i*partitionDimension - d/2), float64(j*partitionDimension - d/2), float64(k*partitionDimension - d/2)},
						MaxVertex: mgl64.Vec3{float64((i+1)*partitionDimension - d/2), float64((j+1)*partitionDimension - d/2), float64((k+1)*partitionDimension - d/2)},
					},
				}
			}
		}
	}
	return &SpatialPartition{
		Partitions:         partitions,
		PartitionDimension: partitionDimension,
		PartitionCount:     partitionCount,
	}
}

// QueryCollisionCandidates queries for collision candidates that have been stored in
// the spatial partition
func (s *SpatialPartition) QueryCollisionCandidates(entity entities.Entity) []entities.Entity {
	return nil
}

func (s *SpatialPartition) Initialize(world World) {
	entityList := world.QueryEntity(components.ComponentFlagCollider | components.ComponentFlagTransform)
	_ = entityList

	// fmt.Println("----------")
	// fmt.Println(s.Partitions)
	// fmt.Println("----------")
	// populate partitions with all of the entities
}
