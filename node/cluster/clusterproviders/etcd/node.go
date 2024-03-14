package etcd

import (
	"encoding/json"

	"github.com/dfklegend/cell2/node/cluster"
)

type Node struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Host     string            `json:"host"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Services []string          `json:"services"`
	Meta     map[string]string `json:"-"`
	Alive    bool              `json:"alive"`
	State    int               `json:"state"`
}

func NewNode(name, host string, port int, services []string) *Node {
	return &Node{
		ID:       name,
		Name:     name,
		Address:  host,
		Host:     host,
		Port:     port,
		Services: services,
		Meta:     map[string]string{},
		Alive:    true,
	}
}

func NewNodeFromBytes(data []byte) (*Node, error) {
	n := Node{}
	if err := json.Unmarshal(data, &n); err != nil {
		return nil, err
	}
	return &n, nil
}

func (n *Node) GetAddress() (host string, port int) {
	host = n.Host
	port = n.Port
	if host == "" {
		host = n.Address
	}
	return
}

func (n *Node) Equal(other *Node) bool {
	if n == nil || other == nil {
		return false
	}
	if n == other {
		return true
	}
	return n.ID == other.ID
}

func (n *Node) GetMeta(name string) (string, bool) {
	if n.Meta == nil {
		return "", false
	}
	val, ok := n.Meta[name]
	return val, ok
}

func (n *Node) MemberStatus() *cluster.Member {
	host, port := n.GetAddress()
	services := n.Services
	if services == nil {
		services = []string{}
	}
	return &cluster.Member{
		Id:       n.ID,
		Host:     host,
		Port:     int32(port),
		Services: services,
		State:    n.State,
	}
}

func (n *Node) SetMeta(name string, val string) {
	if n.Meta == nil {
		n.Meta = map[string]string{}
	}
	n.Meta[name] = val
}

func (n *Node) Serialize() ([]byte, error) {
	data, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (n *Node) Deserialize(data []byte) error {
	return json.Unmarshal(data, n)
}

func (n *Node) IsAlive() bool {
	return n.Alive
}

func (n *Node) SetAlive(alive bool) {
	n.Alive = alive
}

func (n *Node) SetState(state int) {
	n.State = state
}
