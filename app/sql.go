package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// TODO: Parameterize these constants
const (
	instance = "nijk-225007:asia-northeast1:nijk-master"
	user     = "nijk"
)

func openDB() *sql.DB {
	s := os.Getenv("DEV")
	isDev := s != ""
	if isDev {
		return openDBDev()
	}
	return openDBProd()
}

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
