package adapters

import (
	"bytes"
	"log"
	"net/http"
	"net/url"

	"github.com/ysmood/gson"
)

var _ Adapter = &Dnspod{}

// Dnspod ...
type Dnspod struct {
	token string
	log   *log.Logger
}

// Err ...
type Err struct {
	gson.JSON
}

// Error ...
func (e *Err) Error() string {
	return e.JSON.JSON("", "")
}

// SetRecord ...
func (pod *Dnspod) SetRecord(subDomain, domainName, ip string) error {
	recordID, err := pod.getRecordID(subDomain, domainName)
	if err != nil {
		return err
	}

	pod.log.Println("Record.Modify", subDomain, domainName, ip)

	_, err = pod.req("Record.Modify", &url.Values{
		"sub_domain":  {subDomain},
		"domain":      {domainName},
		"record_id":   {recordID},
		"record_type": {"A"},
		"record_line": {"默认"},
		"value":       {ip},
	})

	return err
}

func (pod *Dnspod) getRecordID(subDomain, domainName string) (string, error) {
	pod.log.Println("Record.List", subDomain, domainName)

	data, err := pod.req("Record.List", &url.Values{
		"sub_domain": {subDomain},
		"domain":     {domainName},
	})

	if err != nil {
		if e, ok := err.(*Err); ok && e.Get("code").Str() == "10" {
			pod.log.Println("Record.Create", subDomain, domainName)

			data, err = pod.req("Record.Create", &url.Values{
				"sub_domain":  {subDomain},
				"domain":      {domainName},
				"record_type": {"A"},
				"record_line": {"默认"},
				"value":       {"0.0.0.0"},
			})

			if err != nil {
				return "", err
			}

			return data.Get("record.id").Str(), nil
		}

		return "", err
	}

	return data.Get("records.0.id").Str(), nil
}

func (pod *Dnspod) req(path string, params *url.Values) (gson.JSON, error) {
	params.Add("login_token", pod.token)
	params.Add("format", "json")
	body := bytes.NewBufferString(params.Encode())

	req, err := http.NewRequest(http.MethodPost, "https://dnsapi.cn/"+path, body)
	if err != nil {
		return gson.JSON{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return gson.JSON{}, err
	}

	data := gson.New(res.Body)
	if data.Get("status.code").Str() != "1" {
		return data, &Err{data.Get("status")}
	}

	return data, err
}
