/*
@File: generator_test.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-07-04 02:42  Lucien     1.0         Init
*/
package generator

import (
	"testing"
)

func Test_generator(t *testing.T) {
	for i := 0; i < 8; i++ {
		str := generator(8)
		t.Log(str)
	}
}
