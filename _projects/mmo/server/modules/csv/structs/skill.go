package structs

import (
	"strings"

	"mmo/modules/fight/common/skilldef"
)

// ----

type SkillType struct {
	Type int
}

func (t *SkillType) UnmarshalCSV(csv string) (err error) {
	t.Type = GetIndexFromStrings(skilldef.SkillTypeNames, csv, skilldef.SkillTypeInvalid)
	return nil
}

// ----

type TarType struct {
	Type int
}

func (t *TarType) UnmarshalCSV(csv string) (err error) {
	t.Type = GetIndexFromStrings(skilldef.TarTypeNames, csv, skilldef.TarTypeInvalid)
	return nil
}

// ----

type DmgType struct {
	Type int
}

func (t *DmgType) UnmarshalCSV(csv string) (err error) {
	t.Type = GetIndexFromStrings(skilldef.DmgTypeNames, csv, skilldef.DmgTypeInvalid)
	return nil
}

// ----

type SubSkillsType struct {
	Skills []string
}

func (t *SubSkillsType) UnmarshalCSV(csv string) (err error) {
	if strings.TrimSpace(csv) == "" {
		return nil
	}
	t.Skills = strings.Split(csv, ",")
	return nil
}
