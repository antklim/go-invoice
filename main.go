package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/antklim/go-invoice/cli"
)

// TODO: start runner in a separate go-routine to handle user input without
// blocking main routine. Add separate channels to handle user commands errors
// and OS signals, like SIGTERM.

func runner(exit chan<- struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("> ")
	for scanner.Scan() {
		switch input := scanner.Text(); {
		case input == "exit":
			exit <- struct{}{}
			return
		case input == "help":
			cli.HelpCmd()
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
		fmt.Println()
	case <-exit:
	}
	fmt.Println("\nBye!")
}
