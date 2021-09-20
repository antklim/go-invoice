package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

func runner(exit chan<- struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "exit" {
			exit <- struct{}{}
		}
	}
}

func main() {
	fmt.Println("Welcome to go-invoice.")
	fmt.Println(`Type "help" for more information.`)

	exit := make(chan struct{}, 1)
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go runner(exit)

	select {
	case <-osSignals:
		fmt.Println("received system signal")
	case <-exit:
		fmt.Println(`received "exit" command`)
	}
}
