package util

import (
	"strconv"
	"strings"
)

func Unit2Int(str any) int {
	switch str.(type) {
	case string:
		num := 0
		unit := ""
		for i, v := range []byte(str.(string)) {
			if v < 48 || v > 57 {
				num, _ = strconv.Atoi(str.(string)[:i])
				unit = str.(string)[i:]
			}
		}
		switch strings.ToUpper(unit) {
		case "K", "KB":
			return num * 1024
		case "M", "MB":
			return num * 1024 * 1024
		case "G", "GB":
			return num * 1024 * 1024 * 1024
		case "T", "TB":
			return num * 1024 * 1024 * 1024 * 1024
		case "PB", "EB":
			return num * 1024 * 1024 * 1024 * 1024 * 1024
		case "E", "EBP":
			return num * 1024 * 1024 * 1024 * 1024 * 1024 * 1024
		case "Z", "ZB":
			return num * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024
		case "Y", "YB":
			return num * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024
		default:
			return num
		}
	case int:
	case int64:
		return str.(int)
	}
	return 0
}
