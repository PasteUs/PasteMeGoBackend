package main

import (
	"github.com/PasteUs/PasteMeGoBackend/config"
	"github.com/PasteUs/PasteMeGoBackend/model/paste"
	"github.com/PasteUs/PasteMeGoBackend/router"
	"github.com/PasteUs/PasteMeGoBackend/v2/model"
)

func main() {
	model.Init() // v2 paste model
	paste.Init() // v3 paste model
	router.Run(config.Get().Address, config.Get().Port)
}
