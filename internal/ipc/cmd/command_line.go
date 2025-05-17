package cmd

import (
	"bytes"
	"os/exec"
)

// ExecCommand 执行命令行命令并返回其标准输出和错误输出
func ExecCommand(command string, args ...string) (string, string, error) {
	cmd := exec.Command(command, args...)

	// 创建一个缓冲区来捕获命令的输出
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行命令
	err := cmd.Run()
	if err != nil {
		return "", "", err
	}

	// 返回命令的输出
	return stdout.String(), stderr.String(), nil
}
