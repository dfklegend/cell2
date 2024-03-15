using Phoenix.Core;
using Phoenix.Utils;
using Phoenix.Game.Card;

public static class GMCmdUtils
{
    public static void Output(string content)
    {
        HEventUtil.Dispatch(GlobalEvents.It.events, new HEventBattleLog(content));
    }
}
