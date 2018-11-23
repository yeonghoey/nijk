from collections import Counter, defaultdict
import concurrent
import concurrent.futures
from functools import lru_cache
import json
import math


JOB_SIZE = 1


def main(contexts_path, target_path):
    print('- build_index')
    with open(contexts_path) as f:
        termctxs, termidfs, ctxvecs = build_index(f)
    print('o build_index')

    print('- build_termscores')
    termscores = build_termscores(termctxs, termidfs, ctxvecs)
    print('o build_termscores')

    with open(target_path, 'w') as f:
        json.dump(termscores, f, ensure_ascii=False, indent=2)

    with open(target_path) as f:
        print(f.read())


def build_index(contexts):
    termctxs = defaultdict(list)
    ctxinfos = {}
    totallen = 0
    for ctxid, context in enumerate(contexts):
        terms = context.split()
        ctxlen = len(terms)
        ctxcnts = Counter(terms)
        ctxinfos[ctxid] = (ctxlen, ctxcnts)
        totallen += ctxlen
        # Update document frequencies
        for term in ctxcnts:
            termctxs[term].append(ctxid)

    numctx = len(ctxinfos)
    termidfs = {
        term: math.log(1 + (numctx - len(ctxs) + 0.5) / (len(ctxs) + 0.5))
        for term, ctxs in termctxs.items()
    }

    avglen = totallen / numctx
    ctxvecs = {
        ctxid: bm25vec(ctxlen, ctxcnts, avglen)
        for ctxid, (ctxlen, ctxcnts) in ctxinfos.items()
    }

    return termctxs, termidfs, ctxvecs


def bm25vec(ctxlen, ctxcnts, avglen, k=1.2, b=.75):
    vector = {}
    common = k*(1 - b + b*(ctxlen/avglen))
    total = 0
    for term, cnt in ctxcnts.items():
        bm25 = (k+1) * cnt / (cnt + common)
        vector[term] = bm25
        total += bm25
    normalized = {term: bm25/total for term, bm25 in vector.items()}
    return normalized


def build_termscores(termctxs, termidfs, ctxvecs):
    termscores = {}
    futures = []
    with concurrent.futures.ProcessPoolExecutor() as executor:
        jobs = list(chunks(list(termctxs.items()), JOB_SIZE))
        for job in jobs:
            f = executor.submit(process_term, termctxs, termidfs, ctxvecs, job)
            futures.append(f)
        print('jobs submitted: %d (%d each)' % (len(jobs), JOB_SIZE))
        for i, f in enumerate(futures, 1):
            part = f.result()
            termscores.update(part)
            print('%d/%d' % (i, len(futures)))
    return termscores


def process_term(termctxs, termidfs, ctxvecs, job):
    termscores = {}

    @lru_cache(maxsize=1000)
    def similarity(c1, c2):
        vec1 = ctxvecs[c1]
        vec2 = ctxvecs[c2]
        return sum(termidfs[t] * vec1[t] * vec2[t] for t in vec1 if t in vec2)

    for term1, ctxs1 in job:
        scores = []
        for term2, ctxs2 in termctxs.items():
            if term1 == term2:
                continue
            score = 0
            for ctxid1 in ctxs1:
                for ctxid2 in ctxs2:
                    if ctxid1 == ctxid2:
                        continue
                    c1, c2 = ((ctxid1, ctxid2) if ctxid1 < ctxid2 else
                              (ctxid2, ctxid1))
                    score += similarity(c1, c2)
            if score > 0:
                scores.append((score, term2))
        termscores[term1] = sorted(scores, reverse=True)[:10]
    return termscores


def chunks(l, n):
    for i in range(0, len(l), n):
        yield l[i:i+n]


if __name__ == '__main__':
    main(contexts_path='contexts.txt', target_path='termscores.json')
