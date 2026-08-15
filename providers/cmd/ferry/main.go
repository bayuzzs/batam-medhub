package main

import (
	"log"

	"batam-medhub/providers/internal/providerapp"
)

func main() {
	if err := providerapp.Run(providerapp.Identity{ID: "ferry-demo-01", Type: "FERRY"}); err != nil {
		log.Fatal(err)
	}
}
