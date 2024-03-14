package factory

import (
	"mlrs/b3"
	"mlrs/b3/actions"
	"mlrs/b3/composites"
	"mlrs/b3/config"
	"mlrs/b3/core"
	"mlrs/b3/decorators"
)

func createBaseStructMaps() *b3.RegisterStructMaps {
	st := b3.NewRegisterStructMaps()
	//actions
	st.Register("Error", &actions.Error{})
	st.Register("Failure", &actions.Failure{})
	st.Register("Running", &actions.Running{})
	st.Register("Success", &actions.Success{})
	st.Register("Wait", &actions.Wait{})
	st.Register("Log", &actions.Log{})
	//composites
	st.Register("MemPriority", &composites.MemPriority{})
	st.Register("MemSequence", &composites.MemSequence{})
	st.Register("Priority", &composites.Priority{})
	st.Register("Sequence", &composites.Sequence{})

	//decorators
	st.Register("Inverter", &decorators.Inverter{})
	st.Register("Limiter", &decorators.Limiter{})
	st.Register("MaxTime", &decorators.MaxTime{})
	st.Register("Repeater", &decorators.Repeater{})
	st.Register("RepeatUntilFailure", &decorators.RepeatUntilFailure{})
	st.Register("RepeatUntilSuccess", &decorators.RepeatUntilSuccess{})
	return st
}

var baseMap = createBaseStructMaps()
var registerMap = b3.NewRegisterStructMaps()

func RegisterMap() *b3.RegisterStructMaps {
	return registerMap
}

func CreateBevTreeFromConfig(config *config.BTTreeCfg, extMap *b3.RegisterStructMaps) *core.BehaviorTree {
	tree := core.NewBeTree()
	tree.Load(config, createBaseStructMaps(), extMap)
	return tree
}

func CreateBevTreeFromProjectData(projectData *config.BTProjectCfg) *core.BehaviorTree {

	var tree *core.BehaviorTree
	for _, treeConfig := range projectData.Trees {
		tree = CreateBevTreeFromConfig(&treeConfig, registerMap)
	}
	return tree
}
