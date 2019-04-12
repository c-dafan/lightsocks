package main

import (
	"fmt"
	"github.com/c-dafan/lightsocks"
	"github.com/c-dafan/lightsocks/cmd"
	"log"
	"net"
)

var version = "master"

func main() {
	log.SetFlags(log.Lshortfile)

	// 服务端监听端口随机生成
	//port, err := freeport.GetFreePort()
	//if err != nil {
	//	// 随机端口失败就采用 7448
	//	port = 7448
	//}
	// 默认配置
	config := &cmd.Config{
		ListenAddr: "127.0.0.1:39126",
		// 密码随机生成
		Password: "YEobfrT6VP6u252mLK3LEotXJzMVWFk8fK/ZAOXCbQkCXSnwtjuewIFiiNGY4oMwHt9w5ytb9ORR65AYN0zKEz69HCJoKs6b3VYyDqyj9+/Fls3u7NIBqbB/VcMhdUTXk6FHIGU4Xy2OisyElI31U/EZY5WCjPjyv97mvhq4uwQfPxRkUqCJNEm8mrMmyPNsq1w1Oe23z+EjQA/jNjrJhmH/4OgKl9hCQV4vRtTEuXqfx0uqRTF9nEOPxmuxSG6SonS6CME9pXaZ+wP8snkkhbUo3P0uatp46QbWZ01xkacXctMLbw1Qc08RaYcH0OqkWvkW9neAJR0Me2bVBahOEA==",
	}
	// config.ReadConfig()
	config.SaveConfig(nil)

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
