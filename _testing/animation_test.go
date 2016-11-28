package sometest

import (
	"testing"
	"time"

	"github.com/kkevinchou/kito/animation"
)

func TestBlah(t *testing.T) {
	a := animation.Load("../assets/animations/ant")
	a.Update(time.Second * 1)
	t.Fail()
}
