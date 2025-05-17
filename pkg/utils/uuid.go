package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"time"
)

func GenerateUUID() (string, error) {
	// 获取本地网络接口的MAC地址
	ifas, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get local network interfaces: %v", err)
	}
	if len(ifas) == 0 {
		return "", fmt.Errorf("no network interface found")
	}

	var mac string
	// 尝试找到一个有效的MAC地址
	for _, ifa := range ifas {
		if len(ifa.HardwareAddr) != 0 {
			mac = hex.EncodeToString(ifa.HardwareAddr)[:6] // 只取MAC地址的前6个字符
			break
		}
	}
	if mac == "" {
		return "", fmt.Errorf("no valid network interface with a MAC address found")
	}

	// 使用当前时间（取纳秒的后8位）和随机数生成UUID
	timestamp := time.Now().UnixNano() % 1e8 // 取后8位数字确保长度和增加随机性
	randomBytes := make([]byte, 4)           // 产生4个随机字节
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	uuid := fmt.Sprintf("%x-%x-%x",
		timestamp,   // 时间戳后8位
		randomBytes, // 4个随机字节
		mac,         // MAC地址前6个字符
	)

	return uuid, nil
}
