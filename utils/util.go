package utils

import (
	"errors"
	"strings"
)

func Exe2DeviceID(execution string) (string, error) {
	index := strings.Index(execution, "_")
	if index == -1 {
		return "", errors.New("execution格式错误")
	}
	return execution[0:index], nil
}
