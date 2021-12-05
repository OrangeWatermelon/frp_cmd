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

## 使用方式
从base64编码导入配置
```
frpc loadini -i <base64ini>

help:
-i, --ini string           base64 ini
```
启用端口转发
```
frpc fp -s 1.1.1.1:1234 -r 1234 -l 8080 [--lip 192.168.1.3] [-t f86bc7ff68aff3ad] [-n zz]

help:
-r, --rp string            remote port
-l, --lp string            local port
--lip string           local ip (default "127.0.0.1")
-t, --token string         token
-n, --name string          name (default 随机生成)
```
启用socks代理
```
frpc socks -s 1.1.1.1:1234 -r 1234 [-t f86bc7ff68aff3ad] [-n zz] [-u z] [-p z]

help:
-s, --server_addr string   frp server's address (default "127.0.0.1:7000")
-r, --rp string            remote port
-t, --token string         token
-n, --name string          name (default 随机生成)
-u, --user string          user
-p, --pwd string           password
```

## 改动
只改了frpc，
1. cmd/frpc/sub下增加loadini.go、socks5.go、portforward.go
2. pkg/config/parse.go中添加ParseClientConfig1函数，用于从字符串 解析配置文件
3. 修改cmd/frpc/sub/root.go，注释了protocol命令

## 不足
影响了frp原本功能，删掉了protocol命令
因为子命令参数password和protrol的shorthand都是p，冲突了
平时使用frp都是通过配置文件使用，所以就直接删掉了