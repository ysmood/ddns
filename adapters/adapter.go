package adapters

import (
	"log"

	"github.com/ysmood/gson"
)

// Adapter ...
type Adapter interface {
	// SetRecord ...
	SetRecord(subDomain, domainName, ip string) error
}

// New ...
func New(adapterName string, config gson.JSON, log *log.Logger) (adapter Adapter) {
	switch adapterName {
	case "dnspod":
		adapter = &Dnspod{token: config.Get("token").Str(), log: log}
	}

	return
}
