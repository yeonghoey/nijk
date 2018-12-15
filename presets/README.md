# Presets
Each preset should only contain lines in following format:

```
<Name> <Extractor> <URL>
```

- `<Name>` is generally composed of `<project>-v<version>`, like `cpython-v3.7.1`.
- `<Extractor>` should be only of the subdirectory name in `../extractors`.
- `<URL>` is a URL for downloading the source code of the project.

Reading these lines, [../scripts/compile_collection.py](../scripts/compile_collection.py) compiles a collection.
