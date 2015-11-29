# Linecommand


Linecommand is a libraray that helps setup commands for interactive CLI applications. Linecommand is built ontop of Readline (github.com/chzyer/readline) and the API is modeled after Cobra (https://github.com/spf13/cobra).


##Usage
To use linecommand:

1. Simply create an instance of ```linecommand.App{}``` 
2. Add an instance of ```linecommand.Command{}``` for each command you want.
3. Call ```app.Run()```

Example:
```golang

package main

import (
	"fmt"
	"github.com/jhalickman/linecommand"
)

func main() {
	app := linecommand.App{}

	echoCommand := linecommand.Command{
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

	app.AddCommand(echoCommand)

	app.Run()
}

```
