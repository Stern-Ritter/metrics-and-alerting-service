package utils

import (
	"strings"
)

func Contains(s []string, v string) bool {
	for _, e := range s {
		if strings.EqualFold(e, v) {
			return true
		}
	}
	return false
}
