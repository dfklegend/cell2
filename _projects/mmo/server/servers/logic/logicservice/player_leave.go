package logicservice

func (p *LogicPlayer) OnLogout() {
	p.ns.GetLogger().Infof("%v OnLogout", p.uid)
}
