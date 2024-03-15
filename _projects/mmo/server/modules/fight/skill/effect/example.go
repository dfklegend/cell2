package effect

import (
	"fmt"

	"github.com/dfklegend/cell2/utils/jsonutils"
	l "github.com/dfklegend/cell2/utils/logger"

	"mmo/modules/csv/base"
	"mmo/modules/csv/entry"
	"mmo/modules/fight/common"
	"mmo/modules/fight/common/skilleffect"
)

func init() {
	Register(skilleffect.OpExample, &exampleOp{})
}

type exampleOp struct {
}

// exampleArgs 实现接口 IArgsFormatter
type exampleArgs struct {
	Rate   float32
	Damage int

	//扩展字段处理后的数据结构
	second jsonData
}

type jsonData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (a *exampleOp) Format(cfg *base.IArgs) {
	cfg.FormatArgs(&exampleArgs{
		Rate:   0,  //默认值
		Damage: 10, //默认值
	})

	//处理extArgs，可以约定自己的格式
	extArgs := cfg.ExtArgs
	//测试json处理
	extArgs = "{\"name\":\"tom\",\"age\":18}"
	args := cfg.ArgsImpl.(*exampleArgs)
	var jd jsonData
	jsonutils.Unmarshal([]byte(extArgs), &jd)
	args.second = jd
	fmt.Println(args)
}

func (a *exampleOp) Apply(caster common.ICharacter, tar common.ICharacter, cfg *entry.SkillEffect, skillLv int) {
	if tar == nil {
		return
	}
	args := cfg.ArgsImpl.(*exampleArgs)
	rate := args.Rate
	damage := args.Damage
	l.Log.Infof("rate : %v, damage : %v", rate, damage)
}
