# Path for local copies of source projects specified in targets
EXTRACTORS = extractors/
LOCAL = .local/

collections/%.txt: targets/%.txt
	pipenv run python 'scripts/collect.py' '$<' '$(LOCAL)'
	pipenv run python 'scripts/extract.py' '$<' '$(LOCAL)' '$@' '$(EXTRACTORS)'
