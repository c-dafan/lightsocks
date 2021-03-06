package main

import (
	"fmt"
	l4g "github.com/alecthomas/log4go"
	"github.com/c-dafan/lightsocks"
	"github.com/c-dafan/lightsocks/cmd"
	"github.com/phayes/freeport"
	"log"
	"net"
	"os"
)

var version = "master"

func main() {
	log.SetFlags(log.Lshortfile)
	err := os.Mkdir("logs", os.ModePerm)
	if err != nil {
		if !os.IsExist(err) {
			err = os.Mkdir("logs", os.ModePerm)
			if err != nil {
				return
			}
		}
	}
	// 服务端监听端口随机生成
	port, err := freeport.GetFreePort()
	if err != nil {
		// 随机端口失败就采用 7448
		port = 7448
	}
	// 默认配置
	config := &cmd.Config{
		ListenAddr: fmt.Sprintf(":%d", port),
		// 密码随机生成
		Password: lightsocks.RandPassword(),
	}
	config.ReadConfig("server.json")
	config.SaveConfig("server.json")
	l4g.LoadConfiguration("example.xml")
	defer l4g.Close()
	// 启动 server 端并监听
	lsServer, err := lightsocks.NewLsServer(config.Password, config.ListenAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(lsServer.Listen(func(listenAddr net.Addr) {
		log.Println("使用配置：", fmt.Sprintf(`
本地监听地址 listen：
%s
密码 password：
%s
	`, listenAddr, config.Password))
		log.Printf("lightsocks-server:%s 启动成功 监听在 %s\n", version, listenAddr.String())
	}))
}
