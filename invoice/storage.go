package invoice

type Storage interface {
	AddInvoice(Invoice) error
}
