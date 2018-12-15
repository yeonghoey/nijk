# Extractor `py`

Extractor `py` extracts some contexts from Python 3 source code:
- A class name and its methods names.
- A function name and its argument names.

Things filtered:
- Tests, like modules whose names start with `test_` or end with `_test`.
- Dunder names like `__init__`, and `__str__` in method and function names.
- `self` and `cls` in argument names
- Contexts with a single term.


For example,
```
OrderedDict clear popitem move_to_end keys items values pop setdefault copy fromkeys
popitem last
move_to_end key last
```
