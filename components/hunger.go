package components

type HungerComponent struct {
	Value int
}

func (h *HungerComponent) SetHunger(value int) {
	h.Value = value
}

func (h *HungerComponent) GetHunger() int {
	return h.Value
}

func (h *HungerComponent) UpdateHunger(delta int) {
	h.Value += delta
}
