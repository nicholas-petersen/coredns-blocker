package blocker

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Blocker struct {
	Next  plugin.Handler
	Hosts map[string]struct{}
}

func newBlocker(file string) (*Blocker, error) {
	hosts, err := scanHosts(file)
	if err != nil {
		return nil, err
	}

	blocker := Blocker{
		Hosts: hosts,
	}

	return &blocker, nil
}

func scanHosts(file string) (map[string]struct{}, error) {
	hosts := map[string]struct{}{}

	ignore := []string{
		"#",
		"127",
		"255",
		"::1",
		"fe",
		"ff",
		"0.0.0.0 0.0.0.0",
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if ignoreLine(ignore, scanner.Text()) {
			continue
		}

		host, _ := strings.CutPrefix(scanner.Text(), "0.0.0.0")
		hosts[strings.TrimSpace(host)] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return hosts, fmt.Errorf("reading hosts file : %s", err)
	}

	return hosts, nil
}

func ignoreLine(ignore []string, text string) bool {
  for _, i := range ignore {
		if strings.HasPrefix(i, text) {
			return true
		}
	}

	return false
}

func (b Blocker) Name() string {
	return "blocker"
}

func (b Blocker) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	qname := state.Name()

	answer := []dns.RR{}

	if _, ok := b.Hosts[strings.TrimSuffix(qname, ".")]; ok {
			answer = append(answer, &dns.A{
				Hdr: dns.RR_Header{
					Name:   qname,
					Rrtype: dns.TypeA,
					Ttl:    3600,
					Class:  dns.ClassINET,
				},
				
				A: net.IPv4zero,
			})

		blockerCount.WithLabelValues(qname).Add(1)
		m := new(dns.Msg)
		m.SetReply(r)
		m.Answer = answer

		w.WriteMsg(m)
		return dns.RcodeSuccess, nil
	}

	return plugin.NextOrFailure(b.Name(), b.Next, ctx, w, r)
}
