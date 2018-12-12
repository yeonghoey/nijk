package main

import (
	"database/sql"
	"fmt"
	"log"

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
	r.GET("/:term", func(c *gin.Context) {
		term := c.Param("term")
		c.JSON(200, gin.H{
			"par": parTerms(term),
			"syn": synTerms(term),
			"sim": simTerms(term),
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

func simTerms(this string) []string {
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
