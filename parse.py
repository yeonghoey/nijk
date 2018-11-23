import ast
from itertools import chain
from pathlib import Path
import re


def main(source, target):
    lines = []
    for pyfilename in source.glob('**/*.py'):
        print(f"Parse '{pyfilename}'")
        for context in process(pyfilename):
            lines.append(' '.join(context))

    with open(target, 'w') as f:
        f.write('\n'.join(lines))


def process(pyfilename):
    with open(pyfilename) as f:
        source = f.read()

    root = ast.parse(source, pyfilename)

    for a in ast.walk(root):
        if isinstance(a, ast.FunctionDef):
            yield chain([a.name],
                        [x.arg for x in a.args.args],
                        split(a.name),
                        *[split(x.arg) for x in a.args.args])


def split(s):
    return re.split(r'(?<=[a-zA-Z0-9])_(?=[a-zA-Z0-9])', s)


if __name__ == '__main__':
    main(source=Path('./model-projects'), target=Path('./contexts.txt'))
