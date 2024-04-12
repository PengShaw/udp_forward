package main

import (
	"log"

	"github.com/PengShaw/udp_forward/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
