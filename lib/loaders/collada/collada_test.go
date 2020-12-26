package collada_test

import (
	"fmt"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/animation"
	"github.com/kkevinchou/kito/lib/loaders/collada"
)

func TestCollada(t *testing.T) {
	_, err := collada.ParseCollada("sample/model.dae")
	if err != nil {
		t.Fatal(err)
	}
}

func TestJointHierarchy(t *testing.T) {
	c, err := collada.ParseCollada("sample/model.dae")
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

func TestJointHierarchy2(t *testing.T) {
	c, err := collada.ParseCollada("sample/model.dae")
	if err != nil {
		t.Fatal(err)
	}

	var maxZ float32
	verticesAffectedByJoint := []int{}
	for i, ids := range c.JointIDs {
		if c.PositionSourceData[i].Z() > maxZ {
			maxZ = c.PositionSourceData[i].Z()
			fmt.Println(maxZ)
		}
		for _, id := range ids {
			if id == 15 {
				verticesAffectedByJoint = append(verticesAffectedByJoint, i)
			}
		}
	}

	for _, index := range verticesAffectedByJoint {
		v := c.PositionSourceData[index]
		fmt.Println(index, v)
	}

	// fmt.Println(verticesAffectedByJoint)
	t.Fail()
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
