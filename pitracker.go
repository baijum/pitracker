package main

import (
	"log"
	"os"

	"github.com/baijum/pitracker/tracker"
)

var (
	boltDBFile string
)

func init() {
	boltDBFile = os.Getenv("PITRACKER_BOLTDB_FILE")
}

func main() {
	db, err := tracker.OpenBoltDB(boltDBFile)
	if err != nil {
		log.Fatal("Failed to connect Bolt DB: ", err)
	}
	defer db.Close()

	err = tracker.CreateBucket(db, "project")
	if err != nil {
		log.Fatal("Unable to create 'project' bucket: ", err)
	}

	err = tracker.CreateBucket(db, "item")
	if err != nil {
		log.Fatal("Unable to create 'item' bucket: ", err)
	}

	tracker.WebMain(db)
}
