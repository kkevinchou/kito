package console

var GlobalConsole *Console = &Console{}

type ConsoleItem struct {
	Command string
	Output  string
}

type Console struct {
	ConsoleItems []*ConsoleItem
	Input        string
}

func (c *Console) Send() {
	c.ConsoleItems = append(c.ConsoleItems, &ConsoleItem{Command: c.Input})
	c.Input = ""
}
