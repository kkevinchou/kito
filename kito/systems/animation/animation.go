package animation

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/lib/libutils"
	"github.com/kkevinchou/kito/lib/modelspec"
)

type World any

type AnimationSystem struct {
	*base.BaseSystem
	world    World
	entities []entities.Entity
}

func NewAnimationSystem(world World) *AnimationSystem {
	return &AnimationSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
		entities:   []entities.Entity{},
	}
}

func (s *AnimationSystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	if componentContainer.AnimationComponent != nil {
		s.entities = append(s.entities, entity)
	}
}

func (s *AnimationSystem) Update(delta time.Duration) {
	for _, entity := range s.entities {
		componentContainer := entity.GetComponentContainer()
		animationComponent := componentContainer.AnimationComponent
		tpcComponent := componentContainer.ThirdPersonControllerComponent

		targetAnimation := "Idle"
		if !libutils.Vec3IsZero(tpcComponent.Velocity) {
			targetAnimation = "Walk"
		}

		if animationComponent.CurrentAnimation != targetAnimation {
			animationComponent.CurrentAnimation = targetAnimation
			animationComponent.ElapsedTime = 0
			animationComponent.Animation = animationComponent.Animations[targetAnimation]
		}

		animationComponent.ElapsedTime += delta
		for animationComponent.ElapsedTime.Milliseconds() > animationComponent.Animation.Length().Milliseconds() {
			animationComponent.ElapsedTime = time.Duration(animationComponent.ElapsedTime.Milliseconds()-animationComponent.Animation.Length().Milliseconds()) * time.Millisecond
		}

		pose := calculateCurrentAnimationPose(animationComponent.ElapsedTime, animationComponent.Animation.KeyFrames())
		animationTransforms := computeJointTransforms(animationComponent.Animation.RootJoint(), pose)
		animationComponent.AnimationTransforms = animationTransforms
	}
}

// applyPoseToJoints returns the set of transforms that move the joint from the bind pose to the given pose
func computeJointTransforms(joint *modelspec.JointSpec, pose map[int]mgl32.Mat4) map[int]mgl32.Mat4 {
	animationTransforms := map[int]mgl32.Mat4{}
	computeJointTransformsHelper(joint, mgl32.Ident4(), pose, animationTransforms)
	return animationTransforms
}

func computeJointTransformsHelper(joint *modelspec.JointSpec, parentTransform mgl32.Mat4, pose map[int]mgl32.Mat4, transforms map[int]mgl32.Mat4) {
	localTransform := pose[joint.ID]

	if _, ok := pose[joint.ID]; !ok {
		// panic(fmt.Sprintf("joint with id %d does not have a pose", joint.ID))
		localTransform = joint.BindTransform
	}

	// model-space transform that includes all the parental transforms
	// and the local transform, not meant to be used to transform any vertices
	// until we multiply it by the inverse bind transform
	poseTransform := parentTransform.Mul4(localTransform)

	for _, child := range joint.Children {
		computeJointTransformsHelper(child, poseTransform, pose, transforms)
	}

	// this is the model-space transform that can finally be used to transform
	// any vertices it influences
	transforms[joint.ID] = poseTransform.Mul4(joint.InverseBindTransform)
}

func calculateCurrentAnimationPose(elapsedTime time.Duration, keyFrames []*modelspec.KeyFrame) map[int]mgl32.Mat4 {
	var startKeyFrame *modelspec.KeyFrame
	var endKeyFrame *modelspec.KeyFrame
	var progression float32

	// iterate backwards looking for the starting keyframe
	for i := len(keyFrames) - 1; i >= 0; i-- {
		keyFrame := keyFrames[i]
		if elapsedTime >= keyFrame.Start || i == 0 {
			startKeyFrame = keyFrame
			if i < len(keyFrames)-1 {
				endKeyFrame = keyFrames[i+1]
				progression = float32((elapsedTime - startKeyFrame.Start).Milliseconds()) / float32((endKeyFrame.Start - startKeyFrame.Start).Milliseconds())
			} else {
				// interpolate towards the first kf, assume looping animations
				endKeyFrame = keyFrames[0]
				progression = 0
			}
			break
		}
	}

	// progression = 0
	// startKeyFrame = keyFrames[0]
	return interpolatePoses(startKeyFrame, endKeyFrame, progression)
}

func interpolatePoses(k1, k2 *modelspec.KeyFrame, progression float32) map[int]mgl32.Mat4 {
	interpolatedPose := map[int]mgl32.Mat4{}
	for jointID := range k1.Pose {
		k1JointTransform := k1.Pose[jointID]
		k2JointTransform := k2.Pose[jointID]

		// WTF - this lerp doesn't look right when interpolating keyframes???
		// rotationQuat := mgl32.QuatLerp(k1JointTransform.Rotation, k2JointTransform.Rotation, progression)

		rotationQuat := libutils.QInterpolate(k1JointTransform.Rotation, k2JointTransform.Rotation, progression)
		rotation := rotationQuat.Mat4()

		translation := k1JointTransform.Translation.Add(k2JointTransform.Translation.Sub(k1JointTransform.Translation).Mul(progression))
		scale := k1JointTransform.Scale.Add(k2JointTransform.Scale.Sub(k1JointTransform.Scale).Mul(progression))
		// scale = mgl32.Vec3{0.5, 0.5, 0.5}

		interpolatedPose[jointID] = mgl32.Translate3D(translation.X(), translation.Y(), translation.Z()).Mul4(rotation).Mul4(mgl32.Scale3D(scale.X(), scale.Y(), scale.Z()))
	}
	// pause hip joint
	// interpolatedPose[42] = mgl32.Ident4()
	return interpolatedPose
}
