package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
)

// TODO: Parameterize these constants
const (
	instance = "nijk-225007:asia-northeast1:nijk-master"
	user     = "nijk"
	preset   = "python"
)

type binDB struct {
	db   *sql.DB
	last string
	err  error
}

// Exec is a wrapper function which executes sql.DB.Exec and keep its error,
// when only there was no error before.
func (bdb *binDB) Exec(query string, args ...interface{}) {
	if bdb.err == nil {
		_, err := bdb.db.Exec(query, args...)
		bdb.last = query
		bdb.err = err
	}
}

func main() {
	cfg := mysql.Cfg(instance, user, "")
	cfg.DBName = "nijk"
	db, err := mysql.DialCfg(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bdb := binDB{db, "", nil}
	for _, relation := range []string{"paradigmatic", "syntagmatic"} {
		bdb.Exec(dropTable(preset, relation))
		bdb.Exec(createTable(preset, relation))
	}
}

func dropTable(preset, relation string) string {
	template := `
DROP TABLE IF EXISTS %s_%s`
	return fmt.Sprintf(template, preset, relation)
}

func createTable(preset, relation string) string {
	template := `
CREATE TABLE %s_%s
(
  this   VARCHAR(64),
  that   VARCHAR(64),
  score  DOUBLE,
  PRIMARY KEY (this, that),
  INDEX (score)
);`
	return fmt.Sprintf(template, preset, relation)
}
