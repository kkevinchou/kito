package animation

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/lib/model"
	"github.com/kkevinchou/kito/systems/base"
)

type World interface{}

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

		animationComponent.ElapsedTime += delta
		for animationComponent.ElapsedTime.Milliseconds() > animationComponent.Animation.Length.Milliseconds() {
			animationComponent.ElapsedTime = time.Duration(animationComponent.ElapsedTime.Milliseconds()-animationComponent.Animation.Length.Milliseconds()) * time.Millisecond
		}

		// keyFrames := animationComponent.Animation.KeyFrames
		// // set translation to 0 for all keyframes for the spine2
		// for z := 0; z < len(keyFrames); z++ {
		// 	keyFrames[z].Pose[38].Translation = mgl32.Vec3{0, 0, 0}
		// 	keyFrames[z].Pose[38].Rotation = mgl32.QuatIdent()
		// }

		pose := calculateCurrentAnimationPose(animationComponent.ElapsedTime, animationComponent.Animation.KeyFrames)

		// jointID := 0
		// fmt.Println("---------------")
		// fmt.Println("POSE FOR JOINT\n", pose[jointID])
		animationTransforms := applyPoseToJoints(animationComponent.Animation.RootJoint, pose, map[int]bool{1000: true})
		// fmt.Println("ANIMATION TRANSFORM FOR JOINT\n", animationTransforms[jointID])

		// 38, 50, 51

		animationComponent.AnimationTransforms = animationTransforms

		// kf1Pose := calculateCurrentAnimationPose(0, keyFrames)
		// fmt.Println("---------------------------------")
		// fmt.Println(keyFrames[0].Pose[0].Rotation)
		// fmt.Println(keyFrames[0].Pose[0].Translation)
		// spine2 := kf1Pose[0]
		// kf1AnimationTransforms := applyPoseToJoints(animationComponent.Animation.RootJoint, kf1Pose)
		// for i, _ := range animationComponent.AnimationTransforms {
		// 	if i != 0 {
		// 		animationComponent.AnimationTransforms[i] = kf1AnimationTransforms[i]
		// 	}
		// }
	}
}

// applyPoseToJoints returns the set of transforms that move the joint from the bind pose to the given pose
func applyPoseToJoints(joint *model.Joint, pose map[int]mgl32.Mat4, frozenJoints map[int]bool) map[int]mgl32.Mat4 {
	animationTransforms := map[int]mgl32.Mat4{}
	applyPoseToJointsHelper(joint, mgl32.Ident4(), pose, animationTransforms, frozenJoints)
	return animationTransforms
}

func applyPoseToJointsHelper(joint *model.Joint, parentTransform mgl32.Mat4, pose map[int]mgl32.Mat4, transforms map[int]mgl32.Mat4, frozenJoints map[int]bool) {
	localTransform := pose[joint.ID]

	// TODO: is this what i actually want to do? probably need to refactor this
	// when we start keyframing joints at different frames
	if _, ok := pose[joint.ID]; !ok {
		// fmt.Println("HI")
		localTransform = joint.LocalBindTransform
	}

	poseTransform := parentTransform.Mul4(localTransform) // model-space relative to the origin

	if joint.ID == 15 {
		// fmt.Println("POSE LOCAL\n", pose[joint.ID])
		// fmt.Println("POSE TRANSFORM\n", parentTransform.Mul4((joint.LocalBindTransform)))
		// fmt.Println("FINAL TRANSFORM\n", poseTransform.Mul4(joint.InverseBindTransform))
	}

	if _, ok := frozenJoints[joint.ID]; ok {
		// poseTransform = parentTransform.Mul4((joint.LocalBindTransform))

		// fmt.Println("---------------------------")
		// fmt.Println("SELF\n", joint.LocalBindTransform)
		// _, r := libutils.Decompose(joint.LocalBindTransform)
		// fmt.Println("SELF QUAT\n", r)
		// fmt.Println("PARENT\n", parentTransform)

		// fmt.Println("OPPOSITE ========================")
		// fmt.Println("SELF\n", joint.LocalBindTransform)
		// _, r := libutils.Decompose(joint.LocalBindTransform)
		// fmt.Println("SELF QUAT\n", r)
		// fmt.Println("POSE LEFT LOCAL\n", pose[15])
	}
	for _, child := range joint.Children {
		applyPoseToJointsHelper(child, poseTransform, pose, transforms, frozenJoints)
	}
	// if joint.ID == 0 {
	// 	fmt.Println("LOCAL BIND TRANSFORM\n", joint.LocalBindTransform)
	// 	fmt.Println("POSE TRANSFORM\n", poseTransform)
	// 	fmt.Println("INVERSE BIND TRANSFORM\n", joint.InverseBindTransform)
	// 	fmt.Println("OUTPUT TRANSFORM\n", poseTransform.Mul4(joint.InverseBindTransform))
	// }
	transforms[joint.ID] = poseTransform.Mul4(joint.InverseBindTransform) // model-space relative to the bind pose
}

func calculateCurrentAnimationPose(elapsedTime time.Duration, keyFrames []*model.KeyFrame) map[int]mgl32.Mat4 {
	var startKeyFrame *model.KeyFrame
	var endKeyFrame *model.KeyFrame
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

	return interpolatePoses(startKeyFrame, endKeyFrame, progression)
	// _ = progression
	// return interpolatePoses(keyFrames[0], keyFrames[1], 0)
}

func interpolatePoses(k1, k2 *model.KeyFrame, progression float32) map[int]mgl32.Mat4 {
	interpolatedPose := map[int]mgl32.Mat4{}
	for jointID := range k1.Pose {
		k1JointTransform := k1.Pose[jointID]
		k2JointTransform := k2.Pose[jointID]

		// WTF - this lerp doesn't look right when interpolating keyframes???
		// rotationQuat := mgl32.QuatLerp(k1JointTransform.Rotation, k2JointTransform.Rotation, progression)

		rotationQuat := qInterpolate(k1JointTransform.Rotation, k2JointTransform.Rotation, progression)
		rotation := rotationQuat.Mat4()

		translation := k1JointTransform.Translation.Add(k2JointTransform.Translation.Sub(k1JointTransform.Translation).Mul(progression))

		interpolatedPose[jointID] = mgl32.Translate3D(translation.X(), translation.Y(), translation.Z()).Mul4(rotation)
	}
	return interpolatedPose
}

// Quaternion interpolation, reimplemented from: https://github.com/TheThinMatrix/OpenGL-Animation/blob/dde792fe29767192bcb60d30ac3e82d6bcff1110/Animation/animation/Quaternion.java#L158
func qInterpolate(a, b mgl32.Quat, blend float32) mgl32.Quat {
	var result mgl32.Quat = mgl32.Quat{}
	var dot float32 = a.W*b.W + a.V.X()*b.V.X() + a.V.Y()*b.V.Y() + a.V.Z()*b.V.Z()
	blendI := float32(1) - blend
	if dot < 0 {
		result.W = blendI*a.W + blend*-b.W
		result.V = mgl32.Vec3{
			blendI*a.V.X() + blend*-b.V.X(),
			blendI*a.V.Y() + blend*-b.V.Y(),
			blendI*a.V.Z() + blend*-b.V.Z(),
		}
	} else {
		result.W = blendI*a.W + blend*b.W
		result.V = mgl32.Vec3{
			blendI*a.V.X() + blend*b.V.X(),
			blendI*a.V.Y() + blend*b.V.Y(),
			blendI*a.V.Z() + blend*b.V.Z(),
		}
	}
	result.Normalize()
	return result
}
