package singleton

import (
	"github.com/kkevinchou/kito/components/singleton"
)

type Singleton struct {
	*singleton.InputComponent
}

func New() *Singleton {
	return &Singleton{
		InputComponent: singleton.New(),
	}
}
