package kito

type InputPoller func() []interface{}

func NullInputPoller() []interface{} {
	return nil
}
