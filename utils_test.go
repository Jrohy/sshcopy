package main

import (
	"testing"
	"time"
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
