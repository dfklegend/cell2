package skilldef

const (
	FormulaInvalid = iota
	FromulaNoDmg
	FromulaNormal
	FormulaMelee
	FormulaRange
)

var (
	FormulaNames = []string{
		"",
		"无伤害",
		"普通",
		"物理近战",
		"物理远程",
	}
)
