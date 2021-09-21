package cli

import "fmt"

type Command struct {
	name string
	desc string
	run  func()
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

var Commands = map[string]Command{
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
		if cmd, ok := Commands[name]; ok {
			fmt.Printf("%-25s%s\n", cmd.name, cmd.desc)
		}
	}
}
