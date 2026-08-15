package main

import (
	"log"

	"batam-medhub/providers/internal/providerapp"
)

func main() {
	if err := providerapp.Run(providerapp.Identity{ID: "hospital-demo-01", Type: "HOSPITAL"}); err != nil {
		log.Fatal(err)
	}
}
