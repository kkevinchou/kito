package main

import (
	"time"

	"github.com/kkevinchou/ant/behavior"
)

func main() {
	a := behavior.Sequence{}
	a.Tick(time.Second * 1)
}
