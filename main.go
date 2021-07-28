package main

import (
    "github.com/PasteUs/PasteMeGoBackend/config"
    "github.com/PasteUs/PasteMeGoBackend/model"
    "github.com/PasteUs/PasteMeGoBackend/server"
)

func main() {
    model.Init()
    server.Run(config.Get().Address, config.Get().Port)
}
