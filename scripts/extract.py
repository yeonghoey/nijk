from pathlib import Path
import subprocess
import sys


def main(target, local, collection, extractors):
    with open(target) as tf,\
         open(collection, 'w') as cf:
        for line in tf:
            name, lang, _ = line.strip().split()
            path = local / name
            extractor = extractors / lang / 'run'
            assert path.exists()
            assert extractor.exists()
            run(extractor, path, cf)


def run(extractor, path, cf):
    p = subprocess.run([str(extractor), str(path)], stdout=subprocess.PIPE)
    cf.write(p.stdout.decode())


if __name__ == '__main__':
    target = Path(sys.argv[1])
    local = Path(sys.argv[2])
    collection = Path(sys.argv[3])
    extractors = Path(sys.argv[4])
    main(target, local, collection, extractors)
