package adapters

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

// NewDnspod ...
func NewDnspod(opts *Options) Adapter {
	return &dnspod{opts}
}

type dnspod struct {
	opts *Options
}

// SetRecord ...
func (pod *dnspod) SetRecord(ip string) error {
	recordID, err := pod.getRecordID()

	if err != nil {
		return err
	}

	pod.req("Record.Modify", url.Values{
		"domain":      {pod.opts.DomainName},
		"sub_domain":  {pod.opts.SubDomain},
		"record_id":   {recordID},
		"record_type": {"A"},
		"record_line": {"默认"},
		"value":       {ip},
	})

	return nil
}

func (pod *dnspod) getRecordID() (string, error) {
	data, err := pod.req("Record.List", url.Values{
		"domain":     {pod.opts.DomainName},
		"sub_domain": {pod.opts.SubDomain},
	})

	if err != nil {
		if err.Error() == "No records" {
			data, err = pod.req("Record.Create", url.Values{
				"domain":      {pod.opts.DomainName},
				"sub_domain":  {pod.opts.SubDomain},
				"record_type": {"A"},
				"record_line": {"默认"},
				"value":       {"0.0.0.0"},
			})

			if err != nil {
				return "", err
			}

			return data.Get("record.id").String(), nil
		}

		return "", err
	}

	return data.Get("records.0.id").String(), nil
}

func (pod *dnspod) req(path string, param url.Values) (data gjson.Result, err error) {
	param["login_token"] = []string{pod.opts.Token}
	param["format"] = []string{"json"}

	if param["sub_domain"][0] == "" {
		delete(param, "sub_domain")
	}

	var r *http.Response

	r, err = http.PostForm(
		fmt.Sprintf("https://dnsapi.cn/%s", path),
		param,
	)

	if err != nil {
		return data, err
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return data, err
	}

	data = gjson.Parse(string(body))

	if err != nil {
		return data, err
	}

	if data.Get("status.code").String() != "1" {
		return data, errors.New(data.Get("status.message").String())
	}

	return data, err
}

