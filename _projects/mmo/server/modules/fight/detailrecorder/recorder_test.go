package detailrecorder

import (
	"log"
	"testing"

	"mmo/modules/fight/utils"
)

func Test_Normal(t *testing.T) {
	r := &Recorder{}
	r.timeProvider = utils.NewTestTimeProvider()

	r.log(" test")
	r.log(" test1")

	r.Visit(func(info string) {
		log.Println(info)
	})
}
