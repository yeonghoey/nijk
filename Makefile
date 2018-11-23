.PHONY: init run

run: termscores.json init
	pipenv run python run.py

init:
	pipenv install

termscores.json: contexts.txt process.py
	pipenv run python process.py

contexts.txt: model-projects parse.py
	pipenv run python parse.py
