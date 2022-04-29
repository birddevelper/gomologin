package gologin

import (
	"strconv"
	"time"
)

func currentTimeStamp() string{
	now := time.Now()
	unix := now.Unix()
	timestamp := strconv.FormatInt(unix, 10)

	return timestamp
}
