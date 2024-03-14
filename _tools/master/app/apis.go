package master

import (
	"master/cmdservice"
	"master/consolecmd"
	"master/telnetcmd"
	"master/webservice"
)

func RegisterAllAPIEntries() {
	cmdservice.RegisterEntry()
	consolecmd.RegisterEntry()
	telnetcmd.RegisterEntry()
	webservice.RegisterEntry()
}
