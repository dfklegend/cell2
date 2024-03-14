package config

import (
	"errors"
	"fmt"
)

//BTProjectCfg 工程json类型
type BTProjectCfg struct {
	Version     string            `json:"version"`
	Scope       string            `json:"scope"`
	Select      string            `json:"selectedTree"`
	ID          string            `json:"id"`
	Trees       []BTTreeCfg       `json:"trees"`
	CustomNodes []BTCustomNodeCfg `json:"custom_nodes"`
}

func (b *BTProjectCfg) GetDefaultSelectTreeCfg() (*BTTreeCfg, error) {
	for _, tree := range b.Trees {
		if tree.ID == b.Select {
			return &tree, nil
		}
	}
	return nil, errors.New("default select tree not exist")
}

func (b *BTProjectCfg) FindTreeCfg(id string) (*BTTreeCfg, error) {
	for _, tree := range b.Trees {
		if tree.ID == id {
			return &tree, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Can't find tree {%v}", id))
}
