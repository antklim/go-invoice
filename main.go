package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: start runner in a separate go-routine to handle user input without
// blocking main routine. Add separate channels to handle user commands errors
// and OS signals, like SIGTERM.

var commands = [][2]string{
	{"create", "Create new invoice"},
	{"view", "View invoice."},
	{"issue", "Issue invoice."},
	{"pay", "Pay invoice."},
	{"cancel", "Cancel invoice."},
	{"update-customer", "Update invoice customer."},
	{"add-item", "Add invoice item."},
	{"delete-item", "Delete invoice item."},
	{"exit", "Exit go-invoice."},
}

func runner() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprint(os.Stdout, ">>> ")
	cmdString, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Fprintln(os.Stdout, "you entered ", cmdString)
}

func main() {
	fmt.Println("go-invoice - track your invoices easy.")
	fmt.Println("Please use `help` to show available commands.")
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")
	fmt.Println("available commands:")
	for _, cmd := range commands {
		fmt.Printf("command: %s,\tdescription: %s\n", cmd[0], cmd[1])
	}
	defer fmt.Println("Bye!")
	runner()
}
