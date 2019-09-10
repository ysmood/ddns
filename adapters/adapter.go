package adapters

// Adapter ...
type Adapter interface {
	// SetRecord ...
	SetRecord(subDomain, domainName, ip string) error
}

// New ...
func New(adapterName, config string) (adapter Adapter) {
	switch adapterName {
	case "dnspod":
		adapter = &Dnspod{token: config}
	}

	return
}
