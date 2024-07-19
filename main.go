package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Println("Welcome to Monkey language , user is: ", user.Username)
	repl.CStart(os.Stdin, os.Stdout)
}
