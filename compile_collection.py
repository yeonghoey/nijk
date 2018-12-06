#!/usr/bin/env python

from pathlib import Path
import subprocess
import sys
from tempfile import TemporaryFile
from zipfile import ZipFile

import click
import requests


def pathlibdir(ctx, param, value):
    p = Path(value)
    p.mkdir(parents=True, exist_ok=True)
    return p


@click.command()
@click.option('--presets',
              default='./presets',
              type=click.Path(file_okay=False),
              callback=pathlibdir)
@click.option('--extractors',
              default='./extractors',
              type=click.Path(file_okay=False),
              callback=pathlibdir)
@click.option('--collections',
              default='./collections',
              type=click.Path(file_okay=False),
              callback=pathlibdir)
@click.option('--local',
              default='./.local',
              type=click.Path(file_okay=False),
              callback=pathlibdir)
@click.argument('name')
def main(presets, extractors, collections, local, name):
    filename = f'{name}.txt'
    preset = presets / filename
    collection = collections / filename
    with open(preset) as pf,\
         open(collection, 'w') as cf:
        for line in pf:
            name, lang, url = line.strip().split()

            # Download if the source project does not exist in local
            source_path = local / name
            if not source_path.exists():
                with TemporaryFile() as f:
                    response = requests.get(url)
                    f.write(response.content)
                    f.seek(0)
                    ZipFile(f).extractall(source_path)

            # Extract contexts from the source project
            extractor = extractors / lang / 'run'
            assert extractor.exists()
            cmd = [str(extractor), str(source_path)]
            p = subprocess.run(cmd, stdout=subprocess.PIPE)
            cf.write(p.stdout.decode())


if __name__ == '__main__':
    main()
