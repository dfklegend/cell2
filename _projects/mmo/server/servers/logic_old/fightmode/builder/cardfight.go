package builder

import (
	"mmo/messages/cproto"
	"mmo/servers/logic_old/define"
	"mmo/servers/logic_old/fightmode"
	"mmo/servers/logic_old/fightmode/modes"
)

func BuildCardFight(player define.ILogicPlayer, card0, card1 *cproto.CharCard) *modes.CardFight {
	data := &fightmode.CardFightData{}
	data.Cards = make([]*cproto.CharCard, 2)
	data.Cards[0] = card0
	data.Cards[1] = card1

	mode := modes.NewCardFight()
	mode.InitData(player, data)
	return mode
}
