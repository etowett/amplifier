package helpers

import (
	"log"
	"time"
)

var (
	LocKenya *time.Location
)

func init() {
	var err error
	LocKenya, err = time.LoadLocation("Africa/Nairobi")
	if err != nil {
		log.Fatal("Failed to load location Africa/Nairobi")
	}
}
