package skilleffect

import (
	"mmo/modules/csv/structs"
	"mmo/modules/fight/common/skilleffect"
)

// ----

type ApplyTime struct {
	Type int
}

func (t *ApplyTime) UnmarshalCSV(csv string) (err error) {
	t.Type = structs.GetIndexFromStrings(skilleffect.ApplyTimeNames, csv, skilleffect.ATInvalid)
	return nil
}

// ----

type TarType struct {
	Type int
}

func (t *TarType) UnmarshalCSV(csv string) (err error) {
	t.Type = structs.GetIndexFromStrings(skilleffect.TarTypeNames, csv, skilleffect.TarTypeInvalid)
	return nil
}

// ----

type OpType struct {
	Type int
}

func (t *OpType) UnmarshalCSV(csv string) (err error) {
	t.Type = structs.GetIndexFromStrings(skilleffect.OpNames, csv, skilleffect.OpInvalid)
	return nil
}
