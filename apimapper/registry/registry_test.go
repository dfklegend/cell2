package registry

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/utils/serialize/proto/msgs"
)

type EntryProto struct {
	api.APIEntry
}

func (e *EntryProto) Join(d api.IContext, msg *msgs.TestHello, cbFunc apientry.HandlerCBFunc) {
	log.Printf("in EntryProto.join\n")

	apientry.CheckInvokeCBFunc(cbFunc, nil, &msgs.TestHello{
		I: msg.I * 2,
	})
}

func MakeName(serviceType string, postfix string) string {
	return fmt.Sprintf("%v.%v", serviceType, postfix)
}

func TestMake(t *testing.T) {
	Registry.AddCollection(MakeName("chat", "remote")).
		Register(&EntryProto{},
			apientry.WithGroupName("entry"),
			apientry.WithNameFunc(strings.ToLower)).
		Register(&EntryProto{},
			apientry.WithGroupName("logic"),
			apientry.WithNameFunc(strings.ToLower)).
		Build()

	remote := Registry.GetCollection(MakeName("chat", "remote"))
	assert.Equal(t, true, remote != nil)
}
