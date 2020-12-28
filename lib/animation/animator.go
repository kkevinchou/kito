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
	for a.ElapsedTime.Milliseconds() > a.Animation.Length.Milliseconds() {
		a.ElapsedTime = time.Duration(a.ElapsedTime.Milliseconds()-a.Animation.Length.Milliseconds()) * time.Millisecond
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

func (a *Animator) PlayAnimation(animation *Animation) {
	a.Animation = animation
	a.ElapsedTime = 0
}

func (a *Animator) calculateCurrentAnimationPose() map[int]mgl32.Mat4 {
	var startKeyFrame *KeyFrame
	var endKeyFrame *KeyFrame
	var progression float32

	// iterate backwards looking for the starting keyframe
	for i := len(a.Animation.KeyFrames) - 1; i >= 0; i-- {
		keyFrame := a.Animation.KeyFrames[i]
		if a.ElapsedTime >= keyFrame.Start || i == 0 {
			startKeyFrame = keyFrame
			if i < len(a.Animation.KeyFrames)-1 {
				endKeyFrame = a.Animation.KeyFrames[i+1]
				progression = float32((a.ElapsedTime - startKeyFrame.Start).Milliseconds()) / float32((endKeyFrame.Start - startKeyFrame.Start).Milliseconds())
			} else {
				// interpolate towards the first kf, assume looping animations
				endKeyFrame = a.Animation.KeyFrames[0]
				progression = 0
			}
			break
		}
	}

	return InterpolatePoses(startKeyFrame, endKeyFrame, progression)
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
func InterpolatePoses(k1, k2 *KeyFrame, progression float32) map[int]mgl32.Mat4 {
	interpolatedPose := map[int]mgl32.Mat4{}
	for jointID := range k1.Pose {
		k1JointTransform := k1.Pose[jointID]
		k2JointTransform := k2.Pose[jointID]

		// WTF - this lerp doesn't look right when interpolating keyframes???
		// rotationQuat := mgl32.QuatLerp(k1JointTransform.Rotation, k2JointTransform.Rotation, progression)

		rotationQuat := qinterpolate(k1JointTransform.Rotation, k2JointTransform.Rotation, progression)
		rotation := rotationQuat.Mat4()

		translation := k1JointTransform.Translation.Add(k2JointTransform.Translation.Sub(k1JointTransform.Translation).Mul(progression))

		interpolatedPose[jointID] = mgl32.Translate3D(translation.X(), translation.Y(), translation.Z()).Mul4(rotation)
	}
	return interpolatedPose
}

// reimplemented from: https://github.com/TheThinMatrix/OpenGL-Animation/blob/dde792fe29767192bcb60d30ac3e82d6bcff1110/Animation/animation/Quaternion.java#L158
func qinterpolate(a, b mgl32.Quat, blend float32) mgl32.Quat {
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
