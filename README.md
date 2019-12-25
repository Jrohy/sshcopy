# sshcopy
[![Go Report Card](https://goreportcard.com/badge/github.com/Jrohy/sshcopy)](https://goreportcard.com/report/github.com/Jrohy/sshcopy)
[![Downloads](https://img.shields.io/github/downloads/Jrohy/sshcopy/total.svg)](https://img.shields.io/github/downloads/Jrohy/sshcopy/total.svg)

自动生成密钥和拷贝密钥到远程服务器(ssh-copy-id), 支持并发批量设置服务器免密

## 运行方式
### 1. 命令行参数
```bash
./sshcopy -ip [ip] -user [user] -port [port] -pass [pass] [-h|--help]
所有参数支持多个参数传参, 空格隔开, 例如 -ip "ip1 ip2" -port "port1 port2"
    -h, --help           查看帮助
    -ip                  server ip, 不传脚本进入交互输入模式
    -user                server user, 多个user时和ip按顺序匹配, user数量不足用最后一个, 不传默认所有ip user为root
    -port                server port, 多个port时和ip按顺序匹配, port数量不足用最后一个, 不传默认所有ip port为22
    -pass                server password, 多个password时和ip按顺序匹配, pass数量不足用最后一个, 不传脚本会提示输入服务器密码
```

### 2. 交互输入模式
直接运行`./sshcopy`即可, 脚本会提示输入要进行免密的服务器ip、端口、用户等信息的
