package main

import (
	"log"

	"batam-medhub/providers/internal/providerapp"
)

func main() {
	if err := providerapp.Run(providerapp.Identity{ID: "transport-demo-01", Type: "TRANSPORT"}); err != nil {
		log.Fatal(err)
	}
}
