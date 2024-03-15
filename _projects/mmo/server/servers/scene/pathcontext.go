package sceneservice

import (
	"mmo/servers/scene/blockgraph"
)

type Point struct {
	X int
	Y int
}

type PathContext struct {
	graph          *blockgraph.Graph
	tempSkipPoints []Point
}

func newPathContext(graph *blockgraph.Graph) *PathContext {
	return &PathContext{
		graph:          graph,
		tempSkipPoints: make([]Point, 10),
	}
}

func (c *PathContext) UpdateInfo() {
	//
}

func (c *PathContext) IsInBlock(x, y int) bool {
	if len(c.tempSkipPoints) > 0 {
		if c.isInSkipPoint(x, y) {
			return false
		}
	}
	return c.graph.IsInBlock(x, y)
}

func (c *PathContext) IsNearEnough(x, y int) bool {
	return false
}

func (c *PathContext) AddTempSkipPoint(x, y int) {
	c.tempSkipPoints = append(c.tempSkipPoints, Point{x, y})
}

func (c *PathContext) ClearTempSkipPoints() {
	c.tempSkipPoints = make([]Point, 10)
}

func (c *PathContext) isInSkipPoint(x, y int) bool {
	for _, v := range c.tempSkipPoints {
		if x == v.X && y == v.Y {
			return true
		}
	}
	return false
}
