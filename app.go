package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// TODO: Parameterize these constants
const (
	presetsPath = "presets"
	presetExt   = ".txt"
	instance    = "nijk-225007:asia-northeast1:nijk-master"
	user        = "nijk"
	topN        = 50
)

var presets map[string][]presetItem
var db *sql.DB

type presetItem struct {
	Project string
	URL     string
}

type formParam struct {
	Preset string
	This   string
}

type sectionParam struct {
	Title string
	Thats []string
}

func init() {
	presets = loadPresets()

	isDev := os.Getenv("DEV") != ""
	if isDev {
		db = openDBDev()
	} else {
		gin.SetMode(gin.ReleaseMode)
		db = openDBProd()
	}
}

func loadPresets() map[string][]presetItem {
	presets := make(map[string][]presetItem)
	for _, presetFile := range listPresetFiles() {
		preset := strings.TrimSuffix(presetFile, presetExt)
		items := loadPresetItems(presetFile)
		presets[preset] = items
	}
	return presets
}

func listPresetFiles() []string {
	files, err := ioutil.ReadDir(presetsPath)
	if err != nil {
		log.Fatal(err)
	}

	presetFiles := make([]string, 0, len(files))
	for _, file := range files {
		filename := file.Name()
		ext := filepath.Ext(filename)
		if ext == presetExt {
			presetFiles = append(presetFiles, filename)
		}
	}
	return presetFiles
}

func loadPresetItems(presetFile string) []presetItem {
	p := path.Join(presetsPath, presetFile)
	f, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	reader := bufio.NewReader(f)
	items := []presetItem{}
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		fs := strings.Fields(line)
		item := presetItem{Project: fs[0], URL: fs[2]}
		items = append(items, item)
	}
	return items
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

func main() {
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"title": strings.Title,
		"formParam": func(p, t string) formParam {
			return formParam{p, t}
		},
		"sectionParam": func(t string, ts []string) sectionParam {
			return sectionParam{t, ts}
		},
	})

	router.LoadHTMLGlob("_templates/*")
	router.GET("/", handleIndex)
	router.GET("/:preset", handlePreset)
	router.GET("/:preset/:this", handleThis)
	router.Run()
}

func handleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"presets": presets,
	})
}

func handlePreset(c *gin.Context) {
	preset := c.Param("preset")
	items, ok := presets[preset]
	if !ok {
		return
	}

	if !redirectOnQuery(c) {
		c.HTML(http.StatusOK, "preset.html", gin.H{
			"preset": preset,
			"items":  items,
		})
	}
}

func handleThis(c *gin.Context) {
	preset := c.Param("preset")
	this := c.Param("this")

	if _, exists := presets[preset]; !exists {
		return
	}

	if !redirectOnQuery(c) {
		paradigmatic := make(chan []string)
		go func() { paradigmatic <- queryParadigmatic(preset, this) }()
		syntagmatic := make(chan []string)
		go func() { syntagmatic <- querySyntagmatic(preset, this) }()
		containing := make(chan []string)
		go func() { containing <- queryContaining(preset, this) }()

		c.HTML(http.StatusOK, "this.html", gin.H{
			"preset":       preset,
			"this":         this,
			"paradigmatic": <-paradigmatic,
			"syntagmatic":  <-syntagmatic,
			"containing":   <-containing,
		})
	}
}

func redirectOnQuery(c *gin.Context) bool {
	q := c.Query("q")
	if q == "" {
		return false
	}

	preset := c.Param("preset")
	location := fmt.Sprintf("/%s/%s", preset, q)
	c.Redirect(http.StatusMovedPermanently, location)
	return true
}

func queryParadigmatic(preset, this string) []string {
	table := fmt.Sprintf("`%s_paradigmatic`", preset)
	q := fmt.Sprintf(`
SELECT that
FROM %s
WHERE this=?
ORDER BY score DESC
LIMIT %d;`, table, topN)
	return query(q, this)
}

func querySyntagmatic(preset, this string) []string {
	table := fmt.Sprintf("`%s_syntagmatic`", preset)
	q := fmt.Sprintf(`
SELECT that
FROM %s
WHERE this=?
ORDER BY score DESC
LIMIT %d;`, table, topN)
	return query(q, this)
}

func queryContaining(preset, this string) []string {
	table := fmt.Sprintf("`%s_paradigmatic`", preset)
	q := fmt.Sprintf(`
SELECT DISTINCT this
FROM %s
WHERE this<>? AND this LIKE ?
LIMIT %d;`, table, topN)
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
