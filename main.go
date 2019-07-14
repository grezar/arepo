package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Printf("Invalid number of arguments (expected 2, got %d)\n", len(os.Args))
        os.Exit(1)
    }
    os.Exit(run())
}

func run() int {
    return 0
}
