package util

import (
    "fmt"
    "math/rand"
    "regexp"
    "strings"
    "time"
)

func init() {
    rand.Seed(time.Now().UnixNano()) // 连续生成随机数，以当前纳秒数作为随机数种子
}

// Parse token 有两种情况
// 情况 1: 只有 key; 情况 2: 有 key 和 password
func Parse(token string) (string, string) {
    buf := strings.Split(token, ",")
    if len(buf) == 1 { // 只有 key
        return buf[0], ""
    } else {
        return buf[0], buf[1] // 有 key 和 password
    }
}

func ValidChecker(key string) (string, error) {
    if len(key) > 8 || len(key) < 3 {
        return "", ErrWrongLength // key's length should at least 3 and at most 8
    }
    flag, err := regexp.MatchString("^[0-9a-z]{3,8}$", key) // 正则匹配
    if err != nil {
        return "", err
    }
    if !flag {
        return "", ErrWrongFormat // key's format checking failed, should only contains digital or lowercase letters
    }
    flag, err = regexp.MatchString("[a-z]", key) // 正则匹配
    if err != nil {
        return "", err
    }
    if !flag { // only digit
        return "permanent", nil // key 中只包括数字的是永久的
    }
    return "temporary", nil // key 中至少包括字母的是阅后即焚的
}

func LogFormat(IP string, format string, a ...interface{}) string {
    content := fmt.Sprintf(format, a...)
    return fmt.Sprintf("[%s] %s", IP, content)
}
