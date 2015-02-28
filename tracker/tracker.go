package tracker

import (
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

var (
	DB *bolt.DB
)

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

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
}

// EmberClassic returns a new Negroni instance with the default
// middleware already in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging
// Static - Static File Serving
func EmberClassic() *negroni.Negroni {
	return negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.NewStatic(http.Dir("web/dist")))
}

func WebMain() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/projects", GetAllProjectsHandler).Methods("GET")
	// r.HandleFunc("/api/v1/projects", CreateProjectHandler).Methods("POST")
	// r.HandleFunc("/api/v1/projects/{project}", GetProjectHandler).Methods("GET")
	// r.HandleFunc("/api/v1/projects/{project}", UpdateProjectHandler).Methods("PUT")
	// r.HandleFunc("/api/v1/projects/{project}", ArchiveProjectHandler).Methods("DELETE")
	// r.HandleFunc("/api/v1/items", GetAllItemsHandler).Methods("GET")
	// r.HandleFunc("/api/v1/items", CreateItemHandler).Methods("POST")
	// r.HandleFunc("/api/v1/items/{item}", GetItemHandler).Methods("GET")
	// r.HandleFunc("/api/v1/items/{item}", UpdateItemHandler).Methods("PUT")
	n := EmberClassic()
	n.UseHandler(r)
	n.Run(":3000")
}
