package main

import (
    "bytes"
    "fmt"
    "log"
    "os/exec"
)

func main() {
    // define the command that you want to run
    // here "git checkout develop"
    cmd := exec.Command("git", "checkout", "main")
    // specify the working directory of the command
    cmd.Dir = "/home/robert/Code/dexetera"
    // create a buffer to store the output of your process
    var out bytes.Buffer
    // define the process standard output
    cmd.Stdout = &out
    // Run the command
    err := cmd.Run()
    if err != nil {
        // error case : status code of command is different from 0
        log.Fatal("git checkout err:", err)
    }
    fmt.Println(out.String())
}
