package common

// IPlugin 技能,buf可以自带一个plugin
// plugin通过监听玩家事件来执行一些逻辑
//
// 比如: 某个技能3下攻击之后，必爆击一次
type IPlugin interface {
	Init(character ICharacter)
	OnStart()
	OnTriggle()
	OnStop()
}
