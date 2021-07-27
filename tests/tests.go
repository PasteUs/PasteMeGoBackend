package tests

import (
    "os"
    "path"
    "runtime"
    "testing"
)

func init() {
    if os.Getenv("UNITTEST") == "1" {
        testing.Init()

        _, filename, _, _ := runtime.Caller(0)
        dir := path.Join(path.Dir(filename), "..")
        err := os.Chdir(dir)
        if err != nil {
            panic(err)
        }
    }
}
