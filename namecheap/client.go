package namecheap

import (
	"fmt"
	"net/http"
	"time"
)

type DnsOptions struct {
	Host     string
	Domain   string
	Password string
	IP       string
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

func createUpdateRequest(options DnsOptions) string {
	return fmt.Sprintf(
		"https://dynamicdns.park-your-domain.com/update?host=%s&domain=%s&password=%s&ip=%s",
		options.Host,
		options.Domain,
		options.Password,
		options.IP,
	)
}

func (c DnsClient) Update(options DnsOptions) error {

	url := createUpdateRequest(options)

	resp, err := c.httpClient.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return nil
}
