//go:build windows
// +build windows

package utils

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procCreateMutex  = kernel32.NewProc("CreateMutexW")
	procReleaseMutex = kernel32.NewProc("ReleaseMutex")
	procCloseHandle  = kernel32.NewProc("CloseHandle")
)

//   📋 工作原理：
//   1. 第一个进程启动：
//     - 调用 CreateMutexW("Global\\OSSInsightTestEmailSender")
//     - Windows 内核创建互斥体，返回句柄
//     - 程序继续执行，发送邮件
//     - 完成后释放互斥体
//   2. 第二个进程尝试启动（几乎同时）：
//     - 调用 CreateMutexW("Global\\OSSInsightTestEmailSender")
//     - Windows 内核发现互斥体已存在，返回
//   ERROR_ALREADY_EXISTS (错误码 183)
//     - 程序检测到错误，输出警告并立即退出
//     - 不会发送邮件
//   🔍 代码位置：
//   - 互斥体实现：utils/mutex_windows.go
//   - 使用位置：cmd/test-email/main.go:17-23

type Mutex struct {
	handle syscall.Handle
}

// CreateNamedMutex 创建一个命名互斥体（Windows 系统级别）
func CreateNamedMutex(name string) (*Mutex, error) {
	mutexName, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return nil, fmt.Errorf("无法转换互斥体名称: %v", err)
	}

	// 调用 CreateMutexW
	ret, _, err := procCreateMutex.Call(
		0,                                  // 默认安全属性
		0,                                  // 不是初始拥有者
		uintptr(unsafe.Pointer(mutexName)), // 互斥体名称
	)

	handle := syscall.Handle(ret)
	if handle == 0 {
		return nil, fmt.Errorf("创建互斥体失败: %v", err)
	}

	// 检查是否已经存在（ERROR_ALREADY_EXISTS = 183）
	if err != nil && err.(syscall.Errno) == 183 {
		// 互斥体已存在，说明程序已在运行
		procCloseHandle.Call(uintptr(handle))
		return nil, fmt.Errorf("程序已在运行中（互斥体: %s）", name)
	}

	return &Mutex{handle: handle}, nil
}

// Release 释放互斥体
func (m *Mutex) Release() error {
	if m.handle == 0 {
		return nil
	}

	procReleaseMutex.Call(uintptr(m.handle))
	procCloseHandle.Call(uintptr(m.handle))
	m.handle = 0
	return nil
}
