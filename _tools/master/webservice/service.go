package webservice

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/flamego/flamego"

	"github.com/dfklegend/cell2/actorex"
	"github.com/dfklegend/cell2/actorex/service"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/builtin/msgs"
	"github.com/dfklegend/cell2/nodectrl"
	"github.com/dfklegend/cell2/nodectrl/define"
	"github.com/dfklegend/cell2/utils/jsonutils"
	"github.com/dfklegend/cell2/utils/logger"
)

const HttpPort = 8080

type WebService struct {
	*service.Service
	cmdServicePID *actor.PID
}

func (s *WebService) Start() {
	logger.Log.Infof("web service start")
	s.cmdServicePID = actor.NewPID(actorex.LocalAddress, "cmdservice")

	port := app.Node.GetMasterInfo().HttpPort

	go s.startHttpServer(port)
}

func (s *WebService) retire(c flamego.Context) string {
	p := c.Request().PostFormValue("p")
	ch := make(chan string, 1)
	time.AfterFunc(time.Second*3, func() {
		ch <- ""
	})

	s.Post(func() {
		cluster := app.Node.GetCluster()
		member := cluster.GetMembers()[fmt.Sprintf("%v@%v", cluster.GetName(), p)]

		host := member.Host
		port := member.Port
		pid := actor.NewPID(fmt.Sprintf("%v:%v", host, port), define.NodeAdmin)

		s.RequestEx(pid, "ctrl.cmd", &msgs.CtrlCmd{
			Cmd: "web_retire",
		}, func(err error, res any) {
			ctrlCmdAck := res.(*msgs.CtrlCmdAck)
			fmt.Println(ctrlCmdAck.Result)
			ch <- fmt.Sprintf("%v", ctrlCmdAck.Result)
		})
	})

	json := <-ch
	return fmt.Sprintf("[%s]", json)
}

func (s *WebService) exit(c flamego.Context) string {
	p := c.Request().PostFormValue("p")
	ch := make(chan string, 1)
	time.AfterFunc(time.Second*3, func() {
		ch <- ""
	})

	s.Post(func() {
		cluster := app.Node.GetCluster()
		member := cluster.GetMembers()[fmt.Sprintf("%v@%v", cluster.GetName(), p)]

		host := member.Host
		port := member.Port
		pid := actor.NewPID(fmt.Sprintf("%v:%v", host, port), define.NodeAdmin)

		s.RequestEx(pid, "ctrl.cmd", &msgs.CtrlCmd{
			Cmd: "web_exit",
		}, func(err error, res any) {
			ctrlCmdAck := res.(*msgs.CtrlCmdAck)
			fmt.Println(ctrlCmdAck.Result)
			ch <- fmt.Sprintf("%v", ctrlCmdAck.Result)
		})
	})

	json := <-ch
	return fmt.Sprintf("[%s]", json)
}

func (s *WebService) getNodesStatus() string {
	ch := make(chan string, 1)
	time.AfterFunc(time.Second*3, func() {
		ch <- "timeout"
	})

	cluster := app.Node.GetCluster()
	members := cluster.GetMembers()

	size := len(members)

	var results []nodectrl.NodeStatus

	s.Post(func() {
		for _, member := range members {
			id := member.Id
			if strings.Contains(id, app.MasterId) {
				ch <- app.MasterId
				continue
			}
			host := member.Host
			port := member.Port
			pid := actor.NewPID(fmt.Sprintf("%v:%v", host, port), define.NodeAdmin)

			s.RequestEx(pid, "ctrl.cmd", &msgs.CtrlCmd{
				Cmd: "web_nodes",
			}, func(err error, res any) {
				ctrlCmdAck := res.(*msgs.CtrlCmdAck)
				ch <- fmt.Sprintf("%v", ctrlCmdAck.Result)
			})
		}
	})
	for size > 0 {
		json := <-ch
		size--
		if json == app.MasterId {
			continue
		}
		if json == "timeout" {
			continue
		}
		var status nodectrl.NodeStatus
		jsonutils.Unmarshal([]byte(json), &status)
		results = append(results, status)
	}
	if len(results) == 0 {
		return "[]"
	}
	sort.Slice(results, func(i, j int) bool {
		time1, _ := strconv.Atoi(results[i].Time)
		time2, _ := strconv.Atoi(results[j].Time)
		return time1 > time2
	})
	marshal := jsonutils.Marshal(results)
	return marshal
}

func (s *WebService) startHttpServer(port int) {

	f := flamego.Classic()

	f.Use(flamego.Static(flamego.StaticOptions{
		Directory: "webservice/vue-admin-webapp/dist",
	}))

	f.Any("/login", login)
	f.Any("/getInfo", getInfo)
	f.Any("/getRoles", getRoles)
	f.Any("/getCardsData", getCardsData)
	f.Any("/getTableData", getTableData)
	f.Any("/getLineData", getLineData)
	f.Any("/getPageData1", getPageData1)
	f.Any("/getPageData2", getPageData2)
	f.Any("/getBarData", getBarData)

	f.Any("/api/nodes", s.getNodesStatus)
	f.Any("/api/retire", s.retire)
	f.Any("/api/exit", s.exit)

	if port == 0 {
		port = HttpPort
	}
	f.Run(port)
}
