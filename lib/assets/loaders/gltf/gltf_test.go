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
	m, err := gltf.ParseGLTF("../../../../_assets/gltf/alpha.gltf", &gltf.ParseConfig{TextureCoordStyle: gltf.TextureCoordStyleOpenGL})
	if err != nil {
		t.Error(err)
	}
	_ = m

	fmt.Println(m.RootTransforms)

	t.Fail()
}
