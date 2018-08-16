package ip

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func Providers(providers ...Provider) Provider {
	return ProvidersManager{
		Providers: providers,
	}
}

type ProvidersManager struct {
	Providers []Provider
}

func (p ProvidersManager) ExternalIP() (string, error) {

	for _, provider := range p.Providers {
		ip, err := provider.ExternalIP()
		if err == nil {
			return ip, nil
		} else {
			fmt.Printf("provider %s not working: %s\n", provider.Name(), err.Error())
		}
	}

	return "", fmt.Errorf("no working provider")
}

func (p ProvidersManager) Name() string {
	return "providers-manager"
}

type Provider interface {
	ExternalIP() (string, error)
	Name() string
}

func NewPlainIPProvider(url string) Provider {
	return PlainIPProvider{
		client: http.Client{
			Timeout: time.Second * 10,
		},
		url: url,
	}
}

type PlainIPProvider struct {
	url    string
	client http.Client
}

func (p PlainIPProvider) ExternalIP() (string, error) {

	resp, err := p.client.Get(p.url)

	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	bodyString := string(bodyBytes)

	return bodyString, nil
}

func (p PlainIPProvider) Name() string {
	return p.url
}
