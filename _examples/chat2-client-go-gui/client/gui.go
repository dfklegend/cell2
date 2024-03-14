package client

import (
	"log"

	g "github.com/AllenDang/giu"
	"github.com/lucasb-eyer/go-colorful"
)

//GUIStartup 图形界面启动
func GUIStartup(host string) {

	cfg.UserName = "tong"
	cfg.Host = host

	w := g.NewMasterWindow("ChatApp", 800, 600, 0)
	w.SetCloseCallback(func() bool {
		log.Println("ChatApp Window Close")
		return true
	})
	w.Run(guiLoop)
}

func guiLoop() {

	w := g.SingleWindow()

	header := g.Row(
		g.InputText(&cfg.Host),
		g.Button("连接").OnClick(func() {
			log.Println("连接")

			//连接 gateway
			cfg.newGateway = true
			//go start(cfg.Host)
			start(cfg.Host)
		}),
		g.Button("断开连接").OnClick(func() {
			log.Println("断开连接")
			if cfg.Client == nil {
				return
			}
			cfg.Client.close()
		}))

	msgArea := g.Child().Size(0, 500).Layout(func() []g.Widget {
		var contents []g.Widget
		for _, history := range cfg.ChatHistory {
			var textWidget g.Widget
			if history.Name == cfg.UserName {
				green, _ := colorful.Hex("#7CFC00")
				textWidget = g.Style().SetColor(g.StyleColorText, green).To(g.BulletTextf("[%v][%s]: %s", history.Time.Format("15:04:05"), history.Name, history.Msg))
			} else {
				textWidget = g.BulletTextf("[%v][%s]: %s", history.Time.Format("15:04:05"), history.Name, history.Msg)
			}
			contents = append(contents, textWidget)
		}
		return contents
	}()...)

	bottom := g.Row(
		g.InputText(&cfg.UserName).Size(100),
		g.Row(
			g.InputText(&cfg.ChatMsg).Hint("开始聊天"),
			g.Event().OnKeyPressed(g.KeyEnter, func() {
				if len(cfg.ChatMsg) == 0 {
					return
				}
				//发送聊天信息
				cfg.Client.sendChatMsg(cfg.ChatMsg)
				//清除发送框
				cfg.ChatMsg = ""
			}),
		),
		g.Button("Send").OnClick(func() {
			if len(cfg.ChatMsg) == 0 {
				return
			}
			//发送聊天信息
			cfg.Client.sendChatMsg(cfg.ChatMsg)
			//清除发送框
			cfg.ChatMsg = ""
		}),
		g.Button("Clear").OnClick(func() {
			cfg.ChatHistory = []ChatMsg{}
		}),
	)

	w.Layout(
		g.Align(g.AlignCenter).To(g.Labelf("%s %s", cfg.Host, cfg.Room)),
		header,
		msgArea,
		bottom,
		//g.PrepareMsgbox(),
	)
}
