package cmd

import (
	"bufio"
	"fmt"
	"os"
)

var running = true

func getInput() string {
	in := bufio.NewReader(os.Stdin)
	str, _, err := in.ReadLine()
	if err != nil {
		return err.Error()
	}
	return string(str)
}

func processInput() {
	for running {
		doOnce()
	}
}

func doOnce() {
	fmt.Print("please input:")
	s := getInput()
	fmt.Println(s)
	// 传到cmdMgr中
	DispatchCmd(s)
}

//	简单的提供一个console

func StartConsoleCmd() {
	running = true
	go processInput()
}

func LoopConsoleCmd() {
	running = true
	processInput()
}

func StopConsoleCmd() {
	running = false
}
