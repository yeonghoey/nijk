package main

import (
	"bufio"
	"log"
	"os"
	"sync/atomic"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	"github.com/yeonghoey/nijk/scorer"
)

// TODO: Parameterize these constants
const (
	instance = "nijk-225007:asia-northeast1:nijk-master"
	user     = "nijk"

	numWorkers = 128

	k = 1.2
	b = 0.75

	paradigmaticThreshold = 0.75
	syntagmaticThreshold  = 0.5
)

func main() {
	cfg := mysql.Cfg(instance, user, "")
	cfg.DBName = "nijk"
	db, err := mysql.DialCfg(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize
	bdb := binDB{db, "", nil}

	// TODO: Parameterize the preset name
	bdb.Exec(queryCreateTable("python", "paradigmatic"))
	bdb.Exec(queryCreateTable("python", "syntagmatic"))

	pStmt := bdb.Prepare(queryInsert("python", "paradigmatic"))
	defer pStmt.Close()
	sStmt := bdb.Prepare(queryInsert("python", "syntagmatic"))
	defer sStmt.Close()

	if bdb.err != nil {
		log.Fatalf("%s: %v\n", bdb.last, bdb.err)
	}

	reader := bufio.NewReader(os.Stdin)
	collection := scorer.NewCollection(reader, k, b)

	var processed int32
	log.Printf("Run Paradigmatic")
	processed = 0
	collection.Paradigmatic(numWorkers, func(a, b string, score float64) {
		if score < paradigmaticThreshold {
			return
		}

		_, err := pStmt.Exec(a, b, score)
		if err != nil {
			log.Printf("Paradigmatic(%s, %s)=%.3f, err=%v", a, b, score, err)
		}

		n := atomic.AddInt32(&processed, 1)
		if n%1000 == 0 {
			log.Printf("%d processed", n)
		}
	})

	log.Printf("Run Syntagmatic")
	processed = 0
	collection.Syntagmatic(numWorkers, func(a, b string, score float64) {
		if score < syntagmaticThreshold {
			return
		}
		_, err := sStmt.Exec(a, b, score)
		if err != nil {
			log.Printf("Syntagmatic(%s, %s)=%.3f, err=%v", a, b, score, err)
		}
		n := atomic.AddInt32(&processed, 1)
		if n%1000 == 0 {
			log.Printf("%d processed", n)
		}
	})
}
