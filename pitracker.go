package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/baijum/plus"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	signingKey   = "SOME KEY"
	privateKey   []byte
	publicKey    []byte
	clientID     string
	clientSecret string
)

type AuthToken struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Item struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Name struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
}

type Email struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Profile struct {
	DisplayName string  `json:"displayName"`
	Name        Name    `json:"name"`
	Emails      []Email `json:"emails"`
	Gender      string  `json:"gender"`
	URL         string  `json:"url"`
}

func AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err == nil && token.Valid {
		next(w, r)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	}
}

/*
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
*/

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	x, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error reading code in the request body: ", err)
	}
	code := string(x)

	accessToken, idToken, err := plus.GetTokens(code, clientID, clientSecret)
	if err != nil {
		log.Fatal("error exchanging code for access token: ", err)
	}
	gplusID, err := plus.DecodeIDToken(idToken)
	if err != nil {
		log.Fatal("Error decoding ID token: ", err)
	}

	baseUrl := "https://www.googleapis.com/plus/v1"
	apiFull := fmt.Sprintf("%s/people/%s", baseUrl, gplusID)

	url := fmt.Sprintf("%s?access_token=%s", apiFull, accessToken)

	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err, "Error", 500)
	}

	defer resp.Body.Close()

	var profile Profile

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&profile)

	if err != nil {
		log.Fatal("unable to decode body: ", err)
	}
	fmt.Printf("profile: %v\n", profile)
	fmt.Printf("gplusID: %s\n", gplusID)
	var displayName string
	var id int
	DB.QueryRow(`SELECT displayname
		from "member"
		WHERE plusid = $1`, gplusID).Scan(&displayName)

	if displayName == "" {
		fmt.Printf("displayName: %s\n", displayName)
		DB.QueryRow(`INSERT INTO "member" (
			plusid, email, displayname,
			familyname, givenname,
			gender, url) VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id`,
			gplusID, profile.Emails[0].Value, profile.DisplayName,
			profile.Name.FamilyName, profile.Name.GivenName,
			profile.Gender, profile.URL).Scan(&id)

	}
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["sub"] = gplusID
	token.Claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	tokenString, _ := token.SignedString(privateKey)
	fmt.Printf("gplusID: %v\n", gplusID)
	fmt.Printf("tokenString: %v\n", tokenString)

	authToken, err := json.Marshal(AuthToken{true, tokenString, "Logged in"})
	if err != nil {
		log.Fatal("Unable to marhal token")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(authToken))
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

func GetAllProjectsHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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

var rsaPrivateKey []byte

func init() {
	privateKey, _ = ioutil.ReadFile("test/id_rsa")
	publicKey, _ = ioutil.ReadFile("test/id_rsa.pub")
	clientID = os.Getenv("GPLUS_CLIENT_ID")
	clientSecret = os.Getenv("GPLUS_CLIENT_SECRET")
}

func main() {
	var err error
	rsaPrivateKey, err = ioutil.ReadFile("test/id_rsa")
	if err != nil {
		log.Fatal(err.Error())
	}

	openDB()

	defer DB.Close()

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
	r.Handle("/api/v1/projects",
		negroni.New(negroni.HandlerFunc(AuthMiddleware),
			negroni.HandlerFunc(GetAllProjectsHandler))).Methods("GET")
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
