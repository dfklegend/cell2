package modes

import (
	"github.com/dfklegend/cell2/node/app"
	l "github.com/dfklegend/cell2/utils/logger"

	mymsg "mmo/messages"
	"mmo/messages/cproto"
	"mmo/servers/logic_old/define"
	"mmo/servers/logic_old/fightmode"
)

type CardFight struct {
	*BaseMode

	Cards []*cproto.CharCard
}

func NewCardFight() *CardFight {
	return &CardFight{
		BaseMode: newBaseMode(),
	}
}

func (f *CardFight) InitData(player define.ILogicPlayer, d fightmode.ISceneFightData) {
	f.BaseMode.InitData(player, d)

	data := d.(*fightmode.CardFightData)
	f.Cards = data.Cards
}

func (f *CardFight) SendInitDataToScene(cb func()) {
	ns := f.player.GetNodeService()
	app.Request(f.player.GetNodeService(), "scene.remote.initcardfight", f.sceneService, &mymsg.SceneInitCardFight{
		UId:     f.uid,
		LogicId: ns.Name,
		SceneId: f.sceneId,
		Token:   f.token,
		Cards:   f.Cards,
	}, func(err error, ret any) {
		cb()
	})
}

func (f *CardFight) OnFightResult(result fightmode.ISceneFightResult) {
	data := result.(*mymsg.SceneCardResult)
	l.L.Infof("card result: %v", data)
	if data.Winner < 0 || int(data.Winner) >= len(f.Cards) {
		return
	}

	for i := 0; i < len(f.Cards); i++ {
		winner := i == int(data.Winner)
		id := f.Cards[i].Id
		card := f.player.GetCharCard().GetCard(id)
		card.Stat.Total++
		if winner {
			card.Stat.Win++
		}
	}

	f.player.GetCharCard().UpdateAndRefresh()
}
