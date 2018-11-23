.PHONY: init run

init:
	pipenv install

run: termscores.json
	pipenv run python run.py

termscores.json: contexts.txt process.py
	pipenv run python process.py

contexts.txt: model-projects parse.py
	pipenv run python parse.py
