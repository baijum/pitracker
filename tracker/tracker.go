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

package tracker

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

var (
	DB    *bolt.DB
	uiDir string
)

type Project struct {
	Id          string `json:"_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Open Bold DB and return reference and error if any
func OpenBoltDB(file string) (*bolt.DB, error) {
	// Open the Bolt data file. It will be created if it doesn't exist.
	// timeout option prevent an indefinite wait for DB file availability
	DB, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 1 * time.Second})
	return DB, err
}

func CreateBucket(db *bolt.DB, name string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func GetAllProjectsHandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	w.Header().Set("Content-Type", "application/json")

	var pl []Project

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("project"))
		b.ForEach(func(k, v []byte) error {
			log.Printf("key=%s, value=%s\n", k, v)
			var p Project
			json.Unmarshal(v, &p)
			pl = append(pl, p)
			return nil
		})
		return nil
	})

	t := make(map[string][]Project)
	t["projects"] = pl
	out, err := json.Marshal(t)
	if err != nil {
		log.Fatal("Unable to marhal")
	}
	log.Printf("Out: %s", out)
	w.Write(out)
}

func CreateProjectHandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	var pl []Project
	t1 := make(map[string][]Project)

	decoder := json.NewDecoder(r.Body)
	var t map[string]Project
	err := decoder.Decode(&t)
	if err != nil {
		log.Fatal("Unable to decode body")
	}
	project := t["project"]
	log.Printf("Project: %+v", project)

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("project"))
		v := b.Get([]byte(project.Id))
		if v != nil {
			return errors.New("Project already exists")
		}
		return nil
	})

	if err != nil {
		log.Printf("Error: %v", err)
		// Return the struct or empty ?
		w.WriteHeader(409)
		return
	}

	p1, _ := json.Marshal(project)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("project"))
		err := b.Put([]byte(project.Id), []byte(p1))
		return err
	})

	pl = append(pl, project)

	t1["projects"] = pl

	out, err := json.Marshal(t1)
	if err != nil {
		log.Fatal("Unable to marhal")
	}
	log.Printf("Out: %s", out)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// EmberClassic returns a new Negroni instance with the default
// middleware already in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging
// Static - Static File Serving
func EmberClassic(dir string) *negroni.Negroni {
	return negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.NewStatic(http.Dir(dir)))
}

func GetProjectHandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	vars := mux.Vars(r)
	pn := vars["project"]

	t := make(map[string]Project)

	_ = db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("project"))
		v1 := c.Get([]byte(pn))
		var p Project
		json.Unmarshal(v1, &p)
		t["project"] = p
		return nil
	})

	out, err := json.Marshal(t)
	if err != nil {
		log.Fatal("Unable to marhal")
	}
	_ = out

	log.Printf("%s", out)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func UpdateProjectHandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	vars := mux.Vars(r)
	pd := vars["project"]

	var pl []Project

	var t map[string]Project
	//t = make(map[string][]Project)
	t1 := make(map[string][]Project)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)
	if err != nil {
		log.Fatal("Unable to decode body")
	}
	project := t["project"]
	project.Id = pd
	log.Printf("Project: %#v\n pd: %+v", project, pd)

	p1, _ := json.Marshal(project)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("project"))
		err := b.Put([]byte(project.Id), []byte(p1))
		return err
	})

	pl = append(pl, project)

	t1["projects"] = pl

	out, err := json.Marshal(t1)
	if err != nil {
		log.Fatal("Unable to marhal")
	}
	log.Printf("Out: %s", out)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

	// err = db.Update(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte("project"))
	// 	err := b.Put([]byte(project.Id), []byte(p1))
	// 	c := tx.Bucket([]byte("project_id_name"))
	// 	_ = c.Put([]byte(project.Id), []byte(project.Id))
	// 	return err
	// })

	// out, err := json.Marshal(t)
	// if err != nil {
	// 	log.Fatal("Unable to marhal")
	// }
	// _ = out

	// log.Printf("%s", out)
	// w.Header().Set("Content-Type", "application/json")
	// w.Write(out)

	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(`{"project": { "id": 1, "name": "ok", "description": "Okay"}}`))
}

func init() {
	uiDir = os.Getenv("PITRACKER_UI_DIR")
}

func WebMain(db *bolt.DB) {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/projects", func(w http.ResponseWriter, r *http.Request) {
		GetAllProjectsHandler(w, r, db)
	}).Methods("GET")
	r.HandleFunc("/api/v1/projects", func(w http.ResponseWriter, r *http.Request) {
		CreateProjectHandler(w, r, db)
	}).Methods("POST")
	r.HandleFunc("/api/v1/projects/{project}", func(w http.ResponseWriter, r *http.Request) {
		GetProjectHandler(w, r, db)
	}).Methods("GET")
	r.HandleFunc("/api/v1/projects/{project}", func(w http.ResponseWriter, r *http.Request) {
		UpdateProjectHandler(w, r, db)
	}).Methods("PUT")
	// r.HandleFunc("/api/v1/projects/{project}", ArchiveProjectHandler).Methods("DELETE")
	// r.HandleFunc("/api/v1/items", GetAllItemsHandler).Methods("GET")
	// r.HandleFunc("/api/v1/items", CreateItemHandler).Methods("POST")
	// r.HandleFunc("/api/v1/items/{item}", GetItemHandler).Methods("GET")
	// r.HandleFunc("/api/v1/items/{item}", UpdateItemHandler).Methods("PUT")

	//h := http.Handle("/introspect", http.StripPrefix("/introspect", boltd.NewHandler(mydb)))

	n := EmberClassic(uiDir)
	n.UseHandler(r)
	n.Run(":3000")
}

func init() {

}
