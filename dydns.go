package main

import (
	"github.com/GaruGaru/DyDns/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
		panic(err)
	}
}
