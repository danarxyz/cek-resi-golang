package utils

import (
	"fmt"
	"time"

	"github.com/goodsign/monday"
)

func ParseTime(dateTimeStr string) time.Time {
	layout := "02 Jan 15:04 MST"
	t, err := monday.Parse(layout, dateTimeStr, monday.LocaleIdID) // Gunakan locale Indonesia
	if err != nil {
		fmt.Println(err)
		return time.Time{}
	}
	return t
}
