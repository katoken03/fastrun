package main

import (
    "fmt"
    "os"
    
    "github.com/kato/fastrun/cmd"
)

var (
    version = "dev"
    commit  = "unknown"
    date    = "unknown"
)

func main() {
    cmd.SetVersionInfo(version, commit, date)
    if err := cmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
