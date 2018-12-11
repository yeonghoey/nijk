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

resource "google_storage_bucket" "nijk-scores" {
  name     = "nijk-scores"
  location = "asia"
}

resource "google_storage_bucket_acl" "nijk-scores-acl" {
  bucket = "${google_storage_bucket.nijk-scores.name}"

  role_entity = [
    "READER:user-${google_sql_database_instance.nijk.service_account_email_address}",
  ]
}

resource "google_storage_default_object_acl" "nijk-scores-acl" {
  bucket = "${google_storage_bucket.nijk-scores.name}"

  role_entity = [
    "READER:user-${google_sql_database_instance.nijk.service_account_email_address}",
  ]
}
