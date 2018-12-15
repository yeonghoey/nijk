# Extractors

Each subdirectory of this represents an extractor.
Each extractor should contain an executable `run` script.

`run` script should do:
- Accept a source code project path as the first argument.
- Traversing the project, extract proper contexts by parsing some source code.
- Print contexts, which is composed of unordered terms separated with spaces to `stdout`.

By specifying the extractor name in a [preset][../presets], you can extract some contexts from projects.

Currently, there is only one extractor, [py](py).
