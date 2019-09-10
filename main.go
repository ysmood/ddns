package main

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/ysmood/ddns/adapters"

	"github.com/ysmood/myip"
	"gopkg.in/alecthomas/kingpin.v2"
)

type ddns struct {
	config       *string
	domainName   *string
	subDomain    *string
	userPublicIP *bool
	adapter      *string
	interval     *int

	err *log.Logger
	std *log.Logger

	ip string
}

func main() {
	service := &ddns{
		config:       kingpin.Flag("config", "the config for the adapter").Short('t').Required().String(),
		domainName:   kingpin.Flag("domain-name", "").Short('d').Required().String(),
		subDomain:    kingpin.Flag("sub-domain", "").Short('s').String(),
		userPublicIP: kingpin.Flag("use-public-ip", "").Short('p').Bool(),
		adapter:      kingpin.Flag("adapter", "").Default("dnspod").String(),
		interval:     kingpin.Flag("interval", "").Default("1").Int(),

		err: log.New(os.Stderr, "", log.LstdFlags),
		std: log.New(os.Stdout, "", log.LstdFlags),
	}

	kingpin.Version("0.0.1")
	kingpin.Parse()

	go service.run()

	runtime.Goexit()
}

func (service *ddns) run() {
	var err error

	for {
		err = service.updateIP()

		if err != nil {
			service.err.Println(err)
		}

		time.Sleep(time.Duration(*service.interval) * time.Second)
	}
}

func (service *ddns) updateIP() (err error) {
	var ip string

	if *service.userPublicIP {
		ip, err = myip.GetPublicIP()

		if err != nil {
			return err
		}

	} else {
		ip, err = myip.GetInterfaceIP()

		if err != nil {
			return err
		}
	}

	if service.ip == ip {
		return
	}

	adapter := adapters.New(*service.adapter, *service.config)

	err = adapter.SetRecord(*service.subDomain, *service.domainName, ip)

	if err != nil {
		return err
	}

	service.ip = ip

	service.std.Printf("set ip: %s.%s -> %s\n", *service.subDomain, *service.domainName, ip)

	return
}
