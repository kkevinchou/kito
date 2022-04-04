package model

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/modelspec"
)

// JointTransform represents the joint-space transformations that should be
// applied to the joint for the KeyFrame it is associated with.
type JointTransform struct {
	Translation mgl32.Vec3
	Scale       mgl32.Vec3
	Rotation    mgl32.Quat
}

type Animation struct {
	name          string
	rootJoint     *modelspec.JointSpec
	animationSpec *modelspec.AnimationSpec
}

func (a *Animation) Name() string {
	return a.name
}

func (a *Animation) RootJoint() *modelspec.JointSpec {
	return a.rootJoint
}

func (a *Animation) KeyFrames() []*modelspec.KeyFrame {
	return a.animationSpec.KeyFrames
}

func (a *Animation) Length() time.Duration {
	return a.animationSpec.Length
}

func NewAnimations(spec *modelspec.ModelSpecification) map[string]*Animation {
	animations := map[string]*Animation{}
	for name, animation := range spec.Animations {
		animations[name] = &Animation{
			name:          name,
			animationSpec: animation,
			rootJoint:     spec.RootJoint,
		}
	}

	return animations
}
