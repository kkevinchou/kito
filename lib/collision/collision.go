package collision

import "github.com/kkevinchou/kito/lib/collision/collider"

type ContactManifold struct {
}

func CheckCollision(capsule collider.Capsule, boundingBox collider.BoundingBox) ContactManifold {

	return ContactManifold{}
}
