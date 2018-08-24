package namecheap

import (
	"fmt"
	"net/http"
	"time"
)

type NamecheapOptions struct {
	Domain   string
	Entries  []string
	Password string
}

type DnsUpdateResult struct {
	Entry   string
	Domain  string
	IP      string
	Status  string
	Success bool
}

func NewDnsClient() DnsClient {
	return DnsClient{
		httpClient: http.Client{
			Timeout: time.Second * 10,
		},
	}
}

type DnsClient struct {
	httpClient http.Client
}

func createUpdateRequest(domain string, entry string, password string, ip string) string {
	return fmt.Sprintf(
		"https://dynamicdns.park-your-domain.com/update?domain=%s&host=%s&password=%s&ip=%s",
		domain,
		entry,
		password,
		ip,
	)
}

func (c DnsClient) Update(options NamecheapOptions, ip string) ([]DnsUpdateResult) {

	results := make([]DnsUpdateResult, len(options.Entries))

	for i, entry := range options.Entries {

		url := createUpdateRequest(options.Domain, entry, options.Password, ip)

		resp, err := c.httpClient.Get(url)

		if err != nil {
			results[i] = DnsUpdateResult{
				Entry:   entry,
				Domain:  options.Domain,
				IP:      ip,
				Status:  err.Error(),
				Success: false,
			}
			continue
		}

		if resp.StatusCode != 200 {
			results[i] = DnsUpdateResult{
				Entry:   entry,
				Domain:  options.Domain,
				Status:  fmt.Sprintf("Unexpected status code %d", resp.StatusCode),
				Success: false,
			}
		} else {
			results[i] = DnsUpdateResult{
				Entry:   entry,
				Domain:  options.Domain,
				Status:  "OK",
				Success: true,
				IP:      ip,
			}
		}

		resp.Body.Close()
	}

	return results
}
