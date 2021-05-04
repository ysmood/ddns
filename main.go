package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/ysmood/ddns/adapters"
	"github.com/ysmood/gson"
	"github.com/ysmood/myip"
)

type Service struct {
	Domain      string
	SubDomain   string
	UsePublicIP bool
	Interval    Duration

	// check adapters.New
	AdapterName   string
	AdapterConfig gson.JSON
}

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	tmp, err := time.ParseDuration(v.(string))
	if err != nil {
		return err
	}
	*d = Duration(tmp)
	return nil
}

func main() {
	b, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	s := Service{}

	err = json.Unmarshal(b, &s)
	if err != nil {
		panic(err)
	}

	s.run()
}

func (s *Service) run() {
	lastIP := ""
	for {
		var err error
		lastIP, err = updateIP(s.UsePublicIP, lastIP, s.AdapterName, s.AdapterConfig, s.SubDomain, s.Domain)
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Duration(s.Interval))
	}
}

func updateIP(publicIP bool, lastIP, adapterName string, config gson.JSON, subDomain, domainName string) (string, error) {
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

func setIP(adapterName string, config gson.JSON, subDomain, domainName, ip string) error {
	adapter := adapters.New(adapterName, config)

	err := adapter.SetRecord(subDomain, domainName, ip)
	if err != nil {
		return err
	}

	log.Printf("[ddns] set ip: %s.%s -> %s\n", subDomain, domainName, ip)
	return nil
}
