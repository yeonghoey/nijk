package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var db *sql.DB
var parStmt *sql.Stmt
var synStmt *sql.Stmt
var simStmt *sql.Stmt

// TODO: Parameterize
const (
	preset = "python"
	topN   = 100
)

func init() {
	db = openDB()

	var err error
	parStmt, err = db.Prepare(fmt.Sprintf(`
SELECT that FROM %s_paradigmatic WHERE this=? ORDER BY score DESC LIMIT %d;`,
		preset, topN))

	if err != nil {

		log.Fatalf("Could not prepare stmt: %v", err)
	}

	synStmt, err = db.Prepare(fmt.Sprintf(`
SELECT that FROM %s_syntagmatic WHERE this=? ORDER BY score DESC LIMIT %d;`,
		preset, topN))
	if err != nil {
		log.Fatalf("Could not prepare stmt: %v", err)
	}

	simStmt, err = db.Prepare(fmt.Sprintf(`
SELECT DISTINCT this FROM %s_paradigmatic WHERE this<>? AND this LIKE ? LIMIT %d;`,
		preset, topN))

	if err != nil {
		log.Fatalf("Could not prepare stmt: %v", err)
	}
}

func main() {
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"title": strings.Title,
	})

	r.LoadHTMLGlob("templates/*")
	r.GET("/:preset/:this", func(c *gin.Context) {
		preset := c.Param("preset")
		this := c.Param("this")
		c.HTML(http.StatusOK, "term.tmpl", gin.H{
			"preset":   preset,
			"this":     this,
			"parTerms": parTerms(this),
			"synTerms": synTerms(this),
			"conTerms": conTerms(this),
		})
	})
	r.Run()
}

func parTerms(this string) []string {
	return query(parStmt, this)
}

func synTerms(this string) []string {
	return query(synStmt, this)
}

func conTerms(this string) []string {
	return query(simStmt, this, "%"+this+"%")
}

func query(stmt *sql.Stmt, args ...interface{}) []string {
	rows, err := stmt.Query(args...)
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
