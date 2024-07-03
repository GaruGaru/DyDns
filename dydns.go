package main

import (
	"context"
	"github.com/garugaru/DyDns/ip"
	"github.com/garugaru/DyDns/namecheap"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Ctx(ctx).
		Timestamp().
		Logger()

	var ipProvider = ip.Providers(
		ip.NewPlainIPProvider("https://api.ipify.org/"),
		ip.NewPlainIPProvider("http://myexternalip.com/raw"),
	)

	var options = namecheap.Options{
		Domain:   os.Getenv("DOMAIN"),
		Entries:  strings.Split(os.Getenv("ENTRIES"), ","),
		Password: os.Getenv("PASSWORD"),
	}

	if len(options.Password) == 0 {
		logger.Fatal().Msg("password is required")
	}

	if len(options.Entries) == 0 {
		logger.Fatal().Msg("atleast 1 entry is required")
	}

	if len(options.Domain) == 0 {
		logger.Fatal().Msg("domain is required")
	}

	delay := 60 * time.Second

	delayEnv := os.Getenv("DELAY")
	var err error
	if len(delayEnv) != 0 {
		delay, err = time.ParseDuration(delayEnv)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to parse delay")
		}
	}

	dnsClient := namecheap.NewDnsClient()
	logger.Info().Msgf("Starting dydns on domain %s with %d entries", options.Domain, len(options.Entries))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		logger.Warn().Str("signal", sig.String()).Msg("received signal")
		cancel()
		done <- true
	}()

	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info().Msg("context done, exiting")
				return
			case <-ticker.C:
				externalIP, err := ipProvider.IP(ctx)
				if err != nil {
					logger.Warn().Err(err).Msg("error retrieving external IP")
					continue
				}

				logger.Info().Str("ip", externalIP).Msg("got IP")
				err = dnsClient.Update(ctx, options, externalIP)
				if err != nil {
					logger.Warn().Err(err).Msg("error updating DNS record")
					continue
				}
			}
		}
	}()

	<-done
	logger.Info().Msg("dydns exiting")
}
