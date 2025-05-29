// config/utils.go
package config

import (
	"strconv"
)

func StrToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 1 // 默认页码或大小
	}
	return i
}
