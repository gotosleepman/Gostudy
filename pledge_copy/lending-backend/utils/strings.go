package utils

import "strconv"

func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func StringToInt(i string) int {
	j, _ := strconv.Atoi(i)
	return j
}
