package stringutil

import "strings"

//RightPad2Len https://github.com/git-time-metric/gtm/blob/master/util/string.go#L55-L61
func RightPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

//LeftPad2Len https://github.com/git-time-metric/gtm/blob/master/util/string.go#L63-L69
func LeftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}
