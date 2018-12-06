collections/%.txt: presets/%.txt
	pipenv run ./compile_collection.py '$<' '$@'
