package main

import (
	"fmt"
	"time"

	"github.com/ysmood/ddns/adapters"
	"github.com/ysmood/kit"
	"github.com/ysmood/myip"
)

func main() {
	app := kit.TasksNew("ddns", "a tool for automate dns setup").Version("v0.2.3")

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
	lastIP := ""
	for {
		var err error
		lastIP, err = updateIP(userPublicIP, lastIP, adapterName, config, subDomain, domainName)
		if err != nil {
			kit.Err(err)
		}

		time.Sleep(interval)
	}
}

func updateIP(publicIP bool, lastIP, adapterName, config, subDomain, domainName string) (string, error) {
	ip, err := getIP(publicIP)
	if err != nil {
		return "", err
	}

	if ip == lastIP {
		return lastIP, nil
	}

	err = setIP(adapterName, config, subDomain, domainName, ip)
	if err != nil {
		return "", err
	}

	return ip, nil
}

func getIP(public bool) (ip string, err error) {
	if public {
		ip, err = myip.GetPublicIP()
	} else {
		ip, err = myip.GetInterfaceIP()
	}
	return
}

func setIP(adapterName, config, subDomain, domainName, ip string) error {
	adapter := adapters.New(adapterName, config)

	err := adapter.SetRecord(subDomain, domainName, ip)
	if err != nil {
		return err
	}

	kit.Log(fmt.Sprintf("[ddns] set ip: %s.%s -> %s\n", subDomain, domainName, ip))
	return nil
}
