package cli

import (
	"fmt"
	"io"
)

const (
	helpFormat = "%-25s%s\n"
)

type Runner interface {
	Run(io.Writer)
}

type command struct {
	name   string
	desc   string
	runner Runner
}

var orderedCommands = []string{
	"create",
	"view",
	"issue",
	"pay",
	"cancel",
	"add-item",
	"delete-item",
	"update-customer",
	"help",
	"exit",
}

var commands = map[string]command{
	"create":          {name: "create", desc: "Create new invoice"},
	"view":            {name: "view", desc: "View invoice."},
	"issue":           {name: "issue", desc: "Issue invoice."},
	"pay":             {name: "pay", desc: "Pay invoice."},
	"cancel":          {name: "cancel", desc: "Cancel invoice."},
	"add-item":        {name: "add-item", desc: "Add invoice item."},
	"delete-item":     {name: "delete-item", desc: "Delete invoice item."},
	"update-customer": {name: "update-customer", desc: "Update invoice customer."},
	"help":            {name: "help", desc: "Print this help message."},
	"exit":            {name: "exit", desc: "Exit go-invoice."},
}

func HelpCmd() {
	for _, name := range orderedCommands {
		if cmd, ok := commands[name]; ok {
			fmt.Printf("%-25s%s\n", cmd.name, cmd.desc)
		}
	}
}

type Cli struct {
	commands  map[string]command
	scommands []command       // slice of commands sorted in order of registration
	exit      chan<- struct{} // exit command notification channel
}

func NewCli(exit chan<- struct{}) *Cli {
	return &Cli{exit: exit}
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

	if cli.commands == nil {
		cli.commands = make(map[string]command)
	}
	cmd := command{name: name, desc: desc, runner: runner}
	cli.commands[name] = cmd
	cli.scommands = append(cli.scommands, cmd)
}

func (cli *Cli) Help() {
	for _, cmd := range cli.scommands {
		fmt.Printf(helpFormat, cmd.name, cmd.desc)
	}
	fmt.Printf(helpFormat, "help", "Print this help message.")
	fmt.Printf(helpFormat, "exit", "Exit go-invoice.")
}

func (cli *Cli) Exit() {
	cli.exit <- struct{}{}
}
