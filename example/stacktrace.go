package main

import (
    "fmt"
    "io/ioutil"

    "git.touchdevops.com/lib/goerror"
)

func main() {
    err := read()
    fmt.Println(err.(*goerror.GoError).StackTrace())
}

func read() error {
    return readError()
}

func readError() error {
    _, err := ioutil.ReadFile("/tmp/notfound")
    if err != nil {
        return goerror.DefineInternalServerError("UnableReadFile", "Not found file").WithCause(err)
    }

    return nil
}
