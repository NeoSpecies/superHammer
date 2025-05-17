package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// getProjectRoot 尝试通过查找 .git 目录或特定的标志文件来确定项目根目录。
func GetProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "duang_root.flag")); err == nil {
			return dir, nil
		}
		if dir == filepath.Dir(dir) {
			return "", fmt.Errorf("project root not found")
		}
		dir = filepath.Dir(dir)
	}
}

// resolvePath 根据项目根目录和相对路径构建完整路径。
func ResolvePath(relativePath string) (string, error) {
	projectRoot, err := GetProjectRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(projectRoot, relativePath), nil
}
