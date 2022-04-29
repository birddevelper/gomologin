package gologin

import (
	"strconv"
	"time"
)

func currentTimeStamp() string {
	now := time.Now()
	unix := now.Unix()
	timestamp := strconv.FormatInt(unix, 10)

	return timestamp
}

func roles_contains(s []string, e string) int {
	i := 0
	for _, a := range s {
		if a == e {
			return i
		}
		i++
	}
	return -1
}
