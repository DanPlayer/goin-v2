package utils

import (
	"time"
)

const (
	DATE_FORMAT             = "2006-01-02"
	DATETIME_FORMAT         = "2006-01-02 15:04:05"
	DATETIMEWITHZONE_FORMAT = "2006-01-02 15:04:05 -07"
	TIME_FORMAT             = "15:04:05"
)

var (
	ASTM, _ = time.LoadLocation("Asia/Shanghai")
)

func Now() time.Time {
	return time.Now().In(ASTM)
}
