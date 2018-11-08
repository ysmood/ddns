package adapters

// Options ...
type Options struct {
	DomainName string
	SubDomain  string
	Token      string
}

// Adapter ...
type Adapter interface {
	SetRecord(ip string) error
}
