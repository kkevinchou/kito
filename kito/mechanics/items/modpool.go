package items

import "math/rand"

type AffixType string

const (
	AffixTypePrefix AffixType = "PREFIX"
	AffixTypeSuffix AffixType = "SUFFIX"
)

type ModPool struct {
	prefixPool map[int]*Mod
	suffixPool map[int]*Mod

	prefixList []*Mod
	suffixList []*Mod
}

func NewModPool() *ModPool {
	return &ModPool{
		prefixPool: map[int]*Mod{},
		suffixPool: map[int]*Mod{},
		prefixList: []*Mod{},
		suffixList: []*Mod{},
	}
}

func (m *ModPool) AddMod(mod *Mod) {
	if mod.AffixType == AffixTypePrefix {
		m.prefixPool[mod.ID] = mod
		m.prefixList = append(m.prefixList, mod)
	} else if mod.AffixType == AffixTypeSuffix {
		m.suffixPool[mod.ID] = mod
		m.suffixList = append(m.suffixList, mod)
	}
}

func (m *ModPool) ChooseMods(count int, maxPrefix int, maxSuffix int) []*Mod {
	if maxPrefix+maxSuffix > count {
		panic("max prefix and max suffix cannot exceed count")
	}

	prefixCount := rand.Intn(maxPrefix + 1)
	suffixCount := count - prefixCount

	// TOOD: shuffle

	mods := []*Mod{}
	if prefixCount > 0 {
		prefixes := m.prefixList[0:prefixCount]
		mods = append(mods, prefixes...)
	}
	if suffixCount > 0 {
		suffixes := m.suffixList[0:suffixCount]
		mods = append(mods, suffixes...)
	}

	return mods
}

type Mod struct {
	ID        int
	AffixType AffixType
	Effect    Effect
}

type Effect interface {
	ApplyDamage()
}
