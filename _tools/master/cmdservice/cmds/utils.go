package cmds

import (
	"fmt"
	"strings"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/node/builtin/msgs"
	"github.com/dfklegend/cell2/node/cluster"
)

func DefaultCmdDo(b ICmdBuilder, pid *actor.PID, id string, member *cluster.Member) {
	b.GetService().RequestEx(pid, "ctrl.cmd", &msgs.CtrlCmd{
		Cmd: b.GetCmd(),
	}, func(err error, raw interface{}) {
		if err != nil {
			nodeName := getNodeFromId(id)
			result := fmt.Sprintf(" %v error: %v", nodeName, err)
			b.AppendResult(nodeName, result)
			b.Done()
			return
		}
		ack := raw.(*msgs.CtrlCmdAck)

		nodeName := getNodeFromId(id)
		result := fmt.Sprintf(" %v %v", nodeName, ack.Result)
		b.AppendResult(nodeName, result)
		b.Done()
	})
}

func getNodeFromId(id string) string {
	subs := strings.Split(id, "@")
	if len(subs) < 2 {
		return subs[0]
	}
	return subs[1]
}
