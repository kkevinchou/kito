package gltf_test

import (
	"testing"

	"github.com/kkevinchou/kito/lib/assets/loaders/gltf"
)

func TestBasic(t *testing.T) {
	_, err := gltf.ParseGLTF("sample/basic.glb")
	if err != nil {
		t.Error(err)
	}

	t.Fail()
}
