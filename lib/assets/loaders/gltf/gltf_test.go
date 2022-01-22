package gltf_test

import (
	"testing"

	"github.com/kkevinchou/kito/lib/assets/loaders/gltf"
)

// bug hint: when a joint is defined but has no poses our
// animation loading code freaks out. i removed the joint animatiosn from the legs
// and it seems to point to the origin afterwards
// this means the original animation looked wonky probably because there was no pose info
// for the joint which our animation loading code did not understand. likely need to see
// how we handled poses where a joint does not have any poses

func TestBasic(t *testing.T) {
	_, err := gltf.ParseGLTF("sample/guard.gltf")
	if err != nil {
		t.Error(err)
	}

	t.Fail()
}
