from collections import Counter, defaultdict
import math
from pprint import pprint

def main(contexts_path):
    with open(contexts_path) as f:
        termctxs, termidfs, ctxvecs = build_index(f)

    termscores = defaultdict(list)
    for term1, ctxs1 in termctxs.items():
        scores = []
        for term2, ctxs2 in termctxs.items():
            if term1 == term2:
                continue
            score = 0
            for c1 in ctxs1:
                for c2 in ctxs2:
                    score += similarity(termidfs, ctxvecs[c1], ctxvecs[c2])
            scores.append((score, term2))
        termscores[term1] = sorted(scores, reverse=True)[:100]
    pprint(termscores)


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
        term: math.log((numctx+1) / len(ctxs))
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


def similarity(termidfs, vec1, vec2):
    return sum(termidfs[t] * vec1[t] * vec2[t] for t in vec1 if t in vec2)


if __name__ == '__main__':
    main(contexts_path='contexts.txt')
