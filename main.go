package main

import (
	"fmt"
	"time"

	"github.com/ysmood/ddns/adapters"

	"github.com/ysmood/kit"
	"github.com/ysmood/myip"
)

func main() {
	app := kit.TasksNew("ddns", "a tool for automate dns setup").Version("0.2.2")

	config := app.Flag("config", "the config for the adapter").Short('t').Required().String()
	adapterName := app.Flag("adapter", "").Default("dnspod").String()
	domainName := app.Flag("domain-name", "").Short('d').Required().String()
	subDomain := app.Flag("sub-domain", "").Short('s').Default("@").String()

	kit.Tasks().App(app).Add(
		kit.Task("run", "auto update dns").Init(func(cmd kit.TaskCmd) func() {
			cmd.Default()

			usePublicIP := cmd.Flag("use-public-ip", "").Short('p').Bool()
			interval := cmd.Flag("interval", "").Default("10s").Duration()

			return func() {
				run(*interval, *usePublicIP, *adapterName, *config, *subDomain, *domainName)
			}
		}),
		kit.Task("set", "set dns to ip").Init(func(cmd kit.TaskCmd) func() {
			ip := cmd.Flag("ip", "ip address to set").Required().String()

			return func() {
				kit.E(setIP(*adapterName, *config, *subDomain, *domainName, *ip))
			}
		}),
	).Do()
}

func run(interval time.Duration, userPublicIP bool, adapterName, config, subDomain, domainName string) {
	var err error

	for {
		err = updateIP(userPublicIP, adapterName, config, subDomain, domainName)

		if err != nil {
			kit.Err(err)
		}

		time.Sleep(interval)
	}
}

func updateIP(userPublicIP bool, adapterName, config, subDomain, domainName string) (err error) {
	var ip, lastIP string

	if userPublicIP {
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

	if lastIP == ip {
		return
	}

	err = setIP(adapterName, config, subDomain, domainName, ip)
	if err != nil {
		return err
	}

	lastIP = ip

	return
}

func setIP(adapterName, config, subDomain, domainName, ip string) error {
	adapter := adapters.New(adapterName, config)

	err := adapter.SetRecord(subDomain, domainName, ip)

	if err != nil {
		return err
	}

	kit.Log(fmt.Sprintf("set ip: %s.%s -> %s\n", subDomain, domainName, ip))

	return nil
}
