package tracker

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

var (
	DB *bolt.DB
)

type Project struct {
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

func GetAllProjectsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello someting"))
}

func CreateProjectHandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	var id int
	type proj struct {
		Id int `json:"id"`
		Project
	}
	var pl []proj
	t1 := make(map[string][]proj)

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
		v := b.Get([]byte(project.Name))
		if v != nil {
			return errors.New("Project already exists")
		}
		return nil
	})

	if err != nil {
		log.Printf("Error: %v", err)
		// Return the struct or empty ?
		return
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("project"))
		err := b.Put([]byte(project.Name), []byte(project.Description))
		return err
	})

	p := proj{id, Project{project.Name, project.Description}}
	pl = append(pl, p)

	t1["projects"] = pl

	out, err := json.Marshal(t1)
	if err != nil {
		log.Fatal("Unable to marhal")
	}
	log.Printf("Out: %s", out)
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

func WebMain(db *bolt.DB) {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/projects", GetAllProjectsHandler).Methods("GET")
	r.HandleFunc("/api/v1/projects", func(w http.ResponseWriter, r *http.Request) {
		CreateProjectHandler(w, r, db)
	}).Methods("POST")
	// r.HandleFunc("/api/v1/projects/{project}", GetProjectHandler).Methods("GET")
	// r.HandleFunc("/api/v1/projects/{project}", UpdateProjectHandler).Methods("PUT")
	// r.HandleFunc("/api/v1/projects/{project}", ArchiveProjectHandler).Methods("DELETE")
	// r.HandleFunc("/api/v1/items", GetAllItemsHandler).Methods("GET")
	// r.HandleFunc("/api/v1/items", CreateItemHandler).Methods("POST")
	// r.HandleFunc("/api/v1/items/{item}", GetItemHandler).Methods("GET")
	// r.HandleFunc("/api/v1/items/{item}", UpdateItemHandler).Methods("PUT")
	n := EmberClassic("web/dist")
	n.UseHandler(r)
	n.Run(":3000")
}

func init() {

}
