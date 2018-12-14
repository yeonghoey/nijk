.PHONY: dev deploy

dev:
	dev_appserver.py .

deploy:
	go mod tidy
	gcloud app deploy

dumps/%.sql: collections/%.txt
	go run ./scorer/cmd '$*' < '$<' > '$@'

collections/%.txt: presets/%.txt
	pipenv run ./compile_collection.py '$<' '$@'
