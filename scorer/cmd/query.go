package main

import "fmt"

func queryCreateTable(preset, relation string) string {
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

func queryInsert(preset, relation string) string {
	template := `INSERT INTO %s_%s (this, that, score) VALUES (?, ?, ?);`
	return fmt.Sprintf(template, preset, relation)
}
