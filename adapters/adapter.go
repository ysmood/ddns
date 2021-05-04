package adapters

import (
	"github.com/ysmood/gson"
)

// Adapter ...
type Adapter interface {
	// SetRecord ...
	SetRecord(subDomain, domainName, ip string) error
}

// New ...
func New(adapterName string, config gson.JSON) (adapter Adapter) {
	switch adapterName {
	case "dnspod":
		adapter = &Dnspod{token: config.Get("token").Str()}
	}

	return
}
