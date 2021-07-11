package generator

import (
    "github.com/PasteUs/PasteMeGoBackend/model"
    "math/rand"
    "regexp"
)

var table = []rune("qwertyuiopasdfghjklzxcvbnm0123456789")

// Generate a string using lowercase and digits with fixed length
func generator(length uint8) string {
    ret := make([]rune, length)
    for i := uint8(0); i < length; i++ {
        ret[i] = table[rand.Intn(len(table))]
    }
    return string(ret)
}

// Check str is able to insert or not
func check(key string) bool {
    if key[0] == '0' {
        return false
    }
    flag, err := regexp.MatchString("[a-z]", key)
    if err != nil {
        return false
    }
    return flag && !model.Exist(key)
}

// Generator Generate a string that contains at least one alphabet and not occur in temporary database on field key
func Generator() string {
    str := generator(8)
    for !check(str) { // do {...} while (...)
        str = generator(8)
    }
    return str
}
