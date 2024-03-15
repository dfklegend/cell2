package csv

import (
	"fmt"
)

var (
	Skill           = NewSkillCfg()
	SkillEffect     = NewSkillEffectCfg()
	SkillEffects    = NewSkillEffects()
	Buf             = NewBufCfg()
	Equip           = NewEquipCfg()
	Scene           = NewSceneCfg()
	MonsterTemplate = NewMonsterTemplateCfg()
	Monster         = NewMonsterCfg()
)

func init() {

}

func getFinalPath(dir string, file string) string {
	return fmt.Sprintf("%v/%v", dir, file)
}

func LoadAll(dir string) {
	Skill.LoadFromFile(getFinalPath(dir, "skill.csv"))
	SkillEffect.LoadFromFile(getFinalPath(dir, "skilleffect.csv"))
	SkillEffects.Build(Skill, SkillEffect)
	Buf.LoadFromFile(getFinalPath(dir, "buf.csv"))
	Equip.LoadFromFile(getFinalPath(dir, "equip.csv"))
	Scene.LoadFromFile(getFinalPath(dir, "scene.csv"))
	MonsterTemplate.LoadFromFile(getFinalPath(dir, "monster_template.csv"))
	Monster.LoadFromFile(getFinalPath(dir, "monster.csv"))
}
