# Nijk
Nijk is for helping programmers come up with good names.
Nijk analyzes some of the most popular projects to suggest good names for given context.


# Presets
The preset represents for the analysis unit. It is basically a list of projects.
Based on a preset, Nijk extracts identifiers from the project's source code to compose a collection of contexts.

For more details, see [presets](presets).

# Extractors
An extractor is responsible for extracting contexts from a project.
Currently, Nijk supports [Python](extractors/py/run) only.


For more details, see [extractors](extractors).

# Collection
[scripts/compile_collection.py](scripts/compile_collection.py) downloads projects specified in a preset and triggers extractors on those projects.

## Caveat

This is my(yeongho2@illinois.edu) course project for
Text Information Systems of [MCS-DS](https://cs.illinois.edu/academics/graduate/professional-mcs-program/online-master-computer-science-data-science).

