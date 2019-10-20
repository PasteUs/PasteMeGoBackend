/*
@File: flag.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-07-25 08:36  Lucien     1.0         Init
*/
package flag

import (
	"flag"
	"fmt"
)

var (
	version bool
	Config string
)

func init() {
	flag.BoolVar(&version, "version", false, "Print version information")
	flag.StringVar(&Config, "c", "./config.json", "-c <config file>")
}

func Parse() bool { // return true for continue
	flag.Parse()
	if version {
		fmt.Println("3.2.0")
		return false
	}
	return true
}
