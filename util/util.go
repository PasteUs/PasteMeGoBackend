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
	rand.Seed(time.Now().UnixNano())
}

func Parse(token string) (string, string) {
	buf := strings.Split(token, ",")
	if len(buf) == 1 {
		return buf[0], ""
	} else {
		return buf[0], buf[1]
	}
}

func ValidChecker(key string) (string, error) {
	if len(key) > 8 || len(key) < 3 {
		return "", errors.New("wrong length") // key's length should at least 3 and at most 8
	}
	flag, err := regexp.MatchString("^[0-9a-z]{3,8}$", key)
	if err != nil {
		return "", err
	}
	if !flag {
		return "", errors.New("wrong format") // key's format checking failed, should only contains digital or lowercase letters
	}
	flag, err = regexp.MatchString("[a-z]", key)
	if err != nil {
		return "", err
	}
	if !flag { // only digit
		return "permanent", nil
	}
	return "temporary", nil
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
