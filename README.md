# [Lightsocks](https://github.com/gwuhaolin/lightsocks)
一个轻量级网络混淆代理，基于 SOCKS5 协议，可用来代替 Shadowsocks。

- 只专注于混淆，用最简单高效的混淆算法达到目的；
- 不会放大传输流量，传输流量更少更快，占用资源更少（内存占用1M左右，CPU 占用 0.1% 左右）；
- 纯 Golang 编写，跨平台，对Golang感兴趣？请看[Golang 中文学习资料汇总](http://go.wuhaolin.cn/)；

想了解 Lightsocks 的实现原理？请阅读文章：[你也能写个 Shadowsocks](https://github.com/gwuhaolin/blog/issues/12)。 
[源项目请见](https://github.com/gwuhaolin/lightsocks)

## 修改

* 增加日志输出
* 增加指定配置文件

## 使用方法

[详细说明](http://www.zsdf.shop/04/12/Lightsocks%E7%8E%AF%E5%A2%83%E9%85%8D%E7%BD%AE/)

### 启动 lightsocks-server

```bash
go run lightsocks-server
```
在配置文件server.json中可以看到密码，端口信息。

### 启动 lightsocks-local

修改local.json中配置
```json
{
	"listen": ":7448",
	"remote": "",
	"password": ""
}
```
remote改为运行lightsocks-server主机的ip+server.json的listen,

password改为server.json的password。

```bash
go run lightsocks-local
```

### 配置浏览器

[参考](https://github.com/gwuhaolin/lightsocks/wiki/%E6%90%AD%E9%85%8D-Chrome-%E4%BD%BF%E7%94%A8)