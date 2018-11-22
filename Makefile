.PHONY: test

test: contexts.txt
	pipenv run python test.py

contexts.txt: model-projects parse.py
	pipenv run python parse.py
