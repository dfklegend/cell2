package common

const (
	BufAffectNone int = iota
	BufAffectGood     // 有益
	BufAffectBad      // 有害
)

type IBufCtrl interface {
	Destroy()

	AddBuf(caster ICharacter, owner ICharacter, id BufId, level int, stack int)
	RemoveBuf(id BufId)
	HasBuf(id BufId) bool

	// GetBuf level == 0, 代表没有
	GetBuf(id BufId) (level int, stack int)
	Update()
}
