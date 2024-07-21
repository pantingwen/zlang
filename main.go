package main

import (
	"fmt"
	"z/repl"
	"os"
	"os/user"
	"io/ioutil"
	"z/cli"
)

func main() {
	var operation string;
	var fileName string;
	mode := "vm"
	if (len(os.Args) > 1) {
		operation = os.Args[1];
		if operation == "run" || operation == "build" {
			if (len(os.Args) == 2) {
				fmt.Println("please input source file")
				return;
			}
			fileName = os.Args[2]
			if (len(os.Args) == 4) {
				mode = os.Args[3]
			}
		}	else {
			operation = "run"
			fileName = os.Args[1]
			if (len(os.Args) == 3) {
				mode = os.Args[2]
			}
		}
		fileContent, err := ioutil.ReadFile(fileName)
		if err != nil {
			panic(err)
		}
		sourceCode := string(fileContent)
		switch operation {
			case "run":
				cli.RunSourceCode(sourceCode, mode)
				return
			case "build":
				cli.BuildSourceCode(sourceCode, fileName)
			return
		}
	}
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Println("Welcome to z language , user is: ", user.Username, " type code to execute")
	repl.CStart(os.Stdin, os.Stdout)
}
