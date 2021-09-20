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
	{"add-item", "Add invoice item."},
	{"delete-item", "Delete invoice item."},
	{"update-customer", "Update invoice customer."},
	{"help", "Print this help message."},
	{"exit", "Exit go-invoice."},
}

func runner(exit chan<- struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("> ")
	for scanner.Scan() {
		switch input := scanner.Text(); {
		case input == "exit":
			exit <- struct{}{}
			return
		}
		fmt.Printf("> ")
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
		fmt.Println("\nreceived system signal")
	case <-exit:
		fmt.Println(`received "exit" command`)
	}
}
