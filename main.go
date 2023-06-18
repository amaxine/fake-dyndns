package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/miekg/dns"
)

func main() {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion("ns1.google.com.", dns.TypeA)
	r, _, err := c.Exchange(m, net.JoinHostPort("8.8.8.8", "53"))
	if err != nil {
		log.Fatal(err)
	}
	var ipv4 string
	for _, a := range r.Answer {
		if ra, ok := a.(*dns.A); ok {
			ipv4 = ra.A.String()
		}
	}

	config := dns.ClientConfig{
		Servers: []string{"ns1.google.com"},
		Port:    "53",
	}

	m = new(dns.Msg)
	m.SetQuestion("o-o.myaddr.l.google.com.", dns.TypeTXT)
	m.RecursionDesired = false
	r, _, err = c.Exchange(m, net.JoinHostPort(config.Servers[0], config.Port))
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range r.Answer {
		if txt, ok := a.(*dns.TXT); ok {
			fmt.Println(txt.Txt[0])
		}
	}

	config = dns.ClientConfig{
		Servers: []string{ipv4},
		Port:    "53",
	}
	m = new(dns.Msg)
	m.SetQuestion("o-o.myaddr.l.google.com.", dns.TypeTXT)
	m.RecursionDesired = false
	r, _, err = c.Exchange(m, net.JoinHostPort(config.Servers[0], config.Port))
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range r.Answer {
		if txt, ok := a.(*dns.TXT); ok {
			fmt.Println(txt.Txt[0])
			ipv4 = txt.Txt[0]
		}
	}

	tc := dnsimple.StaticTokenHTTPClient(context.TODO(), os.Getenv("DNSIMPLE_TOKEN"))
	client := dnsimple.NewClient(tc)
	_, err = client.Zones.UpdateRecord(context.TODO(), "106272", "hormonal.party", int64(30490577), dnsimple.ZoneRecordAttributes{Content: ipv4})
	if err != nil {
		log.Fatal(err)
	}
}
