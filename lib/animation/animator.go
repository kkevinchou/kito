package animation

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type Animator struct {
	Animation     *Animation
	AnimatedModel *AnimatedModel
	ElapsedTime   time.Duration
}

func NewAnimator(animatedModel *AnimatedModel, animation *Animation) *Animator {
	animator := &Animator{
		Animation:     animation,
		AnimatedModel: animatedModel,
	}
	animator.Init()
	return animator
}

func (a *Animator) Init() {
	a.AnimatedModel.RootJoint.CalculateInverseBindTransform(mgl32.Ident4())
}

func (a *Animator) Update(delta time.Duration) {
	a.ElapsedTime += delta
	if a.ElapsedTime.Milliseconds() > a.Animation.Length.Milliseconds() {
		a.ElapsedTime = 0
	}
	pose := a.calculateCurrentAnimationPose()
	a.ApplyPoseToJoints(a.AnimatedModel.RootJoint, mgl32.Ident4(), pose)
}

func (a *Animator) ApplyPoseToJoints(joint *Joint, parentTransform mgl32.Mat4, pose map[int]mgl32.Mat4) {
	localTransform := pose[joint.ID]
	poseTransform := parentTransform.Mul4(localTransform) // model-space relative to the origin
	for _, child := range joint.Children {
		a.ApplyPoseToJoints(child, poseTransform, pose)
	}
	joint.AnimationTransform = poseTransform.Mul4(joint.InverseBindTransform) // model-space relative to the bind pose
}

// CollectAnimationTransforms recursively collects all of the animation transforms
// which are used for transforming joints from their bind pose to the animation position

// this should technically be a dictionary by jointID
func (a *Animator) CollectAnimationTransforms() map[int]mgl32.Mat4 {
	transforms := map[int]mgl32.Mat4{}
	a.collectAnimationTransforms(a.AnimatedModel.RootJoint, transforms)
	return transforms
}

func (a *Animator) collectAnimationTransforms(joint *Joint, transforms map[int]mgl32.Mat4) {
	transforms[joint.ID] = joint.AnimationTransform
	for _, child := range joint.Children {
		a.collectAnimationTransforms(child, transforms)
	}
}

func (a *Animator) CollectBindPoseAnimationTransforms() map[int]mgl32.Mat4 {
	transforms := map[int]mgl32.Mat4{}
	a.collectBindPoseAnimationTransforms(a.AnimatedModel.RootJoint, transforms)
	return transforms
}

func (a *Animator) collectBindPoseAnimationTransforms(joint *Joint, transforms map[int]mgl32.Mat4) {
	transforms[joint.ID] = joint.LocalBindTransform
	for _, child := range joint.Children {
		a.collectBindPoseAnimationTransforms(child, transforms)
	}
}

func (a *Animator) PlayAnimation(animation *Animation) {
	a.Animation = animation
	a.ElapsedTime = 0
}

func (a *Animator) calculateCurrentAnimationPose() map[int]mgl32.Mat4 {
	// need some logic here for wrapping around, do i know how long
	// a key frame is for?
	// lastKeyFrame := a.Animation.KeyFrames[len(a.Animation.KeyFrames)-1]
	// if a.ElapsedTime > lastKeyFrame.

	// var startKeyFrame *KeyFrame
	// var endKeyFrame *KeyFrame
	// for i := 0; i < len(a.Animation.KeyFrames)-1; i++ {
	// 	keyFrame := a.Animation.KeyFrames[i]
	// 	nextKeyFrame := a.Animation.KeyFrames[i+1]
	// 	if a.ElapsedTime >= keyFrame.Start && a.ElapsedTime < nextKeyFrame.Start {
	// 		startKeyFrame = keyFrame
	// 		endKeyFrame = nextKeyFrame
	// 		break
	// 	}
	// }

	// if a.ElapsedTime > 400*time.Millisecond {
	// 	return InterpolatePoses(a.Animation.KeyFrames[2], a.Animation.KeyFrames[3], 0)
	// }

	// return InterpolatePoses(startKeyFrame, endKeyFrame, 0)
	return InterpolatePoses(a.Animation.KeyFrames[0], a.Animation.KeyFrames[1], 0)
}

// TODO fill in actual interpolation.
// it's kinda confusing that this constructs joint transforms, should probably be some other
// data type since normally JointTransforms are immutable
func InterpolateJointTransform(t1 *JointTransform, t2 *JointTransform) *JointTransform {
	return &JointTransform{Translation: t1.Translation, Rotation: t1.Rotation}
}

func CalculateJointTransformMatrix(t *JointTransform) mgl32.Mat4 {
	translationMatrix := mgl32.Translate3D(t.Translation.X(), t.Translation.Y(), t.Translation.Z())
	transformMatrix := translationMatrix.Mul4(t.Rotation.Mat4())
	return transformMatrix
}

// TODO fill in actual interpolation. This just keeps the first keyframe
func InterpolatePoses(k1, k2 *KeyFrame, progression float64) map[int]mgl32.Mat4 {
	interpolatedPose := map[int]mgl32.Mat4{}
	for jointID, jointTransform := range k1.Pose {
		interpolatedPose[jointID] = CalculateJointTransformMatrix(jointTransform)
		// BUG: my self extracted translations and rotations don't recompute the transform matrix properly
		// for now, using the transform matrix directly from the collada document

		// if jointID == 0 {
		// 	if CalculateJointTransformMatrix(jointTransform) != jointTransform.Transform {
		// 		fmt.Println(jointID)
		// 		fmt.Println(jointTransform.Translation)
		// 		fmt.Println(jointTransform.Rotation)
		// 		fmt.Println(CalculateJointTransformMatrix(jointTransform))
		// 		fmt.Println(jointTransform.Transform)
		// 		panic("WHOA")
		// 	} else {
		// 		fmt.Println("OKAY", jointID)
		// 	}
		// }
		interpolatedPose[jointID] = jointTransform.Transform
	}
	return interpolatedPose
}
