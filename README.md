# 魔改frp

以下都是在frp v0.38.0基础上更改
## 介绍
1. x17的流量特征的问题
   在2021/10/25时发布的新版本，可以在frpc配置disable_custom_tls_first_byte参数，禁用这个字节
2. 增加子命令loadini，可以从命令行导入base64的配置文件
3. 增加socks和pf命令

socks用于开启socks5代理
使用tcp协议
未开启tls
开启了use_encryption和use_compression
未指定连接名时使用随机生成的8个16进制数

pf用于端口转发
使用tcp协议
未开启tls
开启了use_encryption和use_compression
未指定连接名时使用随机生成的8个16进制数


## 改动
只改了frpc，
1. cmd/frpc/sub下增加loadini.go、socks5.go、portforward.go
2. pkg/config/parse.go中添加ParseClientConfig1函数，用于从字符串 解析配置文件
3. 修改cmd/frpc/sub/root.go，注释了protocol命令

## 不足
影响了frp原本功能，删掉了protocol命令
因为子命令参数password和protrol的shorthand都是p，冲突了
平时使用frp都是通过配置文件使用，所以就直接删掉了