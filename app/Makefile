.PHONY: dev deploy

dev:
	dev_appserver.py .

deploy:
	go mod tidy
	gcloud app deploy
