# sshcopy
[![Go Report Card](https://goreportcard.com/badge/github.com/Jrohy/sshcopy)](https://goreportcard.com/report/github.com/Jrohy/sshcopy)
[![Downloads](https://img.shields.io/github/downloads/Jrohy/sshcopy/total.svg)](https://img.shields.io/github/downloads/Jrohy/sshcopy/total.svg)

自动生成密钥和拷贝密钥到远程服务器(ssh-copy-id), 支持并发批量设置服务器免密

## 运行方式
### 1. 命令行参数
```bash
./sshcopy -ip [ip] -user [user] -port [port] -pass [pass] [-h|--help]
    -ip: 不传则脚本进入交互输入模式, 等于什么参数都没传; -user: 不传则默认所有ip user为root; -port: 不传则默认所有ip port为22; -pass: 此选项仅供被其他脚本调用时传参, 手动运行脚本时建议不传, 脚本会提示输入密码
    -h, --help           查看帮助
    -ip                  server ip, 多个ip空格隔开, 例如: -ip "192.168.37.193 192.168.37.100"
    -user                server user, 多个user空格隔开, 和ip按顺序匹配, 匹配数不足用最后一个, 例如: -user "user1 user2"
    -port                server port, 多个port空格隔开, 和ip按顺序匹配, 匹配数不足用最后一个, 例如: -port "port1 port2"
    -pass                server password, 多个password空格隔开, 和ip按顺序匹配, 匹配数不足用最后一个, 例如: -pass "pass1 pass2"
```

### 2. 交互输入模式
直接运行`./sshcopy`即可, 脚本会提示输入要进行免密的服务器ip、端口、用户等信息的
