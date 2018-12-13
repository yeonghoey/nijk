package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	presetsPath = "../presets"
	presetExt   = ".txt"
)

type presetItem struct {
	project string
	url     string
}

func loadPresetMeta() map[string][]presetItem {
	presetMeta := make(map[string][]presetItem)
	for _, presetFile := range listPresetFiles() {
		preset := strings.TrimSuffix(presetFile, presetExt)
		items := loadPresetItems(presetFile)
		presetMeta[preset] = items
	}
	return presetMeta
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
		item := presetItem{project: fs[0], url: fs[2]}
		items = append(items, item)
	}
	return items
}
