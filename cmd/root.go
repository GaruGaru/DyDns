package cmd

import (
	"fmt"
	"github.com/GaruGaru/DyDns/ip"
	"github.com/GaruGaru/DyDns/namecheap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

func init() {
	rootCmd.Flags().String("host", "", "Dns entry name")
	rootCmd.Flags().String("domain", "", "Domain name")
	rootCmd.Flags().String("password", "", "Dynamic dns password")
	rootCmd.Flags().Int("delay", 60, "Dynamic dns password")
	viper.BindPFlags(rootCmd.Flags())
	viper.AutomaticEnv()
	validateFields()
}

func validateFields() {
	mandatoryFields := []string{"host", "domain", "password"}
	for _, field := range mandatoryFields {
		if viper.GetString(field) == "" {
			panic(fmt.Sprintf("Error field %s not provided or empty", field))
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "DyDns",
	Short: "DyDns is a lightweight dynamic dns client",
	Run: func(cmd *cobra.Command, args []string) {

		for {

			ipProvider := ip.Providers(
				ip.NewPlainIPProvider("https://api.ipify.org/"),
				ip.NewPlainIPProvider("http://myexternalip.com/raw"),
			)

			externalIP, err := ipProvider.ExternalIP()

			if err != nil {
				fmt.Println(err.Error())
			}

			options := namecheap.DnsOptions{
				Host:     viper.GetString("host"),
				Domain:   viper.GetString("domain"),
				Password: viper.GetString("password"),
				IP:       externalIP,
			}

			dnsClient := namecheap.NewDnsClient()

			updateErr := dnsClient.Update(options)

			if updateErr == nil {
				fmt.Printf("[OK] %s.%s -> %s\n", options.Host, options.Domain, options.IP)
			} else {
				fmt.Println(updateErr.Error())
			}

			time.Sleep(time.Duration(viper.GetInt("delay")) * time.Minute)
		}

	},
}

func Execute() error {
	return rootCmd.Execute()
}
