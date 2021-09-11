package invoice

type Storage interface {
	AddInvoice(Invoice) error
	FindInvoice(string) (*Invoice, error)
	UpdateInvoice(Invoice) error
}

type StorageFactory interface {
	MakeStorage() Storage
}
