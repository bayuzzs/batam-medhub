package main

import (
	"log"

	"batam-medhub/providers/internal/providerapp"
)

func main() {
	if err := providerapp.Run(providerapp.Identity{ID: "hotel-demo-01", Type: "HOTEL"}); err != nil {
		log.Fatal(err)
	}
}
