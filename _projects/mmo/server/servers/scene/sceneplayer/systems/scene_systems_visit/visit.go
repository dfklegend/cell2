package scene_systems_visit

import (
	"mmo/servers/scene/sceneplayer/systems/baseinfo"
	"mmo/servers/scene/sceneplayer/systems/control"
	"mmo/servers/scene/sceneplayer/systems/daysign"
	"mmo/servers/scene/sceneplayer/systems/test1"
	"mmo/servers/scene/sceneplayer/systems/test2"
)

// Visit 注册所有使用的系统
func Visit() {
	test1.Visit()
	test2.Visit()
	baseinfo.Visit()
	control.Visit()
	daysign.Visit()
}
