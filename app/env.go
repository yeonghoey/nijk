package main

import (
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
)

var isDev = os.Getenv("DEV") != ""

var db *sql.DB

func init() {
	if isDev {
		db = openDBDev()
	} else {
		db = openDBProd()
		gin.SetMode(gin.ReleaseMode)
	}
}
