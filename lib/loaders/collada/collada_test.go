package collada_test

import (
	"fmt"
	"testing"

	"github.com/kkevinchou/kito/lib/loaders/collada"
)

func TestCollada(t *testing.T) {
	c, err := collada.ParseCollada("sample/model.dae")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(c.TextureSourceData)
	// fmt.Println(c.Root.Children[0])

	t.Error()
}
