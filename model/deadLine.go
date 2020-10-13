/*
@File: deadLine.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2020-10-12 20:06  Lx200916   1.0	     Add Burn after reading
*/
package model

import (
	"database/sql/driver"
	"fmt"
	"github.com/wonderivan/logger"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

// Days:Hours 时间字符串
type DeadLine time.Time

func (t *DeadLine) UnmarshalJSON(data []byte) (err error) {
	// 空值不进行解析
	if len(data) == 2 {
		t = nil
		return
	}

	timelists := strings.Split(strings.Trim(string(data), "\""), ":")
	days, err := strconv.Atoi(timelists[0])
	hours := 0
	if len(timelists) > 1 {
		hours, err = strconv.Atoi(timelists[0])
	}
	if err != nil {
		logger.Warn("DeadLine Parse Failed")
	}
	now := time.Now()

	*t = DeadLine(now.Add(time.Duration(24*days+hours) * time.Hour))
	fmt.Println(*t)

	// 指定解析的格式

	return err
}
func (t DeadLine) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}
func (t DeadLine) Value() (driver.Value, error) {
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return time.Time(t).Format(TimeFormat), nil
}

func (t *DeadLine) Scan(v interface{}) error {
	tTime, _ := time.Parse("2006-01-02 15:04:05 +0800 CST", v.(time.Time).String())
	*t = DeadLine(tTime)
	return nil
}

func (t DeadLine) String() string {
	return time.Time(t).Format(TimeFormat)
}
