package build

import (
	"log"
)

/*
	build 包
	编译时，使用-ldflags来注入数据

	目前发现，如果想注入，必须在代码目录下执行
	比如
	xxx
		main
		other
	main包被放到main下执行，如果编译
	go build .\main -ldflags="-X 'main.Version=1.0.0'"  出错
	去到main目录下 执行
	go build -ldflags="-X 'main.Version=1.0.0'" 则正确
*/

var (
	Version   = "not set"
	Time      = "not set"
	GoVersion = "not set"
)

func DumpInfo(version string, time string, goVersion string) {
	log.Println("---- build info ----")
	log.Printf("Version = %v\r\n", version)
	log.Printf("Time = %v\r\n", time)
	log.Printf("GoVersion = %v\r\n", goVersion)
	log.Println("---- build info end ----")
}

func DumpBuildInfo() {
	DumpInfo(Version, Time, GoVersion)
}
