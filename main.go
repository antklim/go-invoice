package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "create", Description: "Create new invoice"},
		{Text: "view", Description: "View invoice."},
		{Text: "issue", Description: "Issue invoice."},
		{Text: "pay", Description: "Pay invoice."},
		{Text: "cancel", Description: "Cancel invoice."},
		{Text: "update-customer", Description: "Update invoice customer."},
		{Text: "add-item", Description: "Add invoice item."},
		{Text: "delete-item", Description: "Delete invoice item."},
		{Text: "exit", Description: "Exit go-invoice."},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func executor(in string) {
	fmt.Println("user entered", in)
}

func main() {
	fmt.Println("go-invoice - track your invoices easy.")
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")
	defer fmt.Println("Bye!")
	p := prompt.New(
		executor,
		completer,
		prompt.OptionTitle("go-invoice-prompt: interactive invoices client"),
		prompt.OptionPrefix(">>> "),
		prompt.OptionInputTextColor(prompt.Blue),
	)
	p.Run()
}
