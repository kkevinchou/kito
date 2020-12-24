package collada_test

import (
	"testing"

	"github.com/kkevinchou/kito/lib/loaders/collada"
)

func TestCollada(t *testing.T) {
	c, err := collada.ParseCollada("sample/model.dae")
	if err != nil {
		t.Fatal(err)
	}
	_ = c

	t.Error()
}
