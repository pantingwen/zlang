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
	if (len(os.Args) > 1) {
		operation = os.Args[1];
		if operation == "run" || operation == "build" {
			if (len(os.Args) == 2) {
				fmt.Println("please input source file")
				return;
			}
			fileName = os.Args[2]
		}	else {
			operation = "run"
			fileName = os.Args[1]
		}
		fileContent, err := ioutil.ReadFile(fileName)
		if err != nil {
			panic(err)
		}
		sourceCode := string(fileContent)
		switch operation {
			case "run":
				cli.RunSourceCode(sourceCode)
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
