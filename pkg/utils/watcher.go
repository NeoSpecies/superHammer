package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

// Watcher 包含文件监视器所需的配置
type Watcher struct {
	Folders []string
	PIDFile string
	Script  string
}

// New 创建一个新的 Watcher
func New(folders []string, pidFile, script string) *Watcher {
	return &Watcher{
		Folders: folders,
		PIDFile: pidFile,
		Script:  script,
	}
}

// Start 开始监视文件变化
func (w *Watcher) Start(ctx context.Context) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for _, folder := range w.Folders {
		err = watcher.Add(folder)
		if err != nil {
			log.Fatalf("Failed to add watcher for folder %s: %v", folder, err)
		}
	}

	// 监听文件系统事件和上下文的取消事件
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("Modified file:", event.Name)
				w.restartPHPProcess()
			}
		case err := <-watcher.Errors:
			log.Println("Error:", err)
		case <-ctx.Done(): // 添加这个case来响应上下文的取消事件
			fmt.Println("File watcher is stopping...")
			return
		}
	}
}

// restartPHPProcess 重启 PHP 进程
func (w *Watcher) restartPHPProcess() {
	// 读取文件中的现有 PID 并结束进程
	if pidData, err := os.ReadFile(w.PIDFile); err == nil {
		if pid, err := strconv.Atoi(string(pidData)); err == nil {
			syscall.Kill(pid, syscall.SIGTERM) // 尝试优雅地结束进程
		}
	}

	// 启动新的 PHP 进程
	cmd := exec.Command("php", w.Script)
	if err := cmd.Start(); err != nil {
		log.Fatalf("PHP process failed to start: %s", err)
	}

	// 将新的 PID 保存到文件
	newPid := cmd.Process.Pid
	os.WriteFile(w.PIDFile, []byte(fmt.Sprintf("%d", newPid)), 0644)

	fmt.Printf("PHP process started with PID: %d\n", newPid)
}
