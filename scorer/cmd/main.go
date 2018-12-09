package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	"github.com/yeonghoey/nijk/scorer"
)

const (
	k = 1.2
	b = 0.75
)

// TODO: Deduplicate queries
const paradigSchema = `
CREATE TABLE IF NOT EXISTS paradigmatic
(
  preset VARCHAR(64),
  this   VARCHAR(64),
  that   VARCHAR(64),
  score  DOUBLE,
  PRIMARY KEY (preset, this, that)
);`

const paradigInsert = `
INSERT INTO paradigmatic (preset, this, that, score) VALUES (?, ?, ?, ?);`

const syntagSchema = `
CREATE TABLE IF NOT EXISTS syntagmatic
(
  preset VARCHAR(64),
  this   VARCHAR(64),
  that   VARCHAR(64),
  score  DOUBLE,
  PRIMARY KEY (preset, this, that)
);`

const syntagInsert = `
INSERT INTO syntagmatic (preset, this, that, score) VALUES (?, ?, ?, ?);`

type dbWrapper struct {
	db   *sql.DB
	last string
	err  error
}

func (dbw *dbWrapper) exec(query string, args ...interface{}) {
	if dbw.err == nil {
		_, err := dbw.db.Exec(query, args...)
		dbw.last = query
		dbw.err = err
	}
}

func main() {
	db, err := mysql.Dial("nijk-225007:asia-northeast1:nijk-master", "nijk")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbw := dbWrapper{db, "", nil}
	dbw.exec("USE nijk;")
	dbw.exec(paradigSchema)
	dbw.exec(syntagSchema)
	if dbw.err != nil {
		log.Fatalf("%s: %v\n", dbw.last, dbw.err)
	}

	reader := bufio.NewReader(os.Stdin)
	collection := scorer.NewCollection(reader, k, b)

	collection.Paradigmatic(func(a, b string, score float64) {
		fmt.Printf("Paradigmatic: %s %s %.2f\n", a, b, score)
		db.Exec(paradigInsert, "python", a, b, score)
	})

	collection.Syntagmatic(func(a, b string, score float64) {
		fmt.Printf("Syntagmatic: %s %s %.2f\n", a, b, score)
		db.Exec(syntagInsert, "python", a, b, score)
	})
}
