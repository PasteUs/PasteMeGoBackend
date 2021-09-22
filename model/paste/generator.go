package paste

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano()) // 连续生成随机数，以当前纳秒数作为随机数种子
}

var (
	charset            = []rune("qwertyuiopasdfghjklzxcvbnm0123456789")
	charsetWithoutZero = []rune("qwertyuiopasdfghjklzxcvbnm123456789")
)

func getOne(cs []rune) rune {
	return cs[rand.Intn(len(cs))]
}

func generator(length int, zeroFirst bool) string {
	ret := make([]rune, length)

	if zeroFirst {
		ret[0] = '0'
	} else {
		ret[0] = getOne(charsetWithoutZero)
	}

	for i := 1; i < length; i++ {
		ret[i] = getOne(charset)
	}
	return string(ret)
}

func Generator(length int, zeroFirst bool, model interface{}) string {
	str := generator(length, zeroFirst)
	for Exist(str, model) {
		str = generator(length, zeroFirst)
	}
	return str
}
