package main_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/baijum/pitracker/tracker"
	"github.com/boltdb/bolt"
)

func tempfile() string {
	f, _ := ioutil.TempFile("", "pitracker-bolt-")
	f.Close()
	os.Remove(f.Name())
	return f.Name()
}

func TestBoltDBFileAvailability(t *testing.T) {
	file := tempfile()
	db1, _ := tracker.OpenBoltDB(file)
	db1.Close()
	db2, _ := tracker.OpenBoltDB(file)
	db2.Close()

	db3, _ := tracker.OpenBoltDB(file)
	db4, err4 := tracker.OpenBoltDB(file)
	if err4 == nil {
		t.Log("DB connection is working: ", err4)
		db4.Close()
		os.Remove(file)
		t.Fail()
	}
	db3.Close()
	os.Remove(file)
}

func TestBucketCreation(t *testing.T) {
	file := tempfile()
	db, _ := tracker.OpenBoltDB(file)
	err := tracker.CreateBucket(db, "items")

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("items"))
		if b == nil {
			t.Log("Unable to create 'items' bucket: ", err)
			db.Close()
			os.Remove(file)
			t.Fail()
		}
		return nil
	})
	db.Close()
	os.Remove(file)
}

func TestReq(t *testing.T) {
	file := tempfile()
	db, _ := tracker.OpenBoltDB(file)

	req, err := http.NewRequest("GET", "/projects", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	tracker.GetAllProjectsHandler(w, req)

	exp := `{"projects": []}`
	out := w.Body.String()

	if out != exp {
		t.Fail()
	}

	db.Close()
	os.Remove(file)
}
