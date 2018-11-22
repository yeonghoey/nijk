import ast
import os
from pathlib import Path


def main(source, target):
    lines = []
    for pyfilename in source.glob('**/*.py'):
        names = process(pyfilename)
        lines.append(' '.join(names))

    with open(target, 'w') as f:
        f.write('\n'.join(lines))


def process(pyfilename):
    with open(pyfilename) as f:
        source = f.read()

    root = ast.parse(source, pyfilename)

    names = []
    for a in ast.walk(root):
        if isinstance(a, ast.Name):
            names.append(a.id.lower())
    return names


if __name__ == '__main__':
    main(source=Path('./model-projects'),
         target=Path('./contexts.txt'))
