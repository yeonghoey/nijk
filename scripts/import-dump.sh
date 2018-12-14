#!/usr/bin/env bash

readonly INSTANCE='nijk-master'
readonly DATABASE='nijk'

readonly PRESET="$1"
readonly LOCAL_PATH="dumps/${PRESET}.sql"
readonly BUCKET_PATH="gs://nijk-scores/${PRESET}.sql"

set -euo pipefail

read -rp "Clear and import data for tables related to '${PRESET}'. Continue? <y/N> " prompt
if [[ ! "${prompt}" =~ [yY](es)* ]]; then
  exit 0
fi

(set -x; go run ./scorer/cmd/schema "${PRESET}")

read -rp "Upload '${LOCAL_PATH}'? <y/N> " prompt
if [[ "${prompt}" =~ [yY](es)* ]]; then
 (set -x; gsutil cp "${LOCAL_PATH}" "${BUCKET_PATH}")
fi

(set -x; gcloud sql import sql "${INSTANCE}" "${BUCKET_PATH}" --database="${DATABASE}")
