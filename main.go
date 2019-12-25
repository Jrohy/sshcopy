package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	wg             sync.WaitGroup
	logger         *log.Logger
	ipStr, passStr string
	portStr        = "22"
	userStr        = "root"
)

func usage() {
	fmt.Println("Usage: " + os.Args[0] + " -ip [ip] -user [user] -port [port] -pass [pass] [-h|--help]")
	flag.PrintDefaults()
	os.Exit(0)
}

func init() {
	logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	flag.StringVar(&ipStr, "ip", "", `server ip, 多个ip空格隔开, 例如: -ip "192.168.37.193 192.168.37.100", 不传则脚本进入交互输入模式(等于什么参数都没传)`)
	flag.StringVar(&userStr, "user", "root", `server user, 多个user空格隔开, 和ip按顺序匹配, 匹配数不足用最后一个, 例如: -user "user1 user2", 不传则默认所有ip user为root`)
	flag.StringVar(&portStr, "port", "22", `server port, 多个port空格隔开, 和ip按顺序匹配, 匹配数不足用最后一个, 例如: -port "port1 port2", 不传则默认所有ip port为22`)
	flag.StringVar(&passStr, "pass", "", `server password, 多个password空格隔开, 和ip按顺序匹配, 匹配数不足用最后一个, 例如: -pass "pass1 pass2", 不传脚本会提示输入服务器密码`)
	flag.Usage = usage
}

func obtainParameter() {
	flag.Parse()
	if ipStr == "" {
		inputReader := bufio.NewReader(os.Stdin)
		fmt.Println("请输入要进行免密的服务器ip, 多个ip空格隔开: ")
		ipStr, _ = inputReader.ReadString('\n')
		temp := ""
		fmt.Printf("默认连接端口: %s, 正确直接回车, 否则输入自定义端口(多个端口空格隔开, 和ip按顺序匹配, 匹配数不足用最后一个): \n", color.CyanString("22"))
		temp, _ = inputReader.ReadString('\n')
		if strings.Replace(temp, "\n", "", -1) != "" {
			portStr = temp
		}
		fmt.Printf("默认连接用户: %s, 正确直接回车, 否则输入自定义用户(多个用户空格隔开, 和ip按顺序匹配, 匹配数不足用最后一个): \n", color.CyanString("root"))
		temp, _ = inputReader.ReadString('\n')
		if strings.Replace(temp, "\n", "", -1) != "" {
			userStr = temp
		}
	}
}

func verifyIP(inputSlice []string) []string {
	result := make([]string, 0, len(inputSlice))
	for _, ip := range inputSlice {
		if len(ip) == 0 {
			continue
		}
		if !CheckIP(ip) {
			logger.Printf("%s不是常规ip, 直接跳过\n", color.YellowString(ip))
			continue
		}
		result = append(result, ip)
	}
	return result
}

func filterValue(index int, totalSize int, key []string) string {
	result := ""
	if len(key) >= totalSize {
		result = key[index]
	} else {
		if index < len(key) {
			result = key[index]
		} else {
			result = key[len(key)-1]
		}
	}
	return result
}

func main() {
	obtainParameter()
	GenerateRsa()

	if ipStr == "" {
		logger.Fatal("必须输入要免密的服务器ip!")
	}
	var serverSlice []Server
	ipSlice := strings.Split(strings.Join(strings.Fields(ipStr), " "), " ")
	userSlice := strings.Split(strings.Join(strings.Fields(userStr), " "), " ")
	portSlice := strings.Split(strings.Join(strings.Fields(portStr), " "), " ")
	passSlice := strings.Split(strings.Join(strings.Fields(passStr), " "), " ")
	qualifiedIP := verifyIP(ipSlice)
	totalSize := len(qualifiedIP)

	for ipIndex, ip := range qualifiedIP {
		wg.Add(1)

		go func(ip string, user string, port string, pass string) {
			defer wg.Done()
			portInt, err := strconv.Atoi(port)
			if err != nil {
				logger.Printf("%s 转换端口出错! 直接跳过\n", color.YellowString(port))
				return
			}
			server := Server{ip: ip, port: portInt, user: user, pass: pass}
			isConnect := server.sshTest()
			if isConnect {
				logger.Printf("%s服务器已经设置为免密!\n", color.MagentaString(server.ip))
			} else {
				serverSlice = append(serverSlice, server)
			}
		}(ip, filterValue(ipIndex, totalSize, userSlice), filterValue(ipIndex, totalSize, portSlice), filterValue(ipIndex, totalSize, passSlice))
	}
	wg.Wait()

	for _, s := range serverSlice {
		s.copySSHID()
	}
}
