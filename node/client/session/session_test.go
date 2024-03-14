package session

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestFrontSession(serverId string, netId uint32) *FrontSession {
	f := &FrontSession{
		Session: nil,
		Data:    NewSessionData(),
	}

	f.Set(KeyServerId, serverId)
	f.Set(KeyNetId, netId)
	return f
}

func Test_BackSession(t *testing.T) {
	//var is IServerSession

	f := newTestFrontSession("test", 1)
	f.Set("k1", 1)
	f.Set("k2", "ffff")

	log.Printf("%v %v\n", f.Get("k1", 0), f.Get("k2", ""))

	str := f.ToJson()
	b := NewBackSession(nil, "server", 0, "")
	b.FromJson(str)
	//is = b
	log.Printf("%+v %v\n", b, b.ToJson())
	b.Set("k3", 2)
	b.Set("k4", 2)
	log.Printf("%+v %v\n", b, b.ToJson())

	assert.Equal(t, float64(1), b.Get("k1", 0))
	assert.Equal(t, 2, b.Get("k4", 0))
}

func Test_UpdateSessionData(t *testing.T) {
	f := newTestFrontSession("test", 1)
	f.Set("f1", 1)
	f.Set("f2", "ffff")

	log.Printf("%v %v\n", f.Get("k1", 0), f.Get("k2", ""))

	b := NewBackSession(nil, "server", 0, "")
	b.Set("b1", 1)
	b.Set("f1", 100)

	f.Data.UpdateFromJson([]byte(b.ToJson()))

	assert.Equal(t, float64(1), f.Get("b1", 0))
	assert.Equal(t, float64(100), f.Get("f1", 0))
}
