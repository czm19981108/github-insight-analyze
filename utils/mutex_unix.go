// +build !windows

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type Mutex struct {
	lockFile *os.File
}

// CreateNamedMutex 创建一个命名互斥体（Unix/Linux/Mac 使用文件锁）
func CreateNamedMutex(name string) (*Mutex, error) {
	// 在 Unix 系统上使用 flock 实现进程互斥
	// 锁文件存放在临时目录
	lockPath := filepath.Join(os.TempDir(), name+".lock")

	// 打开或创建锁文件
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, fmt.Errorf("无法创建锁文件: %v", err)
	}

	// 尝试获取排他锁（非阻塞）
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("程序已在运行中（锁文件: %s）", lockPath)
	}

	return &Mutex{lockFile: file}, nil
}

// Release 释放互斥体
func (m *Mutex) Release() error {
	if m.lockFile == nil {
		return nil
	}

	// 释放文件锁
	syscall.Flock(int(m.lockFile.Fd()), syscall.LOCK_UN)

	// 关闭并删除锁文件
	lockPath := m.lockFile.Name()
	m.lockFile.Close()
	os.Remove(lockPath)

	m.lockFile = nil
	return nil
}
