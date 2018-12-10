provider "google" {
  project = "nijk-225007"
  region  = "asia-northeast1"
  zone    = "asia-northeast1-c"
}

resource "google_sql_database_instance" "nijk" {
  name = "nijk-master"
  database_version = "MYSQL_5_7"
  region = "asia-northeast1"

  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_database" "nijk" {
  instance  = "${google_sql_database_instance.nijk.name}"
  name      = "nijk"
  charset   = "utf8"
  collation = "utf8_bin"
}

resource "google_sql_user" "nijk" {
  instance = "${google_sql_database_instance.nijk.name}"
  name     = "nijk"
}
