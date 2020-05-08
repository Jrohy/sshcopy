package main

import (
	expect "github.com/google/goexpect"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// TimeCost 耗时统计函数, 传普通字符串
func TimeCost(start time.Time, str string) {
	logger.Printf("%s  耗时: %v\n", str, time.Since(start))
}

// TimeCostPTR 耗时统计函数, 传指向指针的指针
func TimeCostPTR(start time.Time, strPtr **string) {
	logger.Printf("%s  耗时: %v\n", **strPtr, time.Since(start))
}

// CheckIP 检测ipv4地址的合法性
func CheckIP(ip string) bool {
	isOk, err := regexp.Match(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`, []byte(ip))
	if err != nil {
		logger.Fatalf("error: %v", err)
	}
	return isOk
}

// IsExists 检测指定路径文件或者文件夹是否存在
func IsExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// GenerateId 自动生成当前用户的密钥
func GenerateId() {
	var timeout = 10 * time.Minute
	rsaPath := filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	rsaPubPath := filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub")
	edPath := filepath.Join(os.Getenv("HOME"), ".ssh", "id_ed25519")
	edPubPath := filepath.Join(os.Getenv("HOME"), ".ssh", "id_ed25519.pub")
	if !IsExists(rsaPath) || !IsExists(rsaPubPath) {
		if !IsExists(edPath) || !IsExists(edPubPath) {
			defer TimeCost(time.Now(), "生成密钥")
			e, _, err := expect.Spawn("ssh-keygen -t ed25519", timeout)
			if err != nil {
				logger.Fatal(err)
			}
			defer e.Close()

			caser := []expect.Caser{
				&expect.BCase{R: "Enter", S: "\n"},
				&expect.BCase{R: "y/n", S: "y\n"},
				&expect.BCase{R: "fingerprint", S: "\n"},
			}

			for {
				output, _, _, err := e.ExpectSwitchCase(caser, timeout)

				if strings.Contains(output, "fingerprint") {
					break
				}
				if err != nil {
					e, _, err = expect.Spawn("ssh-keygen", timeout)
					continue
				}
			}
		}
	}
}
