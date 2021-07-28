package tests

import (
    "github.com/PasteUs/PasteMeGoBackend/model"
    "os"
    "path"
    "runtime"
    "testing"
)

func init() {
    testing.Init()

    _, filename, _, _ := runtime.Caller(0)
    dir := path.Join(path.Dir(filename), "..")
    err := os.Chdir(dir)
    if err != nil {
        panic(err)
    }

    model.Init()
}
