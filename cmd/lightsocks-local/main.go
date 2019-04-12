package main

import (
	"fmt"
	"github.com/c-dafan/lightsocks"
	"github.com/c-dafan/lightsocks/cmd"
	"log"
	"net"
)

const (
	DefaultListenAddr = ":7448"
)

var version = "master"

func main() {
	log.SetFlags(log.Lshortfile)

	// 默认配置
	config := &cmd.Config{
		ListenAddr: DefaultListenAddr,
		Password:   "YEobfrT6VP6u252mLK3LEotXJzMVWFk8fK/ZAOXCbQkCXSnwtjuewIFiiNGY4oMwHt9w5ytb9ORR65AYN0zKEz69HCJoKs6b3VYyDqyj9+/Fls3u7NIBqbB/VcMhdUTXk6FHIGU4Xy2OisyElI31U/EZY5WCjPjyv97mvhq4uwQfPxRkUqCJNEm8mrMmyPNsq1w1Oe23z+EjQA/jNjrJhmH/4OgKl9hCQV4vRtTEuXqfx0uqRTF9nEOPxmuxSG6SonS6CME9pXaZ+wP8snkkhbUo3P0uatp46QbWZ01xkacXctMLbw1Qc08RaYcH0OqkWvkW9neAJR0Me2bVBahOEA==",
		RemoteAddr: "127.0.0.1:39126",
	}
	// config.ReadConfig("F:\\Go_workplace\\src\\github.com\\gwuhaolin\\lightsocks\\cmd\\lightsocks-local\\con.json")
	config.SaveConfig("F:\\Go_workplace\\src\\github.com\\gwuhaolin\\lightsocks\\cmd\\lightsocks-local\\con.json")

	// 启动 local 端并监听
	lsLocal, err := lightsocks.NewLsLocal(config.Password, config.ListenAddr, config.RemoteAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(lsLocal.Listen(func(listenAddr net.Addr) {
		log.Println("使用配置：", fmt.Sprintf(`
本地监听地址 listen：
%s
远程服务地址 remote：
%s
密码 password：
%s
	`, listenAddr, config.RemoteAddr, config.Password))
		log.Printf("lightsocks-local:%s 启动成功 监听在 %s\n", version, listenAddr.String())
	}))
}
