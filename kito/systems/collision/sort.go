package collision

import "github.com/kkevinchou/kito/lib/collision"

type contactsBySeparatingDistance []*collision.Contact

func (c contactsBySeparatingDistance) Len() int {
	return len(c)
}
func (c contactsBySeparatingDistance) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c contactsBySeparatingDistance) Less(i, j int) bool {
	return c[i].SeparatingDistance < c[j].SeparatingDistance
}
