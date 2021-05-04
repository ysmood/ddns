package main

import (
	"encoding/json"
	"io"
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
	Log         bool

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
		lastIP, err = s.updateIP(lastIP)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(time.Duration(s.Interval))
	}
}

func (s *Service) updateIP(lastIP string) (string, error) {
	ip, err := s.getIP()
	if err != nil {
		return "", err
	}

	if ip == lastIP {
		return lastIP, nil
	}

	err = s.setIP(ip)
	if err != nil {
		return "", err
	}

	return ip, nil
}

func (s *Service) getIP() (ip string, err error) {
	if s.UsePublicIP {
		ip, err = myip.GetPublicIP()
	} else {
		ip, err = myip.GetInterfaceIP()
	}
	return
}

func (s *Service) setIP(ip string) error {
	l := log.New(io.Discard, "", 0)
	if s.Log {
		l = log.Default()
	}

	adapter := adapters.New(s.AdapterName, s.AdapterConfig, l)

	err := adapter.SetRecord(s.SubDomain, s.Domain, ip)
	if err != nil {
		return err
	}

	return nil
}
