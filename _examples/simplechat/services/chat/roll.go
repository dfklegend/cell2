package chat

import (
	"fmt"
	"log"
	"math/rand"
	"sort"

	"github.com/dfklegend/cell2/utils/common"
	mymsg "simplechat/messages"
)

/*
	收到/rollstart时，开启
	每个人可以roll一次
	都roll过了，就输出结果
	或者1分钟后自动结束
	可以主动/rollend结束
*/
type RollResult struct {
	ID    int32
	Score int32
}

// roll点控制器
type RollCtrl struct {
	service *Service
	room    *Room

	Running   bool
	Owner     int32
	TimeToEnd int64

	// results
	scores []*RollResult
}

func NewRoll() *RollCtrl {
	return &RollCtrl{
		scores: make([]*RollResult, 0),
	}
}

func (r *RollCtrl) Init(s *Service, room *Room) {
	r.service = s
	r.room = room
}

func (r *RollCtrl) Start(owner int32) bool {
	if r.Running {
		log.Println("has a roll already")
		return false
	}
	r.Running = true
	r.Owner = owner
	r.scores = make([]*RollResult, 0)
	r.TimeToEnd = common.NowMs() + 60*1000

	r.reportLine(fmt.Sprintf("%v 开始roll点", owner))
	return true
}

func (r *RollCtrl) ReqEnd(id int32) {
	if r.Owner != id {
		r.sendToMember(id, "only owner can end roll")
		return
	}
	r.End(true)
}

func (r *RollCtrl) End(needReport bool) {
	r.Running = false
	if len(r.scores) == 0 {
		r.reportLine("没人roll点")
		return
	}
	// sort
	sort.Slice(r.scores, func(i int, j int) bool {
		return r.scores[i].Score > r.scores[j].Score
	})
	// report
	if needReport {
		r.report()
	}
}

func (r *RollCtrl) report() {
	players := r.room.GetPlayers()
	r.reportLine("-------- roll result --------")

	r.reportLine(fmt.Sprintf("胜利者: %v", players.GetPlayerName(r.scores[0].ID)))
	r.reportLine("----")
	for _, v := range r.scores {
		r.reportLine(fmt.Sprintf("%v %v", players.GetPlayerName(v.ID), v.Score))
	}
	// 显示每一个结果
	r.reportLine("-------- result over --------")
}

func (r *RollCtrl) reportLine(str string) {
	msg := &mymsg.Chat{
		ID:   0,
		Name: "系统",
		Str:  str,
	}
	r.room.Broadcast(r.service, msg)
}

func (r *RollCtrl) sendToMember(id int32, str string) {
	msg := &mymsg.Chat{
		ID:   0,
		Name: "系统",
		Str:  str,
	}
	r.room.Send(r.service, id, msg)
}

func (r *RollCtrl) isRolled(id int32) bool {
	for _, v := range r.scores {
		if v.ID == id {
			return true
		}
	}
	return false
}

func (r *RollCtrl) PlayerRoll(id int32) {
	if !r.IsRunning() {
		r.sendToMember(id, "当前没有发起roll点，可以使用/rollbegin开始")
		return
	}

	if r.isRolled(id) {
		r.sendToMember(id, "本轮已经使用过roll")
		return
	}
	result := rand.Intn(100) + 1
	r.scores = append(r.scores, &RollResult{
		ID:    id,
		Score: int32(result),
	})

	r.reportLine(fmt.Sprintf("%v roll %v", id, result))

	if r.isAllRolled() {
		r.End(true)
	}
}

func (r *RollCtrl) isAllRolled() bool {
	for _, v := range r.room.Members {
		if !r.isRolled(v.ID) {
			return false
		}
	}
	return true
}

func (r *RollCtrl) IsRunning() bool {
	return r.Running
}

func (r *RollCtrl) Update() {
	if r.Running && common.NowMs() >= r.TimeToEnd {
		r.End(true)
	}
}
