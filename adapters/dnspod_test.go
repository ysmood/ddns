package adapters_test

import (
	"io"
	"log"
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
	c := adapters.New("dnspod", gson.New(t.Open(false, "dnspod.json")), log.New(io.Discard, "", 0))

	t.Nil(c.SetRecord(t.Srand(8), "ysmood.org", "127.0.0.1"))
}
