package invoice

type Storage interface {
	AddInvoice(Invoice) error
	FindInvoice(string) (*Invoice, error)
}
