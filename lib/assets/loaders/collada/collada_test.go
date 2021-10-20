package collada_test

import (
	"fmt"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/assets/loaders/collada"
	"github.com/kkevinchou/kito/lib/modelspec"
)

// cowboy model
// 1420 faces
// 4260 individual vertices (vertices can be counted multiple times)
// 740 distinct vertices

func TestManWalk(t *testing.T) {
	c, err := collada.ParseCollada("sample/man_walk.dae")
	// printKeyFrames(c)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(c)
	t.Fail()
}

func TestBox(t *testing.T) {
	c, err := collada.ParseCollada("sample/box.dae")
	// printKeyFrames(c)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(c)
	t.Fail()
}

func TestCollada(t *testing.T) {
	c, err := collada.ParseCollada("sample/bob.dae")
	// printKeyFrames(c)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(c)
	t.Fail()
}

func TestCowboy(t *testing.T) {
	c, err := collada.ParseCollada("sample/model.dae")
	// printKeyFrames(c)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(c)
	t.Fail()
}

func printKeyFrames(c *modelspec.ModelSpecification) {
	for _, kf := range c.Animation.KeyFrames {
		fmt.Println(kf.Start, kf.Pose)
	}
	fmt.Println(c.Animation.Length)
}

func TestJointHierarchy(t *testing.T) {
	c, err := collada.ParseCollada("sample/model2.dae")
	if err != nil {
		t.Fatal(err)
	}

	jointIDs := []int{}
	collectIDs(c.RootJoint, &jointIDs)

	for i, jointID := range jointIDs {
		if i != jointID {
			t.Fatalf("joint at position %d does not match joint ID: %d", i, jointID)
		}
	}

	transforms := map[int]mgl32.Mat4{}
	collectBindPoseAnimationTransforms(c.RootJoint, transforms)

	if len(transforms) != 16 {
		t.Fatalf("expected 16 transforms but intead got %d", len(transforms))
	}
}

func collectBindPoseAnimationTransforms(joint *modelspec.JointSpec, transforms map[int]mgl32.Mat4) {
	transforms[joint.ID] = joint.BindTransform
	for _, child := range joint.Children {
		collectBindPoseAnimationTransforms(child, transforms)
	}
}

func collectIDs(joint *modelspec.JointSpec, ids *[]int) {
	*ids = append(*ids, joint.ID)
	for _, child := range joint.Children {
		collectIDs(child, ids)
	}
}
