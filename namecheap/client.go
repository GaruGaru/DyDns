package namecheap

import (
	"fmt"
	"net/http"
	"time"
)

type NamecheapOptions struct {
	Host     string
	Domains  []string
	Password string
}

type DnsUpdateResult struct {
	Host    string
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

func createUpdateRequest(host string, domain string, password string, ip string) string {
	return fmt.Sprintf(
		"https://dynamicdns.park-your-domain.com/update?host=%s&domain=%s&password=%s&ip=%s",
		host,
		domain,
		password,
		ip,
	)
}

func (c DnsClient) Update(options NamecheapOptions, ip string) ([]DnsUpdateResult) {

	results := make([]DnsUpdateResult, len(options.Domains))

	for _, domain := range options.Domains {

		url := createUpdateRequest(options.Host, domain, options.Password, ip)

		resp, err := c.httpClient.Get(url)

		if err != nil {
			results = append(results, DnsUpdateResult{
				Host:    options.Host,
				Domain:  domain,
				IP:      ip,
				Status:  err.Error(),
				Success: false,
			})
		}

		if resp.StatusCode != 200 {
			results = append(results, DnsUpdateResult{
				Host:    options.Host,
				Domain:  domain,
				Status:  fmt.Sprintf("Unexpected status code %d", resp.StatusCode),
				Success: false,
			})
		}

		resp.Body.Close()
	}

	return results
}
