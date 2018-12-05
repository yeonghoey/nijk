from pathlib import Path
import sys
from tempfile import TemporaryFile
from zipfile import ZipFile

import requests


def main(target, local):
    with open(target) as f:
        for line in f:
            name, _, url = line.strip().split()
            path = local / name
            if not path.exists():
                with TemporaryFile() as f:
                    response = requests.get(url)
                    f.write(response.content)
                    f.seek(0)
                    ZipFile(f).extractall(path)


if __name__ == '__main__':
    target = Path(sys.argv[1])
    local = Path(sys.argv[2])
    main(target, local)
