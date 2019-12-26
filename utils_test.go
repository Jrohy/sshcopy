package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	rsaPath    = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	rsaPubPath = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub")
)

func TestCheckIp(t *testing.T) {
	testMap := map[string]bool{
		"192.168.37.163": true,
		"192.2.256.3":    false,
		"192.168":        false,
		"192.168.6o.23":  false,
		"192.168.6.2O":   false,
		"192.I68.6o.23":  false,
		"192.68。6.23":    false,
	}
	for k, v := range testMap {
		if CheckIP(k) != v {
			t.Errorf("当ip为%s时checkIp函数无法通过", k)
		}
	}
}

// 测试id_rsa和id_rsa.pub都不存在的情况
func TestGenerateRsaAll(t *testing.T) {
	if IsExists(rsaPubPath) {
		_ = os.Remove(rsaPath)
	}
	if IsExists(rsaPubPath) {
		_ = os.Remove(rsaPubPath)
	}
	GenerateRsa()
	if !IsExists(rsaPubPath) || !IsExists(rsaPath) {
		t.Error("id_rsa和id_rsa.pub同时不存在时生成密钥失败!")
	}
}

// 测试id_rsa不存在的情况
func TestGenerateRsaPrivate(t *testing.T) {
	var oldTime, newTime time.Time
	if IsExists(rsaPath) {
		_ = os.Remove(rsaPath)
	}
	if IsExists(rsaPubPath) {
		if fileInfo, err := os.Stat(rsaPubPath); err != nil {
			t.Error(err)
		} else {
			oldTime = fileInfo.ModTime()
		}
	}
	GenerateRsa()

	if fileInfo, err := os.Stat(rsaPubPath); err != nil {
		t.Error(err)
	} else {
		newTime = fileInfo.ModTime()
	}
	if !IsExists(rsaPath) || oldTime == newTime {
		t.Error("id_rsa不存在时生成密钥失败!")
	}
}

// 测试id_rsa.pub不存在的情况
func TestGenerateRsaPublic(t *testing.T) {
	var oldTime, newTime time.Time
	if IsExists(rsaPubPath) {
		_ = os.Remove(rsaPubPath)
	}
	if IsExists(rsaPath) {
		if fileInfo, err := os.Stat(rsaPath); err != nil {
			t.Error(err)
		} else {
			oldTime = fileInfo.ModTime()
		}
	}
	GenerateRsa()

	if fileInfo, err := os.Stat(rsaPath); err != nil {
		t.Error(err)
	} else {
		newTime = fileInfo.ModTime()
	}
	if !IsExists(rsaPubPath) || oldTime == newTime {
		t.Error("id_rsa.pub不存在时生成密钥失败!")
	}
}

func TestIsExists(t *testing.T) {
	if !IsExists("utils.go") {
		t.Error("测试isExits函数失败!")
	}
}

func TestTimeCost(t *testing.T) {
	defer TimeCost(time.Now(), "测试耗时函数")
}

func TestTimeCostPTR(t *testing.T) {
	var timeCostPoint *string
	ptr := &timeCostPoint
	defer TimeCostPTR(time.Now(), ptr)
	test := "Test TimeCostPTR函数"
	timeCostPoint = &test
}
