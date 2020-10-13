/*
@File: main.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-06-11 01:27  Lucien     1.0         Init
2019-06-13 01:59  Lucien     1.1         Split function, add mysql.Init()
2019-06-19 19:06  Irene      1.2         Fix package
2019-06-22 00:17  Lucien     1.3         Split into server
2020-10-13 16:02  Lx200916   1.4		 Add Burn after reading

*/
package main

import (
	"github.com/PasteUs/PasteMeGoBackend/config"
	"github.com/PasteUs/PasteMeGoBackend/crontab"
	"github.com/PasteUs/PasteMeGoBackend/flag"
	"github.com/PasteUs/PasteMeGoBackend/server"
)

func main() {
	if flag.Parse() {
		crontab.StartClean("CRON_TZ=Asia/Shanghai 0 16 * * ?")
		server.Run(config.Get().Address, config.Get().Port, config.Get().Log.Path != "")
	}
}
