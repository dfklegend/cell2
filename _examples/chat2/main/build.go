package main

import (
	"log"
)

/*
	build 包
	编译时，使用-ldflags来注入数据
*/

var (
	Version   = "not set"
	Time      = "not set"
	GoVersion = "not set"
)

func DumpInfo() {
	log.Println("---- build info ----")
	log.Printf("Version = %x\r\n", Version)
	log.Printf("Time = %x\r\n", Time)
	log.Printf("GoVersion = %x\r\n", GoVersion)
	log.Println("---- build info end ----")
}
