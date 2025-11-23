package dnstest

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

// Result holds measurements for a DNS test run.
type Result struct {
	Sent       int       `json:"sent"`
	Success    int       `json:"success"`
	SamplesMs  []float64 `json:"samples_ms"`
	AvgMs      float64   `json:"avg_ms"`
	PacketLoss float64   `json:"packet_loss"` // 0.0..1.0
	TTL        uint32    `json:"ttl"`
}

// QuerySamples performs `count` queries to `server` for `name` and returns timing and TTL info.
// server should be an IP or hostname (port 53 will be used unless you include a port).
func QuerySamples(name, server string, count int, useTCP bool, timeout time.Duration) (*Result, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be > 0")
	}

	netProto := "udp"
	if useTCP {
		netProto = "tcp"
	}

	client := &dns.Client{Net: netProto, Timeout: timeout}

	addr := server
	// If server has no port, use 53
	if _, _, err := net.SplitHostPort(server); err != nil {
		addr = net.JoinHostPort(server, "53")
	}

	samples := make([]float64, 0, count)
	var lost int
	var ttl uint32

	for i := 0; i < count; i++ {
		m := new(dns.Msg)
		m.SetQuestion(dns.Fqdn(name), dns.TypeA)

		t0 := time.Now()
		resp, _, err := client.Exchange(m, addr)
		dt := time.Since(t0)

		if err != nil || resp == nil || resp.Rcode != dns.RcodeSuccess || len(resp.Answer) == 0 {
			lost++
			continue
		}

		// record ms
		samples = append(samples, float64(dt.Milliseconds()))

		// capture TTL from first answer if available
		if ttl == 0 {
			if len(resp.Answer) > 0 {
				ttl = resp.Answer[0].Header().Ttl
			}
		}
	}

	success := len(samples)
	var avgMs float64
	if success > 0 {
		var sum float64
		for _, s := range samples {
			sum += s
		}
		avgMs = sum / float64(success)
	}

	res := &Result{
		Sent:       count,
		Success:    success,
		SamplesMs:  samples,
		AvgMs:      avgMs,
		PacketLoss: float64(lost) / float64(count),
		TTL:        ttl,
	}

	return res, nil
}
