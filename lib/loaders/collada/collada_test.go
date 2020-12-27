package collada_test

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/animation"
	"github.com/kkevinchou/kito/lib/loaders/collada"
)

func TestCollada(t *testing.T) {
	_, err := collada.ParseCollada("sample/model2.dae")
	if err != nil {
		t.Fatal(err)
	}

	t.Fail()
}

func TestJointHierarchy(t *testing.T) {
	c, err := collada.ParseCollada("sample/model2.dae")
	if err != nil {
		t.Fatal(err)
	}

	jointIDs := []int{}
	collectIDs(c.Root, &jointIDs)

	for i, jointID := range jointIDs {
		if i != jointID {
			t.Fatalf("joint at position %d does not match joint ID: %d", i, jointID)
		}
	}

	transforms := map[int]mgl32.Mat4{}
	collectBindPoseAnimationTransforms(c.Root, transforms)

	if len(transforms) != 16 {
		t.Fatalf("expected 16 transforms but intead got %d", len(transforms))
	}
}

func collectBindPoseAnimationTransforms(joint *animation.JointSpecification, transforms map[int]mgl32.Mat4) {
	transforms[joint.ID] = joint.BindTransform
	for _, child := range joint.Children {
		collectBindPoseAnimationTransforms(child, transforms)
	}
}

func collectIDs(joint *animation.JointSpecification, ids *[]int) {
	*ids = append(*ids, joint.ID)
	for _, child := range joint.Children {
		collectIDs(child, ids)
	}
}
