package animation

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/libutils"
	"github.com/kkevinchou/kito/lib/model"
	"github.com/kkevinchou/kito/lib/modelspec"
)

type AnimationPlayer struct {
	// stateful data that is manipulated by the Animation System
	elapsedTime         time.Duration
	animationTransforms map[int]mgl32.Mat4
	currentAnimation    string

	// these fields are from the loaded animation and should not be modified
	// Animation  *model.Animation
	animations map[string]*model.Animation
}

func NewAnimationPlayer(animations map[string]*model.Animation) *AnimationPlayer {
	return &AnimationPlayer{
		animations: animations,
	}
}

func (player *AnimationPlayer) AnimationTransforms() map[int]mgl32.Mat4 {
	return player.animationTransforms
}

func (player *AnimationPlayer) PlayAnimation(animationName string) {
	if player.currentAnimation == animationName {
		return
	}

	if _, ok := player.animations[animationName]; ok {
		// player.Animation = animation
		player.currentAnimation = animationName
		player.elapsedTime = 0
	}
}

func (player *AnimationPlayer) Update(delta time.Duration) {
	currAnim := player.animations[player.currentAnimation]

	player.elapsedTime += delta
	for player.elapsedTime.Milliseconds() > currAnim.Length().Milliseconds() {
		player.elapsedTime = time.Duration(player.elapsedTime.Milliseconds()-currAnim.Length().Milliseconds()) * time.Millisecond
	}

	pose := calculateCurrentAnimationPose(player.elapsedTime, currAnim.KeyFrames())
	animationTransforms := computeJointTransforms(currAnim.RootJoint(), pose)
	player.animationTransforms = animationTransforms
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
