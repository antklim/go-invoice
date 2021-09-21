package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/antklim/go-invoice/cli"
)

// TODO: start runner in a separate go-routine to handle user input without
// blocking main routine. Add separate channels to handle user commands errors
// and OS signals, like SIGTERM.

func initCli(exit chan<- struct{}) *cli.Cli {
	c := cli.NewCli(os.Stdin, exit)
	c.Handle("create", "Create new invoice", nil)
	c.Handle("view", "View invoice.", nil)
	c.Handle("issue", "Issue invoice.", nil)
	c.Handle("pay", "Pay invoice.", nil)
	c.Handle("cancel", "Cancel invoice.", nil)
	c.Handle("add-item", "Add invoice item.", nil)
	c.Handle("delete-item", "Delete invoice item.", nil)
	c.Handle("update-customer", "Update invoice customer.", nil)
	return c
}

func main() {
	fmt.Println("Welcome to go-invoice.")
	fmt.Println(`Type "help" for more information.`)

	exit := make(chan struct{}, 1)
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	c := initCli(exit)
	go c.Run()

	select {
	case <-osSignals:
		fmt.Println()
	case <-exit:
	}
	fmt.Println("\nBye!")
}
