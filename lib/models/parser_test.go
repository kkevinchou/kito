package models

import "testing"

func TestParsing(t *testing.T) {
	file := "../../_assets/obj/Oak_Green_01.obj"
	model, err := NewModel(file)
	if err != nil {
		t.Error(err)
	}

	if len(model.Faces) == 0 {
		t.Errorf("Could not load faces")
	}
}
