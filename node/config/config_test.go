package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadServers(t *testing.T) {
	nodes := LoadNodes("../../testdata/config/")
	fmt.Printf("%v\n", nodes)
	for k, v := range nodes.Nodes {
		fmt.Printf("%v %v\n", k, v)
	}

	for k, v := range nodes.Services {
		fmt.Printf("%v %v\n", k, v)
	}

	assert.Equal(t, "gate", nodes.Services["gate-1"].Type)
}

func TestLoadCluster(t *testing.T) {
	info := LoadCluster("../../testdata/config/")
	fmt.Printf("%v\n", info)

	assert.Equal(t, "testcluster", info.Name)
}
