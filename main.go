package main

import (
	"github.com/PasteUs/PasteMeGoBackend/config"
	"github.com/PasteUs/PasteMeGoBackend/router"
)

func main() {
	router.Run(config.Config.Address, config.Config.Port)
}
