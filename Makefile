dumps/%.sql: collections/%.txt
	go run ./scorer/cmd < '$<' > '$@'

collections/%.txt: presets/%.txt
	pipenv run ./compile_collection.py '$<' '$@'
