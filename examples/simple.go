package main

import (
	"fmt"
	"github.com/jhalickman/linecommand"
)

func main() {
	app := linecommand.App{}

	echo_command := linecommand.Command{
		Use:   "echo",
		Short: "echo back what ever you say",
		Long:  "echo back what ever you say",
		Run: func(cmd *linecommand.Command, args string) {
			fmt.Println(args)
		},
	}

	app.DefaultRun = func(args string) {
		fmt.Printf("Oh noes that does not make sense. %s\n", args)
	}

	app.AddCommand(echo_command)

	app.Run()
}
