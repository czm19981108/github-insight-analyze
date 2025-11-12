package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateLockFile 创建锁文件，防止程序重复运行
func CreateLockFile(lockName string) (*os.File, error) {
	lockPath := filepath.Join(os.TempDir(), lockName+".lock")

	// 尝试创建锁文件（如果文件已存在则失败）
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		if os.IsExist(err) {
			return nil, fmt.Errorf("程序已在运行中（锁文件: %s）", lockPath)
		}
		return nil, fmt.Errorf("无法创建锁文件: %v", err)
	}

	return file, nil
}

// RemoveLockFile 删除锁文件
func RemoveLockFile(file *os.File) error {
	if file == nil {
		return nil
	}

	lockPath := file.Name()
	file.Close()
	return os.Remove(lockPath)
}
