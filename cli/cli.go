package cli

import (
	"bufio"
	"fmt"
	"io"
)

const (
	helpFormat = "%-25s%s\n"
	help       = "help"
	exit       = "exit"

	defaultPrompt = "> "
)

var reservedCommands = map[string]struct{}{
	help: {},
	exit: {},
}

type Runner interface {
	Run(io.Writer)
}

type command struct {
	name   string
	desc   string
	runner Runner
}

type Cli struct {
	commands  map[string]command
	scommands []command       // slice of commands sorted in order of registration
	cmdSrc    io.Reader       // commands source
	exit      chan<- struct{} // exit command notification channel
}

func NewCli(r io.Reader, exit chan<- struct{}) *Cli {
	return &Cli{cmdSrc: r, exit: exit}
}

// Handle registers the description and runner for the given command name.
func (cli *Cli) Handle(name, desc string, runner Runner) {
	if name == "" {
		panic("cli: invalid command name")
	}
	if desc == "" {
		panic("cli: invalid command description")
	}
	// if runner == nil {
	// 	panic("cli: nil runner")
	// }
	if _, exist := cli.commands[name]; exist {
		panic("cli: multiple registrations for " + name)
	}
	if _, reserved := reservedCommands[name]; reserved {
		panic("cli: " + name + " is a reserved command")
	}

	if cli.commands == nil {
		cli.commands = make(map[string]command)
	}
	cmd := command{name: name, desc: desc, runner: runner}
	cli.commands[name] = cmd
	cli.scommands = append(cli.scommands, cmd)
}

func (cli *Cli) Run() {
	scanner := bufio.NewScanner(cli.cmdSrc)
	cli.prompt()

	for scanner.Scan() {
		switch input := scanner.Text(); {
		case input == exit:
			cli.exit <- struct{}{}
			return
		case input == help:
			cli.help()
		}
		cli.prompt()
	}
}

func (cli *Cli) help() {
	for _, cmd := range cli.scommands {
		fmt.Printf(helpFormat, cmd.name, cmd.desc)
	}
	fmt.Printf(helpFormat, help, "Print this help message.")
	fmt.Printf(helpFormat, exit, "Exit go-invoice.")
}

func (cli *Cli) prompt() {
	fmt.Print(defaultPrompt)
}
