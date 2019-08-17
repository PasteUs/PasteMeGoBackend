/*
@File: util.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-06-11 02:07  Lucien     1.0         Init
*/
package util

import (
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano()) // 连续生成随机数，以当前纳秒数作为随机数种子
}

// token有两种情况
// 情况1：只有key；情况2：有key和password
func Parse(token string) (string, string) {
	buf := strings.Split(token, ",")
	if len(buf) == 1 { // 只有key
		return buf[0], ""
	} else {
		return buf[0], buf[1] // 有key和password
	}
}

func ValidChecker(key string) (string, error) {
	if len(key) > 8 || len(key) < 3 {
		return "", errors.New("wrong length") // key's length should at least 3 and at most 8
	}
	flag, err := regexp.MatchString("^[0-9a-z]{3,8}$", key) // 正则匹配
	if err != nil {
		return "", err
	}
	if !flag {
		return "", errors.New("wrong format") // key's format checking failed, should only contains digital or lowercase letters
	}
	flag, err = regexp.MatchString("[a-z]", key) // 正则匹配
	if err != nil {
		return "", err
	}
	if !flag { // only digit
		return "permanent", nil // key 中只包括数字的是永久的
	}
	return "temporary", nil // key中至少包括字母的是阅后即焚的
}

func LoggerInfo(IP string, content string) string {
	return fmt.Sprintf("[%s] %s", IP, content)
}

func GetEnvOrFatal(key string) string {
	value, exist := os.LookupEnv(key)
	if !exist {
		logger.Fatal(fmt.Sprintf("Enviromental variable %s is not set", key))
	}
	return value
}
