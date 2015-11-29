package linecommand

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/xlab/closer"
	"log"
	"os/user"
	"path/filepath"
	"strings"
)

//Command is an individual command for your application 
type Command struct {
	Use       string                          // The one-line usage message.
	Short     string                          // The short description shown in the 'help' output.
	Long      string                          // The long message shown in the 'help <this-command>' output.
	Run       func(cmd *Command, args string) // Run runs the command.
	app       *App
//	Completer func(line string) (c []string)
	CompleteOptions []string
}

//App describes the application and holds your commands
type App struct {
	Commands     []Command
	DefaultRun   func(args string)
	line        *readline.Instance
	commandTitle string
	completeOptions []*readline.PrefixCompleter
}

//SetCommandTitle sets the string before > in the repl.
func (a *App) SetCommandTitle(title string) {
	a.commandTitle = title
	a.line.SetPrompt(fmt.Sprintf("%s>", title))
}

//AddCommand adds a new command to your application
func (a *App) AddCommand(command Command) {
	command.app = a
	a.Commands = append(a.Commands, command)
	
	//Completer options
	var children []*readline.PrefixCompleter
	for _, child := range command.CompleteOptions {
		children = append(children, readline.PcItem(child))
	}
	a.completeOptions = append(a.completeOptions, readline.PcItem(command.Use, children...))
}

//Run starts the application run loop
func (a *App) Run() {
	closer.Bind(a.cleanup)
	closer.Checked(a.internalRun, true)
	defer closer.Close()
}

/*func (a *App) Do(lr []rune, pos int) (newLine [][]rune, length int) {
	var c [][]rune
	line := string(lr)

	for _, command := range a.Commands {

		if strings.HasPrefix(strings.ToLower(line), command.Use+" ") {
			c = make([][]rune, 0)
			if command.Completer != nil {
				line = strings.Replace(line, command.Use+" ", "", 1)
				sc := command.Completer(line)
				for _, subCommand := range sc {
					c = append(c, []rune(command.Use + " " + subCommand)[pos:])
				}
			}

			return c, len(line) + 1
		}
		if strings.HasPrefix(command.Use, strings.ToLower(line)) {
			c = append(c, []rune(command.Use)[pos:])
		}
	}

	return c, len(line)
}*/

func (a *App) internalRun() error {
	if a.commandTitle == "" {
		a.commandTitle = ">"
	}
	var historyFile string
	usr, err := user.Current()
	// Only load history if we can get the user
	if err == nil {
		historyFile = filepath.Join(usr.HomeDir, ".command_history")
	}

	l, err := readline.NewEx(&readline.Config{
		Prompt:       a.commandTitle,
		HistoryFile:  historyFile,
		AutoComplete: readline.NewPrefixCompleter(a.completeOptions...),
	})
	if err != nil {
		return err
	}
	a.line = l
	log.SetOutput(l.Stderr())

	for {
		l, e := a.line.Readline()
		if e != nil {
			return e
		}
		if !a.parseCommand(l) {
			break // exit main loop
		}
	}

	return nil
}

func (a *App) parseCommand(cmd string) bool {
	lcmd := strings.TrimSpace(strings.ToLower(cmd))
	if strings.HasPrefix(lcmd, "exit") {
		// signal the program to exit
		return false
	}
	if strings.HasPrefix(lcmd, "help") {
		a.help()
		return true
	}

	for _, command := range a.Commands {
		if strings.HasPrefix(lcmd, command.Use) {
			cmd = strings.TrimSpace(strings.Replace(cmd, command.Use, "", -1))
			command.Run(&command, cmd)
			return true
		}
	}

	if a.DefaultRun != nil {
		a.DefaultRun(cmd)
	} else {
		fmt.Printf("'%s': command not found.\n", lcmd)
	}

	return true
}

func (a *App) help() {
	fmt.Println("Usage:")
	for _, command := range a.Commands {
		fmt.Printf("%s		%s\n", command.Use, command.Short)
	}
	fmt.Println("exit		quit the shell")
	fmt.Println("help		this help text")

}

func (a *App) cleanup() {
	a.line.Close()
}
