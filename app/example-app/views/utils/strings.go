package viewUtils

import "strings"

func FirstToUpper(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.ToUpper(str[:1]) + str[1:]
}
