package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	console "github.com/asynkron/goconsole"

	"github.com/dfklegend/cell2/utils/golua"
	"github.com/dfklegend/cell2/utils/golua/module/golog"
	"mlrs/b3"
	"mlrs/b3/core"
	"mlrs/b3/factory"
	"mlrs/example"
)

const packageName = "example"

func main() {
	fmt.Println("多管火箭发射系统 multiple launch rocket system")
	example.Init()

	golua.InitLuaPathAndCompile("lua", true)

	home, _ := os.UserHomeDir()

	filePath := home + "/com.dfklegend.jgzx.robot/conf.properties"

	file, _ := os.Open(filePath)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	tpl := loadCodeTpl()

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "btFilePath=") {
			b3FilePath := strings.TrimPrefix(line, "btFilePath=")
			println(b3FilePath)

			rawProjectCfg, ok := b3.LoadRawProjectCfg(b3FilePath)
			if ok {

				projectName := rawProjectCfg.Name
				projectPath := rawProjectCfg.Path
				projectData := rawProjectCfg.Data

				fmt.Printf("projectName = %v\n", projectName)
				fmt.Printf("projectPath = %v\n", projectPath)
				//fmt.Printf("projectData = %v\n", projectData)

				//自定义节点注册
				maps := factory.RegisterMap()
				quit := false
				for _, customNode := range projectData.CustomNodes {
					if !maps.CheckElem(customNode.Name) {
						quit = true
						log.Printf("未注册节点类型{%v}", customNode.Name)
						genCode(tpl, customNode)
					}
				}
				if quit {
					os.Exit(0)
					return
				}

				//载入
				tree := factory.CreateBevTreeFromProjectData(&projectData)
				tree.Print()

				//输入板
				board := core.NewBlackboard()

				luaEngine := golua.NewLuaEngine()
				//加载第三方库
				luaEngine.LoadGopherLuaLibs()
				//加载自定义模块
				//luaEngine.LoadModule("go_log", golog.Loader) // 1
				golog.Preload(luaEngine.L) // 2

				board.SetMem("lua", luaEngine)
				luaEngine.DoLuaFile("init.lua")
				//循环每一帧
				for i := 0; i < 50000; i++ {
					tree.Tick(i, board)
					time.Sleep(10)
				}
				luaEngine.Close()
			}
		}
	}

	console.ReadLine()
}
