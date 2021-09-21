package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	helpFormat = "%-25s%s\n"
	help       = "help"
	exit       = "exit"

	defaultPrompt  = "> "
	defaultNameDlm = " " // a delimeter between command name and arguments
	defaultArgsDlm = "," // a delimeter between arguments
	commandParts   = 2   // a number of parte to split a command input
)

var reservedCommands = map[string]struct{}{
	help: {},
	exit: {},
}

type Runner interface {
	Run(io.Writer, ...string)
}

type RunnerFunc func(io.Writer, ...string)

func (f RunnerFunc) Run(out io.Writer, args ...string) {
	f(out, args...)
}

type command struct {
	name   string
	desc   string
	runner Runner
}

type Cli struct {
	commands  map[string]command
	scommands []command       // slice of commands sorted in order of registration
	src       io.Reader       // commands source
	dst       io.Writer       // commands output destination
	exit      chan<- struct{} // exit command notification channel
}

func NewCli(src io.Reader, dst io.Writer, exit chan<- struct{}) *Cli {
	return &Cli{src: src, dst: dst, exit: exit}
}

// Handle registers the description and runner for the given command name.
func (cli *Cli) Handle(name, desc string, runner Runner) {
	if name == "" {
		panic("cli: invalid command name")
	}
	if desc == "" {
		panic("cli: invalid command description")
	}
	if runner == nil {
		panic("cli: nil runner")
	}
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
	cli.helpPrompt()
	cli.prompt()

	scanner := bufio.NewScanner(cli.src)
	for scanner.Scan() {
		input := strings.SplitN(scanner.Text(), defaultNameDlm, commandParts)
		name := input[0]

		switch name {
		case exit:
			cli.exit <- struct{}{}
			return
		case help:
			cli.help()
		default:
			if cmd, ok := cli.commands[name]; ok {
				var args []string
				if len(input) == commandParts {
					args = strings.Split(input[1], defaultArgsDlm)
				}
				cmd.runner.Run(cli.dst, args...)
			} else {
				cli.unknownCommand(name)
			}
		}
		cli.prompt()
	}
}

func (cli *Cli) help() {
	for _, cmd := range cli.scommands {
		fmt.Fprintf(cli.dst, helpFormat, cmd.name, cmd.desc)
	}
	fmt.Fprintf(cli.dst, helpFormat, help, "Print this help message.")
	fmt.Fprintf(cli.dst, helpFormat, exit, "Exit go-invoice.")
}

func (cli *Cli) prompt() {
	fmt.Fprint(cli.dst, defaultPrompt)
}

func (cli *Cli) helpPrompt() {
	fmt.Fprintln(cli.dst, `Type "help" for more information.`)
}

func (cli *Cli) unknownCommand(name string) {
	fmt.Fprintf(cli.dst, "Unknown command %q entered.\n", name)
}
