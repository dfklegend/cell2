package components

import (
	"mmo/servers/scene/define"
)

type PlayerComponent struct {
	*BaseSceneComponent
	player define.IPlayer
}

func NewPlayerComponent(player define.IPlayer) *PlayerComponent {
	return &PlayerComponent{
		BaseSceneComponent: NewBaseSceneComponent(),
		player:             player,
	}
}

func (p *PlayerComponent) GetPlayer() define.IPlayer {
	return p.player
}

func (p *PlayerComponent) OnPrepare() {
	p.BaseSceneComponent.OnPrepare()
}

func (p *PlayerComponent) OnStart() {
}

func (p *PlayerComponent) Update() {
}

func (p *PlayerComponent) LateUpdate() {
}

func (p *PlayerComponent) OnDestroy() {
}
