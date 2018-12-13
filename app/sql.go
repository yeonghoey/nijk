package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// TODO: Parameterize these constants
const (
	instance = "nijk-225007:asia-northeast1:nijk-master"
	user     = "nijk"
	topN     = 100
)

func openDBDev() *sql.DB {
	cfg := mysql.Cfg(instance, user, "")
	cfg.DBName = "nijk"
	db, err := mysql.DialCfg(cfg)
	if err != nil {
		log.Fatalf("Could not open db: %v", err)
	}
	return db
}

func openDBProd() *sql.DB {
	dataSourceName := fmt.Sprintf("%s:@cloudsql(%s)/nijk", user, instance)

	var err error
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("Could not open db: %v", err)
	}
	return db
}

func queryParadigmatic(preset, this string) []string {
	q := fmt.Sprintf(`
SELECT that
FROM %s_paradigmatic
WHERE this=?
ORDER BY score DESC
LIMIT %d;`, preset, topN)
	return query(q, this)
}

func querySyntagmatic(preset, this string) []string {
	q := fmt.Sprintf(`
SELECT that
FROM  %s_syntagmatic
WHERE this=?
ORDER BY score DESC
LIMIT %d;`, preset, topN)
	return query(q, this)
}

func queryContaining(preset, this string) []string {
	q := fmt.Sprintf(`
SELECT DISTINCT this
FROM %s_paradigmatic
WHERE this<>? AND this LIKE ?
LIMIT %d;`, preset, topN)
	return query(q, this, "%"+this+"%")
}

func query(q string, args ...interface{}) []string {
	rows, err := db.Query(q, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()

	thats := []string{}
	for rows.Next() {
		var that string
		if err := rows.Scan(&that); err != nil {
			continue
		}
		thats = append(thats, that)
	}
	return thats
}
