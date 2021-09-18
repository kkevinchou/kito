package model

import (
	"sort"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/modelspec"
)

// if we exceed maxWeights, drop the weakest weights and normalize
// if we're below maxWeights, fill in dummy weights so we always have "maxWeights" number of weights
func FillWeights(jointIDs []int, weights []int, jointWeightsSourceData []float32, maxWeights int) ([]int, []float32) {
	j := []int{}
	w := []float32{}

	if len(jointIDs) <= maxWeights {
		j = append(j, jointIDs...)
		for _, weightIndex := range weights {
			w = append(w, jointWeightsSourceData[weightIndex])
		}
		// fill in empty jointIDs and weights
		for i := 0; i < maxWeights-len(jointIDs); i++ {
			j = append(j, 0)
			w = append(w, 0)
		}
	} else if len(jointIDs) > maxWeights {
		jointWeights := []JointWeight{}
		for i := range jointIDs {
			jointWeights = append(jointWeights, JointWeight{JointID: jointIDs[i], Weight: jointWeightsSourceData[weights[i]]})
		}
		sort.Sort(sort.Reverse(byWeights(jointWeights)))

		// take top 3 weights
		jointWeights = jointWeights[:maxWeights]
		NormalizeWeights(jointWeights)
		for _, jw := range jointWeights {
			j = append(j, jw.JointID)
			w = append(w, jw.Weight)
		}
	}

	return j, w
}

func NormalizeWeights(jointWeights []JointWeight) {
	var totalWeight float32
	for _, jw := range jointWeights {
		totalWeight += jw.Weight
	}

	for i := range jointWeights {
		jointWeights[i].Weight /= totalWeight
	}
}

type byWeights []JointWeight

type JointWeight struct {
	JointID int
	Weight  float32
}

func (s byWeights) Len() int {
	return len(s)
}
func (s byWeights) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byWeights) Less(i, j int) bool {
	return s[i].Weight < s[j].Weight
}

func JointSpecToJoint(js *modelspec.JointSpecification) *Joint {
	j := NewJoint(js.ID, js.Name, js.BindTransform)
	for _, c := range js.Children {
		j.Children = append(j.Children, JointSpecToJoint(c))
	}
	calculateInverseBindTransform(j, mgl32.Ident4())

	return j
}

func calculateInverseBindTransform(joint *Joint, parentBindTransform mgl32.Mat4) {
	bindTransform := parentBindTransform.Mul4(joint.LocalBindTransform) // model-space relative to the origin
	joint.InverseBindTransform = bindTransform.Inv()
	for _, child := range joint.Children {
		calculateInverseBindTransform(child, bindTransform)
	}
}

func getJointMap(joint *Joint, jointMap map[int]*Joint) map[int]*Joint {
	jointMap[joint.ID] = joint
	for _, c := range joint.Children {
		getJointMap(c, jointMap)
	}
	return jointMap
}
