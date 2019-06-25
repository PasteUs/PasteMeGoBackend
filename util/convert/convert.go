/*
@File: strings.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Desciption
------------      -------    --------    -----------
2019-06-25 09:00  Lucien     1.0         Init
*/
package convert

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

func Uint2string(value uint64) string {
	return strconv.FormatUint(value, 10)
}

func String2uint(value string) uint64 {
	ret, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		// TODO
		return 0
	}
	return ret
}

func String2md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
