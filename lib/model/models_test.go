package model_test

import (
	"fmt"
	"testing"

	"github.com/kkevinchou/kito/lib/model"
)

func TestFillWeightsWithNoChange(t *testing.T) {
	jointIDs := []int{0, 1, 2}
	weights := []int{0, 1, 2}
	jointWeightsSourceData := []float32{0.55, 0.25, 0.20}
	maxWeights := 3

	newJointIDs, newWeights := model.FillWeights(jointIDs, weights, jointWeightsSourceData, maxWeights)
	if !IntSliceEqual(jointIDs, newJointIDs) {
		t.Fatal("expected jointIDs to match")
	}

	if !Float32SliceEqual(jointWeightsSourceData, newWeights) {
		t.Fatal("expected weights to match")
	}

}

func TestFillWeightsDroppingWeight(t *testing.T) {
	jointIDs := []int{0, 1, 2, 3}
	weights := []int{0, 1, 2, 3}
	jointWeightsSourceData := []float32{0.55, 0.25, 0.10, 0.20}
	maxWeights := 3

	newJointIDs, newWeights := model.FillWeights(jointIDs, weights, jointWeightsSourceData, maxWeights)

	expectedJointIDs := []int{0, 1, 3}
	if !IntSliceEqual(expectedJointIDs, newJointIDs) {
		t.Fatal("expected jointIDs to match")
	}

	expectedWeights := []float32{0.55, 0.25, 0.2}
	if !Float32SliceEqual(expectedWeights, newWeights) {
		t.Fatal("expected weights to match")
	}
}

func TestFillWeightsWithAddedWeight(t *testing.T) {
	jointIDs := []int{0, 1}
	weights := []int{0, 1}
	jointWeightsSourceData := []float32{0.75, 0.25}
	maxWeights := 3

	newJointIDs, newWeights := model.FillWeights(jointIDs, weights, jointWeightsSourceData, maxWeights)

	expectedJointIDs := []int{0, 1, 0}
	if !IntSliceEqual(expectedJointIDs, newJointIDs) {
		t.Fatal("expected jointIDs to match")
	}

	expectedWeights := []float32{0.75, 0.25, 0}
	if !Float32SliceEqual(expectedWeights, newWeights) {
		t.Fatal("expected weights to match")
	}

	if len(newJointIDs) != maxWeights {
		t.Fatal("expected length of joint ids to match maxWeights")

	}
}
func TestNormalizeWeightsHappyPath(t *testing.T) {
	jointWeights := []model.JointWeight{
		model.JointWeight{JointID: 0, Weight: 0.55},
		model.JointWeight{JointID: 0, Weight: 0.25},
		model.JointWeight{JointID: 0, Weight: 0.20},
	}

	model.NormalizeWeights(jointWeights)
	var expected float32 = 0.55
	if jointWeights[0].Weight != expected {
		t.Fatal(fmt.Sprintf("joint weight should be %f but was instead: %f", expected, jointWeights[0].Weight))
	}
	expected = 0.25
	if jointWeights[1].Weight != expected {
		t.Fatal(fmt.Sprintf("joint weight should be %f but was instead: %f", expected, jointWeights[1].Weight))
	}
	expected = 0.20
	if jointWeights[2].Weight != expected {
		t.Fatal(fmt.Sprintf("joint weight should be %f but was instead: %f", expected, jointWeights[2].Weight))
	}
}

func TestNormalizeWeightsWithAdjustments(t *testing.T) {
	jointWeights := []model.JointWeight{
		model.JointWeight{JointID: 0, Weight: 0.25},
		model.JointWeight{JointID: 0, Weight: 0.25},
	}

	model.NormalizeWeights(jointWeights)
	var expected float32 = 0.5
	if jointWeights[0].Weight != expected {
		t.Fatal(fmt.Sprintf("joint weight should be %f but was instead: %f", expected, jointWeights[0].Weight))
	}
	expected = 0.5
	if jointWeights[1].Weight != expected {
		t.Fatal(fmt.Sprintf("joint weight should be %f but was instead: %f", expected, jointWeights[1].Weight))
	}
}

func IntSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func Float32SliceEqual(a, b []float32) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
