package adapters_test

import (
	"testing"

	"github.com/ysmood/ddns/adapters"
	"github.com/ysmood/got"
	"github.com/ysmood/gson"
)

type DNSPOD struct {
	got.G
}

func Test(t *testing.T) {
	got.Each(t, DNSPOD{})
}

func (t DNSPOD) Basic() {
	c := adapters.New("dnspod", gson.New(t.Open(false, "dnspod.json")))

	t.Nil(c.SetRecord(t.Srand(8), "ysmood.org", "127.0.0.1"))
}
