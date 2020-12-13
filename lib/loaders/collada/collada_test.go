package collada_test

import (
	"fmt"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/loaders/collada"
)

func TestCollada(t *testing.T) {
	rawCollada, err := collada.LoadDocument("sample/model.dae")
	if err != nil {
		t.Error(err)
	}

	collada := collada.ParseCollada(rawCollada)
	fmt.Println(len(collada.Vertices))
	fmt.Println(len(collada.Normals))
	t.Error()
}

func parseNormals(rawCollada *collada.RawCollada) []mgl32.Vec3 {
	return nil
}

func parseAnimations(rawCollada *collada.RawCollada) {
}
