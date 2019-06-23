/*
@File: util.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Desciption
------------      -------    --------    -----------
2019-06-11 02:07  Lucien     1.0         Init
*/
package util

import (
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var table = []rune("qwertyuiopasdfghjklzxcvbnm0123456789")

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
		return "", errors.New("length wrong") // "key's length show greater or equal than 3 and less or equal than 8: " + key)
	}
	if key[0] == '0' {
		return "", errors.New("leading zero") // "permanent key should not have leading zero: " + key)
	}
	flag, err := regexp.MatchString("^[0-9a-z]{3,8}$", key)
	if err != nil {
		return "", err
	}
	if !flag {
		return "", errors.New("reg failed") // "key's format checking failed, should only contains digital or lowercase letters: " + key)
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

// Generate a string using lowercase and digits with fixed length
func generator(length uint8) (string, error) {
	ret := make([]rune, length)
	for i := uint8(0); i < length; i++ {
		ret[i] = table[rand.Intn(len(table))]
	}
	return string(ret), nil
}

// Check str is able to insert or not
func check(str string) bool {
	flag, err := regexp.MatchString("[a-z]", str)
	if err != nil {
		return false
	}
	return flag
}

// Generate a string that contains at least one alphabet and not occur in temporary database on field key
func Generator() (string, error) {
	str, err := generator(8)
	if err != nil {
		return "", err
	}
	for !check(str) { // do {...} while (...)
		str, err = generator(8)
		if err != nil {
			return "", err
		}
	}
	return str, err
}

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
