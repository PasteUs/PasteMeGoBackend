package util

import (
	"crypto/md5"
	"fmt"
	"github.com/PasteUs/PasteMeGoBackend/logging"
	"go.uber.org/zap"
	"strconv"
)

func Uint2string(value uint64) string {
	return strconv.FormatUint(value, 10) // 无符号整型转字符串
}

func String2uint(value string) uint64 {
	ret, err := strconv.ParseUint(value, 10, 64) // 字符串转无符号整型
	if err != nil {
		logging.Panic("convert string to uint failed", zap.String("string", value))
		return 0
	}
	return ret
}

// String2md5 String to MD5
func String2md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
