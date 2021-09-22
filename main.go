package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
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
	c.Handle("add-item", "Add invoice item.", addItemHandler(svc))
	// c.Handle("delete-item", "Delete invoice item.", nil)
	c.Handle("update-customer", "Update invoice customer.", updateCustomerHandler(svc))
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
			fmt.Fprint(out, "create invoice failed: missing customer name\n")
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
			fmt.Fprint(out, "view invoice failed: missing invoice ID\n")
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
			fmt.Fprint(out, "issue invoice failed: missing invoice ID\n")
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
			fmt.Fprint(out, "pay invoice failed: missing invoice ID\n")
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
			fmt.Fprint(out, "cancel invoice failed: missing invoice ID\n")
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

func addItemHandler(svc *invoice.Service) cli.RunnerFunc {
	return func(out io.Writer, args ...string) {
		if len(args) < 4 || args[0] == "" || args[1] == "" || args[2] == "" || args[3] == "" {
			fmt.Fprint(out, "add invoice item failed: missing arguments\n")
			return
		}

		invID, productName := args[0], strings.TrimSpace(args[1])
		price, err := strconv.Atoi(strings.TrimSpace(args[2]))
		if err != nil {
			fmt.Fprintf(out, "add invoice item failed: invalid price argument: %v\n", err)
			return
		}

		qty, err := strconv.Atoi(strings.TrimSpace(args[3]))
		if err != nil {
			fmt.Fprintf(out, "add invoice item failed: invalid qty argument: %v\n", err)
			return
		}

		item, err := svc.AddInvoiceItem(invID, productName, price, qty)
		if err != nil {
			fmt.Fprintf(out, "add invoice item failed: %v\n", err)
			return
		}

		fmt.Fprintf(out, "item %q successfully added to invoice %q\n", item.ID, invID)
	}
}

func updateCustomerHandler(svc *invoice.Service) cli.RunnerFunc {
	return func(out io.Writer, args ...string) {
		if len(args) < 2 || args[0] == "" || args[1] == "" {
			fmt.Fprint(out, "update invoice customer failed: missing invoice ID and/or customer name\n")
			return
		}

		invID, name := args[0], strings.TrimSpace(args[1])
		err := svc.UpdateInvoiceCustomer(invID, name)
		if err != nil {
			fmt.Fprintf(out, "update invoice customer failed: %v\n", err)
			return
		}

		fmt.Fprintf(out, "%q invoice customer successfully updated\n", invID)
	}
}
