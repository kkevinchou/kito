package main

import (
	"time"

	"github.com/kkevinchou/ant/behavior/worker"
)

func main() {
	a := worker.New(1)
	a.Tick(1 * time.Second)
}
