package utils

import "strings"

func SplicingString(str []string, s string) (newstr string) {
	newstr = strings.Join(str, s)
	return
}
