# Nijk
Nijk is for helping programmers come up with good names.
Nijk analyzes some of the most popular projects to suggest good names for given context.

## Architecture
![Architecture](static/architecture.png)


## Quick Start
Nijk is written with multiple tools and programming languages.
To generate a dump, which is the overall output of the *Local* phase, you will need:
- Python 3.6 or later
- Go 1.11 or later

### Python
The Python dependencies is managed with [Pipenv].
If you are familiar with [Pipenv], you can simply run `pipenv install`.
You can also configure your own [Virtualenv] and just run `pip install -r requirements.txt`.


### Go
Go packages for Nijk is organized with Go Modules, which was introduced in Go 1.11.

## Presets
A preset is basically a list of projects to be used as ideal naming examples.
It is also the basic analysis unit of Nijk.
Based on a preset, Nijk extracts identifiers from source code in projects listed, to compose a collection of contexts.

For more details, see [presets](presets).

## Extractors
An extractor is responsible for extracting contexts from a project.
Currently, Nijk supports [Python](extractors/py/run) only.


For more details, see [extractors](extractors).

## Collection
[scripts/compile_collection.py](scripts/compile_collection.py) first downloads projects specified in a preset.
And then, it executes extractors on those projects. Finally, it concatenates the outputs of the executions. 

For more details, see [collections](collections).

## Scorer
Scorer is the core of Nijk. It reads a collection and run Paradigmatic and Syntagmatic Relation Discovery algorithm based on normalized-[BM25](https://en.wikipedia.org/wiki/Okapi_BM25).
Scorer is implemented in Go. It also has a command-line interface which generates SQL dump queries to be imported to a MySQL server.

For more details, see [package scorer](https://godoc.org/github.com/yeonghoey/nijk/scorer).

## Caveat

This is my(yeongho2@illinois.edu) course project for
Text Information Systems of [MCS-DS](https://cs.illinois.edu/academics/graduate/professional-mcs-program/online-master-computer-science-data-science).


[Pipenv]: https://pipenv.readthedocs.io/en/latest/.
[Virtualenv]: https://virtualenv.pypa.io/en/latest/
