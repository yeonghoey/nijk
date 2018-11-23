import json

from prompt_toolkit import PromptSession


session = PromptSession()


def main(termscores_path):
    with open(termscores_path) as f:
        termscores = json.load(f)

    while True:
        term = session.prompt('nijk > ')
        scores = termscores.get(term, [])
        simterms = ['%s(%.1f)' % (simterm, score) for score, simterm in scores]
        for part in chunks(simterms, 5):
            print(' '.join(part))
        print()


def chunks(l, n):
    for i in range(0, len(l), n):
        yield l[i:i+n]


if __name__ == '__main__':
    main(termscores_path='termscores.json')
