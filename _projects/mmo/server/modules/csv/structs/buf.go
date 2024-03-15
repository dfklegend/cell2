package structs

import (
	"mmo/modules/fight/common/bufdef"
)

type BufSkillTarType struct {
	Type int
}

func (t *BufSkillTarType) UnmarshalCSV(csv string) (err error) {
	t.Type = GetIndexFromStrings(bufdef.TarTypeNames, csv, bufdef.TarTypeInvalid)
	return nil
}
