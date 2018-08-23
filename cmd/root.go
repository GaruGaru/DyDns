package cmd

import (
	"fmt"
	"github.com/GaruGaru/DyDns/ip"
	"github.com/GaruGaru/DyDns/namecheap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
	"strings"
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

		ipProvider := ip.Providers(
			ip.NewPlainIPProvider("https://api.ipify.org/"),
			ip.NewPlainIPProvider("http://myexternalip.com/raw"),
		)

		options := namecheap.NamecheapOptions{
			Host:     viper.GetString("host"),
			Domains:  strings.Split(viper.GetString("domain"), ","),
			Password: viper.GetString("password"),
		}

		dnsClient := namecheap.NewDnsClient()

		for {

			externalIP, err := ipProvider.ExternalIP()

			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			updateResults := dnsClient.Update(options, externalIP)

			if len(updateResults) == 0{
				fmt.Println("No update results.")
			}

			for _, result := range updateResults {
				if result.Success {
					fmt.Printf("[OK] Updated %s (%s): %s", result.Host, result.Domain, result.IP)
				} else {
					fmt.Printf("Error updating %s (%s): %s", result.Host, result.Domain, result.Status)
				}
			}

			time.Sleep(time.Duration(viper.GetInt("delay")) * time.Minute)
		}

	},
}

func Execute() error {
	return rootCmd.Execute()
}
