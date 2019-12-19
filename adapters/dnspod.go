package adapters

import (
	"errors"

	"github.com/ysmood/kit"
)

// Dnspod ...
type Dnspod struct {
	token string
}

var _ Adapter = &Dnspod{}

// SetRecord ...
func (pod *Dnspod) SetRecord(subDomain, domainName, ip string) error {
	recordID, err := pod.getRecordID(subDomain, domainName)

	if err != nil {
		return err
	}

	_, err = pod.req("Record.Modify",
		"sub_domain", subDomain,
		"domain", domainName,
		"record_id", recordID,
		"record_type", "A",
		"record_line", "默认",
		"value", ip,
	)

	return err
}

func (pod *Dnspod) getRecordID(subDomain, domainName string) (string, error) {
	data, err := pod.req("Record.List",
		"sub_domain", subDomain,
		"domain", domainName,
	)

	if err != nil {
		if err.Error() == "No records" {
			data, err = pod.req("Record.Create",
				"sub_domain", subDomain,
				"domain", domainName,
				"record_type", "A",
				"record_line", "默认",
				"value", "0.0.0.0",
			)

			if err != nil {
				return "", err
			}

			return data.Get("record.id").String(), nil
		}

		return "", err
	}

	return data.Get("records.0.id").String(), nil
}

func (pod *Dnspod) req(path string, params ...interface{}) (kit.JSONResult, error) {
	params = append(params, "login_token", pod.token, "format", "json")

	data, err := kit.Req("https://dnsapi.cn/" + path).Post().Form(params...).JSON()
	if err != nil {
		return nil, err
	}

	if data.Get("status.code").String() != "1" {
		return data, errors.New(data.Get("status.message").String())
	}

	return data, err
}
