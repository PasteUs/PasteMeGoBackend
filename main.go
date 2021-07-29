package main

import (
    "github.com/PasteUs/PasteMeGoBackend/config"
    "github.com/PasteUs/PasteMeGoBackend/model/v2"
    "github.com/PasteUs/PasteMeGoBackend/server"
)

func main() {
    v2.Init()
    server.Run(config.Get().Address, config.Get().Port)
}
