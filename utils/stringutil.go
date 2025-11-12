package utils

import (
	"fmt"
	"strconv"
)

// StringToInt 尝试将字符串转为int，如果不是数字则返回错误
func StringToInt(s string) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("输入 '%s' 不是有效的整数: %v", s, err)
	}
	return num, nil
}

// IsNumeric 判断字符串是否为整数格式（正数或负数）
func IsNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
