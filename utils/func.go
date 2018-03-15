package utils

import "strings"

func StringContains(str string, arr []string) bool {
	for _, v := range arr {
		if strings.Contains(str, v) {
			return true
		}
	}
	return false
}
