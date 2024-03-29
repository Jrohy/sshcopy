package main

import (
	"fmt"
	"github.com/fatih/color"
	expect "github.com/google/goexpect"
	"github.com/howeyc/gopass"
	"google.golang.org/grpc/codes"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

var (
	timeout = 10 * time.Minute
)

// Server 结构体
type Server struct {
	ip   string
	port int
	user string
	pass string
}

func (server *Server) sshTest() bool {
	connect := true
	var timeCostPoint *string
	ptr := &timeCostPoint
	defer TimeCostPTR(time.Now(), ptr)
	idPath := filepath.Join(os.Getenv("HOME"), ".ssh", "id_ed25519")
	if !IsExists(idPath) {
		idPath = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	}
	key, err := os.ReadFile(idPath)
	if err != nil {
		logger.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		logger.Fatalf("unable to parse private key: %v", err)
	}

	knowHostsPath := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	if !IsExists(knowHostsPath) {
		knowHostsFile, err := os.Create(knowHostsPath)
		if err != nil {
			logger.Fatal(err)
		}
		defer knowHostsFile.Close()
	}

	hostKeyCallback, err := kh.New(knowHostsPath)
	if err != nil {
		logger.Fatal("could not create hostkey callback function: ", err)
	}

	config := &ssh.ClientConfig{
		User: server.user,
		Auth: []ssh.AuthMethod{
			// Add in password check here for moar security.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
	}
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server.ip, server.port), config)
	if err != nil {
		connect = false
	} else {
		client.Close()
	}

	var colorConnect string
	if connect {
		colorConnect = color.GreenString("true")
	} else {
		colorConnect = color.RedString("false")
	}
	result := fmt.Sprintf("ssh连接性测试: '%s@%s -p %d' %s", server.user, server.ip, server.port, colorConnect)
	timeCostPoint = &result
	return connect
}

func (server *Server) copySSHID() {
	defer TimeCost(time.Now(), "ssh copy-id")
	// 调试用
	//e, _, err := expect.SpawnWithArgs([]string{"ssh-copy-id", fmt.Sprintf("%s@%s", server.user, server.ip), "-p", strconv.Itoa(server.port)}, timeout, expect.Verbose(true), expect.DebugCheck(log.New(os.Stdout,"Info:",log.Ldate | log.Ltime | log.Lshortfile)))
	e, _, err := expect.Spawn(fmt.Sprintf("ssh-copy-id %s@%s -p %d", server.user, server.ip, server.port), timeout)
	if err != nil {
		logger.Fatal(err)
	}
	defer e.Close()

	var (
		retryCount = 0
		exit       = false
	)
	caser := []expect.Caser{
		&expect.BCase{R: "password", T: func() (tag expect.Tag, status *expect.Status) {
			password := server.pass
			if retryCount == 3 {
				exit = true
			} else {
				if password == "" {
					if retryCount > 0 {
						logger.Printf("%s, please try again\n\n", color.RedString("Permission denied"))
					}
					tempPass, _ := gopass.GetPasswdPrompt(fmt.Sprintf("请输入'%s@%s -p %s'的密码: ", color.CyanString(server.user), color.CyanString(server.ip), color.CyanString(strconv.Itoa(server.port))), true, os.Stdin, os.Stdout)
					password = string(tempPass)
				}
				_ = e.Send(password + "\n")
				retryCount++
			}
			return expect.OKTag, expect.NewStatus(codes.OK, "")
		}},
		&expect.BCase{R: "yes/no", S: "yes\n"},
	}

	for {
		if exit {
			logger.Println("已经达到3次输错密码! 请重新运行脚本进行免密操作")
			break
		}
		if server.pass != "" && retryCount > 1 {
			logger.Println(fmt.Sprintf("'%s@%s -p %s': %s", color.CyanString(server.user), color.CyanString(server.ip), color.CyanString(strconv.Itoa(server.port)), color.RedString("Permission denied")))
			break
		}
		if output, _, _, err := e.ExpectSwitchCase(caser, timeout); err != nil {
			if strings.Contains(output, "known_hosts") {
				cmd := fmt.Sprintf("sed -i '/%s/d' %s", server.ip, filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
				if _, execError := exec.Command("bash", "-c", cmd).Output(); execError != nil {
					logger.Fatal(execError)
				}
				e, _, _ = expect.Spawn(fmt.Sprintf("ssh-copy-id %s@%s -p %d", server.user, server.ip, server.port), timeout)
				continue
			}
			if !Simple {
				logger.Printf("\n" + output)
			}
			if strings.Contains(output, "added") {
				logger.Println(color.GreenString("成功拷贝密钥!"))
			} else if strings.Contains(output, "exist") {
				logger.Println(color.YellowString("密钥已存在!"))
			} else {
				logger.Println(color.RedString("拷贝密钥失败!"))
			}
			break
		}
	}
}
