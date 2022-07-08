package components

type Action string

var ActionCast Action = "CAST"
var ActionNone Action = "NONE"

// Quick and dirty component meant to store features that are still being interated on
type NotepadComponent struct {
	LastAction Action
}

func (c *NotepadComponent) GetNotepadComponent() *NotepadComponent {
	return c
}

func (c *NotepadComponent) AddToComponentContainer(container *ComponentContainer) {
	container.NotepadComponent = c
}

func (c *NotepadComponent) ComponentFlag() int {
	return ComponentFlagNotepad
}
