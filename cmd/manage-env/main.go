package main

import (
	"context"
	"flag"
	"log"

	"github.com/shaural/GlobePersonalWebsite/server/pkg/db"
)

func main() {
	log.SetPrefix("manage-env:: ")
	initialize := flag.Bool("initialize", false, "set to initialize postgres db")
	load := flag.Bool("load", false, "load db with countries from csv")
	flag.Parse()

	loadDB, err := db.NewDatabase(context.Background())
	if err != nil {
		log.Fatalf("Unable to establish connection to postgres database Error: %v", err)
	}
	defer loadDB.Close()

	if *initialize {
		if err = loadDB.Initialize(); err != nil {
			log.Fatalf("Encountered error initializing. Rolled back any changes. Error: %v", err)
		}
	}
	if *load {
		if err = loadCountries(loadDB); err != nil {
			log.Fatalf("An error occured while loading countries and states: %v", err)
		}
	}
}
