package conv

import "strconv"

func Atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
