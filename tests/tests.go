package tests

import (
    "github.com/PasteUs/PasteMeGoBackend/model/paste"
    "github.com/PasteUs/PasteMeGoBackend/v2/model"
    "testing"
)

func init() {
    testing.Init()
    model.Init()
    paste.Init()
}
