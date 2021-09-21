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
	c.Handle("view", "View invoice.", viewHandler(svc))
	c.Handle("issue", "Issue invoice.", issueHandler(svc))
	c.Handle("pay", "Pay invoice.", payHandler(svc))
	c.Handle("cancel", "Cancel invoice.", cancelHandler(svc))
	// c.Handle("add-item", "Add invoice item.", nil)
	// c.Handle("delete-item", "Delete invoice item.", nil)
	// c.Handle("update-customer", "Update invoice customer.", nil)
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
	return func(out io.Writer, args ...string) {
		if len(args) == 0 || args[0] == "" {
			fmt.Fprintf(out, "create invoice failed: missing customer name\n")
			return
		}

		inv, err := svc.CreateInvoice(args[0])
		if err != nil {
			fmt.Fprintf(out, "create invoice failed: %v\n", err)
			return
		}

		fmt.Fprintf(out, "%q invoice successfully created\n", inv.ID)
	}
}

func viewHandler(svc *invoice.Service) cli.RunnerFunc {
	return func(out io.Writer, args ...string) {
		if len(args) == 0 || args[0] == "" {
			fmt.Fprintf(out, "view invoice failed: missing invoice ID\n")
			return
		}

		invID := args[0]
		inv, err := svc.ViewInvoice(invID)
		if err != nil {
			fmt.Fprintf(out, "view invoice failed: %v\n", err)
			return
		}
		if inv == nil {
			fmt.Fprintf(out, "%q invoice not found\n", invID)
			return
		}

		fmt.Fprintf(out, "%#v\n", inv)
	}
}

func issueHandler(svc *invoice.Service) cli.RunnerFunc {
	return func(out io.Writer, args ...string) {
		if len(args) == 0 || args[0] == "" {
			fmt.Fprintf(out, "issue invoice failed: missing invoice ID\n")
			return
		}

		invID := args[0]
		err := svc.IssueInvoice(invID)
		if err != nil {
			fmt.Fprintf(out, "issue invoice failed: %v\n", err)
			return
		}

		fmt.Fprintf(out, "%q invoice successfully issued\n", invID)
	}
}

func payHandler(svc *invoice.Service) cli.RunnerFunc {
	return func(out io.Writer, args ...string) {
		if len(args) == 0 || args[0] == "" {
			fmt.Fprintf(out, "pay invoice failed: missing invoice ID\n")
			return
		}

		invID := args[0]
		err := svc.PayInvoice(invID)
		if err != nil {
			fmt.Fprintf(out, "pay invoice failed: %v\n", err)
			return
		}

		fmt.Fprintf(out, "%q invoice successfully paid\n", invID)
	}
}

func cancelHandler(svc *invoice.Service) cli.RunnerFunc {
	return func(out io.Writer, args ...string) {
		if len(args) == 0 || args[0] == "" {
			fmt.Fprintf(out, "cancel invoice failed: missing invoice ID\n")
			return
		}

		invID := args[0]
		err := svc.CancelInvoice(invID)
		if err != nil {
			fmt.Fprintf(out, "cancel invoice failed: %v\n", err)
			return
		}

		fmt.Fprintf(out, "%q invoice successfully canceled\n", invID)
	}
}
