package route

import (
	"fmt"
	"log"
	"testing"
	//"reflect"

	"github.com/stretchr/testify/assert"

	"github.com/dfklegend/cell2/utils/logger"
)

func TestErrorParam(t *testing.T) {
	s := TheRouteService
	fmt.Println("TestErrorParam")
	assert.Equal(t, BadRouteParam, s.Route("test", 1))
	fmt.Println("--")
}

func TestMapParam(t *testing.T) {
	s := NewRouteService()

	fmt.Println("TestMapParam")
	m := make(map[string]interface{})
	m["haha"] = 1
	log.Printf("route result:%v\n", s.Route("test", m))
	assert.Equal(t, MissRouteFunc, s.Route("test", m))

	s.Register("test", RouteF)
	assert.Equal(t, "", s.Route("test", m))
	fmt.Println("--")
}

func getServers(serverType string) ([]string, error) {
	all := map[string][]string{
		"test": {
			"s1", "s2",
		},
	}
	return all[serverType], nil
}

func RouteF(serverType string, p IRouteParam) string {
	var id = p.Get("id", 0).(int)
	servers, err := getServers(serverType)

	if servers == nil || len(servers) == 0 || err != nil {
		logger.Log.Errorf("can not find serverType:%v\n", serverType)
		return "null"
	}

	want := fmt.Sprintf("s%v", id)

	for _, v := range servers {
		if want == v {
			return want
		}
	}

	log.Printf("can not find server:%v\n", want)
	return ""
}

func TestRoute(t *testing.T) {
	fmt.Println("TestRoute")
	s := NewRouteService()

	s.Register("test", RouteF)
	m := make(map[string]interface{})
	m["id"] = 1
	assert.Equal(t, "s1", s.Route("test", m))

	m["id"] = 3
	assert.Equal(t, "", s.Route("test", m))

	assert.Equal(t, MissRouteFunc, s.Route("test1", m))

	// string
	assert.Equal(t, "s5", s.Route("test", "s5"))
	fmt.Println("--")
}
