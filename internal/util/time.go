package util

import "time"

func NowInWIB() time.Time {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Now()
	}
	return time.Now().In(location)
}
