package paste

import (
	"math/rand"
)

var (
	charset = []rune("qwertyuiopasdfghjklzxcvbnm0123456789")
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
	for exist(str, model) {
		str = generator(length, zeroFirst)
	}
	return str
}
