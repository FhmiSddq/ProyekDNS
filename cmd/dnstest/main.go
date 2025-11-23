package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/FhmiSddq/ProyekDNS/internal/app/dnstest"
)

func main() {
	server := flag.String("server", "127.0.0.1", "DNS server IP or host")
	domain := flag.String("domain", "example.com", "domain to query")
	count := flag.Int("count", 5, "number of queries to send")
	useTCP := flag.Bool("tcp", false, "use TCP instead of UDP")
	timeout := flag.Duration("timeout", 2*time.Second, "per-query timeout")
	flag.Parse()

	res, err := dnstest.QuerySamples(*domain, *server, *count, *useTCP, *timeout)
	if err != nil {
		log.Fatalf("test failed: %v", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(res); err != nil {
		fmt.Fprintf(os.Stderr, "failed to print result: %v\n", err)
		os.Exit(2)
	}
}
