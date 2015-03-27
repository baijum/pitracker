// PITracker - Project and Issue Tracker
// Copyright (C) 2015 Baiju Muthukadan <baiju@muthukadan.net>

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
