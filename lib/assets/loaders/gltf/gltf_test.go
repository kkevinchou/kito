package gltf_test

import (
	"testing"

	"github.com/kkevinchou/kito/lib/assets/loaders/gltf"
)

func TestBasic(t *testing.T) {
	_, err := gltf.ParseGLTF("sample/cube_anim.gltf")
	if err != nil {
		t.Error(err)
	}

	t.Fail()
}
