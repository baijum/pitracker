package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"strconv"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type TokenStore struct {
	tokens map[string]string
	mu     sync.RWMutex
}

type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Item struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (s *TokenStore) Get(token string) string {
	if IsRedisStore == true {
		t, err := redis.String(RC.Do("GET", token))
		if err != nil {
			log.Println(err)
		}
		return t
	} else {
		s.mu.RLock()
		defer s.mu.RUnlock()
		log.Printf("%+v", s)
		t := s.tokens[token]
		log.Printf("Token requested: %s, %#v", token, token)
		log.Printf("Token retrieved: %s", t)
		return t
	}
}

func (s *TokenStore) Set(token, user string) bool {
	if IsRedisStore == true {
		_, err := redis.String(RC.Do("GET", token))
		if err != nil {
			RC.Do("SET", token, user)
			return true
		} else {
			return false
		}
	} else {
		s.mu.Lock()
		defer s.mu.Unlock()
		_, present := s.tokens[token]
		if present {
			return false
		}
		s.tokens[token] = user
		return true
	}
}

func (s *TokenStore) Put(username string) string {
	token := jwt.New(jwt.GetSigningMethod("HS256"))

	token.Claims["username"] = username
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, _ := token.SignedString(rsaPrivateKey)
	for {
		if s.Set(tokenString, username) {
			return tokenString
		}
	}
}

func NewTokenStore() *TokenStore {
	return &TokenStore{
		tokens: make(map[string]string),
	}
}

func Authorize(w http.ResponseWriter, r *http.Request, p string) error {
	tokenString := r.Header.Get("token")

	token := TS.Get(tokenString)
	log.Printf("TokenString1: %+v", tokenString)
	log.Printf("Token1: %+v", token)
	if token == "" {
		w.WriteHeader(401)
		return errors.New("Unauthorized")
	} else {
		w.WriteHeader(200)
		return nil
	}
}

type AuthToken struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	var realPassword string
	DB.QueryRow(`SELECT password
         FROM "member"
         WHERE username = $1`, username).Scan(&realPassword)

	if password == realPassword {
		tokenString := TS.Put(username)
		log.Println(tokenString)
		authToken, err := json.Marshal(AuthToken{true, tokenString, "Logged in"})
		if err != nil {
			log.Fatal("Unable to marhal token")
		}
		log.Printf("%s", authToken)
		w.Write([]byte(authToken))
		token := TS.Get(tokenString)
		log.Printf("TokenString2: %+v", tokenString)
		log.Printf("Token2: %+v", token)

	} else {
		authToken, err := json.Marshal(AuthToken{false, "", "Incorrect username or password"})
		if err != nil {
			log.Fatal("Unable to marhal token")
		}
		w.Write(authToken)
	}
}

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var (
		id int
	)
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

	DB.QueryRow(`INSERT INTO "project" (
         name, description) VALUES ($1, $2) RETURNING id`,
		project.Name, project.Description).Scan(&id)

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

func GetAllProjectsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		id          int
		name        string
		description string
	)
	type proj struct {
		Id int `json:"id"`
		Project
	}
	var pl []proj
	t := make(map[string][]proj)
	rows, err := DB.Query(`SELECT
         id, name, description FROM "project"`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &description)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name, description)
		p := proj{id, Project{name, description}}
		pl = append(pl, p)
	}
	w.Header().Set("Content-Type", "application/json")
	if pl == []proj(nil) {
		w.Write([]byte(`{"projects": []}`))
	} else {
		log.Printf("%+v", pl)
		t["projects"] = pl
		log.Printf("%+v", t)
		out, err := json.Marshal(t)
		if err != nil {
			log.Fatal("Unable to marhal")
		}
		log.Printf("Out: %s", out)
		w.Write(out)
	}
}

func GetProjectHandler(w http.ResponseWriter, r *http.Request) {
	var (
		id          int
		name        string
		description string
	)
	type proj struct {
		Id int `json:"id"`
		Project
	}

	vars := mux.Vars(r)
	projectId := vars["project"]

	var pl []proj
	t := make(map[string][]proj)
	rows, err := DB.Query(`SELECT
         id, name, description FROM "project"
         WHERE id = $1 AND archived = FALSE`, projectId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &description)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name, description)
		p := proj{id, Project{name, description}}
		pl = append(pl, p)
	}
	log.Printf("%+v", pl)
	t["projects"] = pl
	log.Printf("%+v", t)
	w.Header().Set("Content-Type", "application/json")
	out, err := json.Marshal(t)
	if err != nil {
		log.Fatal("Unable to marhal")
	}
	log.Printf("Out: %s", out)
	w.Write(out)
}

func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	type proj struct {
		Id int `json:"id"`
		Project
	}
	var pl []proj
	t1 := make(map[string][]proj)

	vars := mux.Vars(r)
	id := vars["project"]

	i, err := strconv.Atoi(id)
	if err != nil {
		// handle error
		log.Println(err)
		log.Fatal("Wrong ID")
	}

	decoder := json.NewDecoder(r.Body)
	var t map[string]proj
	err = decoder.Decode(&t)
	if err != nil {
		log.Fatal("Unable to decode body")
	}
	project := t["project"]
	log.Printf("id: %+v", id)
	log.Printf("Project: %+v", project)

	DB.QueryRow(`UPDATE "project"
         SET name = $1, description = $2 WHERE id = $3`,
		project.Name, project.Description, id)

	p := proj{i, Project{project.Name, project.Description}}
	pl = append(pl, p)

	t1["projects"] = pl

	out, err := json.Marshal(t1)
	if err != nil {
		log.Fatal("Unable to marhal")
	}
	log.Printf("Out: %s", out)
	w.Write(out)
}

func ArchiveProjectHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["project"]

	DB.QueryRow(`UPDATE "project"
           SET archived=TRUE WHERE id = $1`,
		id)
}


func GetAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		id          int
		title       string
		description string
	)
	type itm struct {
		Id int `json:"id"`
		Item
	}
	var il []itm
	t := make(map[string][]itm)
	rows, err := DB.Query(`SELECT
         id, title, description FROM "item"`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &title, &description)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, title, description)
		i := itm{id, Item{title, description}}
		il = append(il, i)
	}
	w.Header().Set("Content-Type", "application/json")
	if il == []itm(nil) {
		i := itm{}
		il = append(il, i)
		w.Write([]byte(`{"items": []}`))
	} else {
		log.Printf("%+v", il)
		t["items"] = il
		log.Printf("%+v", t)
		out, err := json.Marshal(t)
		if err != nil {
			log.Fatal("Unable to marhal")
		}
		log.Printf("Out: %s", out)
		w.Write(out)
	}
}

func CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	var id int

	type itm struct {
		Id int `json:"id"`
		Item
	}
	var il []itm
	t1 := make(map[string][]itm)

	decoder := json.NewDecoder(r.Body)
	var t map[string]Item
	err := decoder.Decode(&t)
	if err != nil {
		log.Fatal("Unable to decode body")
	}
	item := t["item"]
	log.Printf("Item: %+v", item)

	DB.QueryRow(`INSERT INTO "item" (
         title, description) VALUES ($1, $2) RETURNING id`,
		item.Title, item.Description).Scan(&id)

	i := itm{id, Item{item.Title, item.Description}}
	il = append(il, i)
	t1["items"] = il
	w.Header().Set("Content-Type", "application/json")
	out, err := json.Marshal(t1)
	if err != nil {
		log.Fatal("Unable to marhal")
	}
	log.Printf("Out: %s", out)
	w.Write(out)
}

func GetItemHandler(w http.ResponseWriter, r *http.Request) {
	var (
		id          string
		title       string
		description string
	)
	type itm struct {
		Id string `json:"id"`
		Item
	}

	vars := mux.Vars(r)
	id = vars["item"]

	var il []itm
	t := make(map[string][]itm)
	rows, err := DB.Query(`SELECT
         id, title, description FROM "item"
         WHERE id = $1`, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &title, &description)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, title, description)
		i := itm{id, Item{title, description}}
		il = append(il, i)
	}
	log.Printf("%+v", il)
	t["items"] = il
	log.Printf("%+v", t)
	w.Header().Set("Content-Type", "application/json")
	out, err := json.Marshal(t)
	if err != nil {
		log.Fatal("Unable to marhal")
	}
	log.Printf("Out: %s", out)
	w.Write(out)
}

func UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
}

var DB *sql.DB
var TS = NewTokenStore()
var RC redis.Conn
var IsRedisStore bool = false

func openDB() {
	var err error
	DB, err = sql.Open("postgres",
		`user=baiju
                 dbname=pitrackerlocal
                 sslmode=disable
                 password='passwd'`)

	if err != nil {
		log.Fatal(err)
	}

}

func openRedisConn() {
	var err error
	RC, err = redis.Dial("tcp", ":6379")

	if err != nil {
		log.Fatal(err)
	}
}

var rsaPrivateKey []byte

func main() {
	var err error
	rsaPrivateKey, err = ioutil.ReadFile("test/id_rsa")
	if err != nil {
		log.Fatal(err.Error())
	}

	openDB()

	defer DB.Close()

	if IsRedisStore == true {
		openRedisConn()

		defer RC.Close()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			// sig is a ^C, handle it
			DB.Close()
			os.Exit(1)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/auth", AuthHandler).Methods("POST")
	// r.HandleFunc("/api/v1/profiles", GetAllProfilesHandler).Methods("GET")
	// r.HandleFunc("/api/v1/profiles", CreateProfileHandler).Methods("POST")
	// r.HandleFunc("/api/v1/profiles/{profile}", GetProfileHandler).Methods("GET")
	r.HandleFunc("/api/v1/projects", GetAllProjectsHandler).Methods("GET")
	r.HandleFunc("/api/v1/projects", CreateProjectHandler).Methods("POST")
	r.HandleFunc("/api/v1/projects/{project}", GetProjectHandler).Methods("GET")
	r.HandleFunc("/api/v1/projects/{project}", UpdateProjectHandler).Methods("PUT")
	r.HandleFunc("/api/v1/projects/{project}", ArchiveProjectHandler).Methods("DELETE")
	r.HandleFunc("/api/v1/items", GetAllItemsHandler).Methods("GET")
	r.HandleFunc("/api/v1/items", CreateItemHandler).Methods("POST")
	r.HandleFunc("/api/v1/items/{item}", GetItemHandler).Methods("GET")
	r.HandleFunc("/api/v1/items/{item}", UpdateItemHandler).Methods("PUT")
	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":3000")
}
