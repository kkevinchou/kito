package animation

import (
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/libutils"
	"github.com/kkevinchou/kito/lib/model"
	"github.com/kkevinchou/kito/lib/modelspec"
)

type AnimationPlayer struct {
	elapsedTime         time.Duration
	animationTransforms map[int]mgl32.Mat4
	currentAnimation    *modelspec.AnimationSpec

	// these fields are from the loaded animation and should not be modified
	animations map[string]*modelspec.AnimationSpec
	rootJoint  *modelspec.JointSpec

	secondaryAnimation *string
	loop               bool
	blendActive        bool
}

// func NewAnimationPlayer(animations map[string]*modelspec.AnimationSpec, rootJoint *modelspec.JointSpec) *AnimationPlayer {
func NewAnimationPlayer(m *model.Model) *AnimationPlayer {
	return &AnimationPlayer{
		animations: m.Animations(),
		rootJoint:  m.RootJoint(),
		loop:       true,
	}
}

func (player *AnimationPlayer) CurrentAnimation() string {
	if player.currentAnimation == nil {
		return ""
	}
	return player.currentAnimation.Name
}

func (player *AnimationPlayer) AnimationTransforms() map[int]mgl32.Mat4 {
	return player.animationTransforms
}

func (player *AnimationPlayer) PlayAnimation(animationName string) {
	if player.currentAnimation != nil && player.currentAnimation.Name == animationName {
		return
	}
	if player.secondaryAnimation != nil && animationName == *player.secondaryAnimation {
		return
	}

	if currentAnimation, ok := player.animations[animationName]; ok {
		player.currentAnimation = currentAnimation
		player.elapsedTime = 0
	} else {
		panic(fmt.Sprintf("failed to find animation %s", animationName))
	}
}

func (player *AnimationPlayer) PlayAndBlendAnimation(animationName string, blendTime time.Duration) {
}

func (player *AnimationPlayer) PlayOnce(animationName string, secondaryAnimation string) {
	local := secondaryAnimation
	player.secondaryAnimation = &local

	if currentAnimation, ok := player.animations[animationName]; ok {
		player.currentAnimation = currentAnimation
		player.elapsedTime = 0
		player.loop = false
	} else {
		panic(fmt.Sprintf("failed to find animation %s", animationName))
	}
}

func (player *AnimationPlayer) Update(delta time.Duration) {
	if player.currentAnimation == nil {
		return
	}

	// TODO(kevin): this code ugly af
	player.elapsedTime += delta
	for player.elapsedTime.Milliseconds() > player.currentAnimation.Length.Milliseconds() {
		player.elapsedTime = time.Duration(player.elapsedTime.Milliseconds()-player.currentAnimation.Length.Milliseconds()) * time.Millisecond

		// if we're not looping, we should have a secondary animation to fall back into
		if !player.loop {
			player.PlayAnimation(*player.secondaryAnimation)
			player.loop = true
			player.secondaryAnimation = nil
		}
	}

	player.animationTransforms = player.computeAnimationTransforms(player.elapsedTime, player.currentAnimation)
}

func (player *AnimationPlayer) computeAnimationTransforms(elapsedTime time.Duration, animation *modelspec.AnimationSpec) map[int]mgl32.Mat4 {
	pose := calculateCurrentAnimationPose(player.elapsedTime, animation.KeyFrames)
	animationTransforms := computeJointTransforms(player.rootJoint, pose)
	return animationTransforms
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
