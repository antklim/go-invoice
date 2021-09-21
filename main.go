package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/antklim/go-invoice/cli"
	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage"
)

// TODO: add integration tests

func initCli(exit chan<- struct{}, svc *invoice.Service) *cli.Cli {
	if exit == nil {
		panic("cli: nil exit channel")
	}
	if svc == nil {
		panic("cli: nil invoice service")
	}

	c := cli.NewCli(os.Stdin, os.Stdout, exit)
	c.Handle("create", "Create new invoice", createHandler(svc))
	c.Handle("view", "View invoice.", nil)
	c.Handle("issue", "Issue invoice.", nil)
	c.Handle("pay", "Pay invoice.", nil)
	c.Handle("cancel", "Cancel invoice.", nil)
	c.Handle("add-item", "Add invoice item.", nil)
	c.Handle("delete-item", "Delete invoice item.", nil)
	c.Handle("update-customer", "Update invoice customer.", nil)
	return c
}

func initService() *invoice.Service {
	f := new(storage.Memory)
	strg := f.MakeStorage()
	return invoice.New(strg)
}

func main() {
	fmt.Println("Welcome to go-invoice.")

	exit := make(chan struct{}, 1)
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	svc := initService()

	c := initCli(exit, svc)
	go c.Run()

	select {
	case <-osSignals:
		fmt.Println()
	case <-exit:
	}
	fmt.Println("\nBye!")
}

func createHandler(svc *invoice.Service) cli.RunnerFunc {
	return func(out io.Writer) {
		inv, err := svc.CreateInvoice("John Doe")
		if err != nil {
			fmt.Fprintf(out, "create invoice failed: %v", err)
			return
		}

		fmt.Fprintf(out, "%q invoice successfully created", inv.ID)
	}
}

// func viewHandler(svc *invoice.Service) cli.RunnerFunc {
// 	return func(out io.Writer) {
// 		inv, err := svc.CreateInvoice("John Doe")
// 		if err != nil {
// 			fmt.Fprintf(out, "create invoice failed: %v", err)
// 			return
// 		}

// 		fmt.Fprintf(out, "%q invoice successfully created", inv.ID)
// 	}
// }
