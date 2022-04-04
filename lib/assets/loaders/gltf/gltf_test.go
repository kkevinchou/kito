package gltf_test

import (
	"fmt"
	"testing"

	"github.com/kkevinchou/kito/lib/assets/loaders/gltf"
)

// bug hint: when a joint is defined but has no poses our
// animation loading code freaks out. i removed the joint animatiosn from the legs
// and it seems to point to the origin afterwards
// this means the original animation looked wonky probably because there was no pose info
// for the joint which our animation loading code did not understand. likely need to see
// how we handled poses where a joint does not have any poses
// C:\Users\kkevi\goprojects\kito\lib\assets\loaders\gltf\gltf_test.go
func TestBasic(t *testing.T) {
	m, err := gltf.ParseGLTF("../../../../_assets/gltf/cube_anim.gltf")
	if err != nil {
		t.Error(err)
	}
	_ = m

	// fmt.Println(m.Meshes)
	// fmt.Println(m.Meshes[0].MeshChunks)
	chunk := m.Meshes[0].MeshChunks[0]
	fmt.Println(len(chunk.VertexIndices))
	fmt.Println(len(chunk.UniqueVertices))
	// fmt.Println(chunk.VertexIndices)
	// fmt.Println(len(chunk.VertexIndices))
	// fmt.Println(chunk.Vertices)
	// for _, v := range chunk.Vertices {
	// 	fmt.Println(v)
	// }
	// fmt.Println(len(chunk.Vertices))

	// fmt.Println(m.Meshes[0].MeshChunks[0].Vertices)

	t.Fail()
}
