package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"

)

var host string
var port int
var user string
var password string
var dbname string

type Entry struct {
	UID  int    `json:"uid"`
	Test string `json:"test"`
}

func main() {
	host = os.Getenv("DBINSTANCE_HOST")
	port, _ = strconv.Atoi(os.Getenv("DBINSTANCE_PORT"))
	user = os.Getenv("DBINSTANCE_USERNAME")
	password = os.Getenv("DBINSTANCE_PASSWORD")
	dbname = os.Getenv("DBINSTANCE_DATABASE")

	//query := `CREATE TABLE IF NOT EXISTS myTable(uid int, test text)`
	//
	//ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancelfunc()
	//db := OpenConnection()
	//
	//_, err := db.ExecContext(ctx, query)
	//if err != nil {
	//	log.Printf("Error %s when creating product table", err)
	//
	//}

	http.HandleFunc("/", GETHandler)
	http.HandleFunc("/insert", POSTHandler)
	http.HandleFunc("/delete/", DELETEHandler)
	http.HandleFunc("/update", PUTHandler)
	fmt.Println("Server is running...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func GETHandler(w http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(w, "<h1>AWS RDS Database</h1>")
	db := OpenConnection()

	rows, err := db.Query(`SELECT "uid", "test" FROM "myTable"`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Retrieved entries from table:myTable in database ack-rds-test-postgre")
	var entries []Entry

	for rows.Next() {

		var entry Entry
		rows.Scan(&entry.UID, &entry.Test)
		fmt.Printf("{uid: %d, test: %s}\n", entry.UID, entry.Test)
		entries = append(entries, entry)
	}

	entryBytes, _ := json.MarshalIndent(entries, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(entryBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var entry Entry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	insertDynStmt := `insert into "myTable"("uid", "test") values($1, $2)`
	_, err = db.Exec(insertDynStmt, entry.UID, entry.Test)
	fmt.Printf("Inserting entry to table:myTable with values {uid: %d, test: %s\n", entry.UID, entry.Test)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}
func DELETEHandler(w http.ResponseWriter, r *http.Request) {
	uids, ok := r.URL.Query()["uid"]

	if !ok || len(uids[0]) < 1 {
		log.Println("Url Param 'uid' is missing")
		return
	}
	uid := uids[0]

	log.Println("Url Param 'uid' is: " + string(uid))
	db := OpenConnection()
	fmt.Printf("Deleting entry from table:myTable where uid=%d in database ack-rds-test-postgre\n", uid)
	deleteStmt := `delete from "myTable" where "uid"=$1`
	_, err := db.Exec(deleteStmt, uid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func PUTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var entry Entry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("Updating entry from table:myTable where uid=%d to {uid: %d, test:%s}\n", entry.UID, entry.UID, entry.Test)
	updateStmt := `update "myTable" set "test"=$1 where "uid"=$2`
	_, err = db.Exec(updateStmt, entry.Test, entry.UID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func OpenConnection() *sql.DB {
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname)

	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully to database:ack-rds-test-postgre on endpoint %s\n", host)

	return db
}
