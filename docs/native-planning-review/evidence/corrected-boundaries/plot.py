#!/usr/bin/env python3
"""Render linear percentile charts from verified raw request archives after timing."""
import argparse
from collections import Counter
import gzip
import hashlib
import itertools
import json
import os
from pathlib import Path

os.environ.setdefault('MPLCONFIGDIR', '/tmp/rdk-native-cleanup-matplotlib')
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
import numpy as np

VARIANTS = ['original', 'intermediate', 'candidate']
LABELS = ['Original native', 'Pre-cleanup native', 'Cleaned native']
COLORS = ['#2563a6', '#7d858d', '#cf6629']
STYLES = ['--', '-.', '-']
WORKS = ['scenes', 'sequence']
TITLES = ['Individual scenes', 'Ordered workload']


def quantiles(values, probabilities):
    values = np.sort(values)
    p = np.asarray(probabilities)
    q = values[np.maximum(0, np.ceil(p*len(values)).astype(int)-1)].copy()
    q[p == .5] = np.median(values)
    return q


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument('--report', type=Path, required=True)
    ap.add_argument('--out', type=Path, required=True)
    ap.add_argument('--variants', nargs='+', required=True)
    ap.add_argument('--baseline', required=True)
    ap.add_argument('--candidate', required=True)
    args = ap.parse_args()
    global VARIANTS, LABELS, COLORS, STYLES
    from analyze import LABELS as names
    palette = dict(zip(names, ['#2563a6','#78838d','#9b72b0','#21a4ad','#3b8b56','#d56b27']))
    VARIANTS = args.variants
    LABELS = [names[v] for v in VARIANTS]
    COLORS = [palette[v] for v in VARIANTS]
    STYLES = ['--' if v == 'original' else ':' if v == 'main' else '-' for v in VARIANTS]
    assert args.baseline in VARIANTS and args.candidate in VARIANTS
    report = json.loads(args.report.read_text())
    blocks = sorted({c['block'] for c in report['configurations'].values()})
    nblocks = len(blocks)
    data = {w: {v: [] for v in VARIANTS} for w in WORKS}
    accepted = {w: {v: [] for v in VARIANTS} for w in WORKS}
    for v in VARIANTS:
        for block in blocks:
            matches = [p for p in report['artifacts'] if p.endswith(f'/{block}-{v}.raw.json.gz')]
            assert len(matches) == 1
            p = Path(matches[0]); payload = p.read_bytes()
            assert hashlib.sha256(payload).hexdigest() == report['artifacts'][str(p)]['sha256']
            raw = json.loads(gzip.decompress(payload))
            group = {w: [] for w in WORKS}
            for name, doc in raw.items():
                if name == 'scene08.json':
                    continue
                w = 'sequence' if name.startswith('sequence/') else 'scenes'
                for row in doc['results']:
                    ms = row['replay_ns']/1e6
                    group[w].append(ms)
                    if row['status'] == 'solved' and row.get('validated', False):
                        accepted[w][v].append(ms)
            for w in WORKS:
                assert len(group[w]) == (275 if w == 'scenes' else 294)
                data[w][v].append(np.asarray(group[w]))
    pooled = {w: {v: np.concatenate(data[w][v]) for v in VARIANTS} for w in WORKS}
    args.out.mkdir(parents=True, exist_ok=True)
    plt.rcParams.update({'font.size': 11, 'axes.spines.top': False, 'axes.spines.right': False})
    def save(fig, name):
        fig.savefig(args.out/(name+'.png'), dpi=170)
        fig.savefig(args.out/(name+'.pdf'))
        plt.close(fig)
    fig, axes = plt.subplots(2, 2, figsize=(12, 8.5), layout='constrained')
    top = np.ceil(max(quantiles(pooled[w][v], [.95])[0] for w in WORKS for v in VARIANTS)/50)*50
    bottom = np.floor(min(quantiles(pooled[w][v], [.95])[0] for w in WORKS for v in VARIANTS)/50)*50
    peak = np.ceil(max(max(pooled[w][v]) for w in WORKS for v in VARIANTS)/100)*100
    for col, w in enumerate(WORKS):
        p = np.unique(np.r_[np.linspace(0, 100, 2001), 50, 95, 99])
        for row in [0, 1]:
            ax = axes[row, col]
            for v, label, color, style in zip(VARIANTS, LABELS, COLORS, STYLES):
                ax.plot(p, quantiles(pooled[w][v], p/100), label=label, color=color, ls=style, lw=2)
            ax.set(xlabel='Percentile of requests', ylabel='Planning latency (ms)')
            ax.grid(alpha=.2)
        axes[0, col].set(xlim=(0, 95), ylim=(0, top), xticks=[0, 25, 50, 75, 90, 95], title=TITLES[col]+' · first 95%')
        axes[0, col].legend(loc='upper left', fontsize=9)
        ax = axes[1, col]
        ax.set(xlim=(95, 100), ylim=(bottom, peak), xticks=[95, 96, 97, 98, 99, 100], title=TITLES[col]+' · slowest 5% (zoom)')
        ax.axvline(99, color='#84909e', ls=':', lw=1)
        lines = ['p99']
        for v, label, color in zip(VARIANTS, LABELS, COLORS):
            q = quantiles(pooled[w][v], [.99])[0]
            ax.plot(99, q, 'o', color=color)
            lines.append(f'{label}: {q:.1f} ms')
        ax.text(.025, .975, '\n'.join(lines), transform=ax.transAxes, va='top', fontsize=9,
                bbox=dict(facecolor='white', edgecolor='#d8dce2', boxstyle='round,pad=.4'))
    fig.suptitle('Planning latency by percentile\nLOWER IS BETTER · linear axes', fontsize=16)
    fig.supxlabel(f'At percentile 95, 95% of requests completed within the plotted time.\nWarm planning replay · all outcomes · {nblocks} process blocks · matched scales within each row', fontsize=9)
    save(fig, 'latency-percentiles')
    probs = np.array([.10, .25, .50, .75, .90, .95, .975, .99, .995])
    fig, axes = plt.subplots(1, 2, figsize=(12, 5.5), layout='constrained')
    effects = {}
    for ax, w, title in zip(axes, WORKS, TITLES):
        delta = quantiles(pooled[w][args.candidate], probs)-quantiles(pooled[w][args.baseline], probs)
        # Exactly n**n paired resamples, coalesced by identical block counts.
        combinations = Counter(tuple(pick.count(i) for i in range(nblocks))
                               for pick in itertools.product(range(nblocks), repeat=nblocks))
        reps = []
        for counts, frequency in combinations.items():
            a, b = [quantiles(np.concatenate([data[w][v][i] for i in range(nblocks)
                                             for _ in range(counts[i])]), probs)
                    for v in [args.baseline, args.candidate]]
            reps.extend([b-a]*frequency)
        intervals = np.array([quantiles(np.asarray(reps)[:, i], [.025, .975]) for i in range(len(probs))])
        ax.vlines(range(len(probs)), intervals[:, 0], intervals[:, 1], color='#596579')
        ax.scatter(range(len(probs)), delta, c=['#16806a' if x < 0 else '#bf4938' for x in delta], zorder=3)
        ax.axhline(0, color='#253140', lw=1.3)
        lo, hi = ax.get_ylim()
        ax.axhspan(lo, 0, color='#16806a', alpha=.045)
        ax.axhspan(0, hi, color='#bf4938', alpha=.045)
        ax.set(ylim=(lo, hi), title=title, ylabel='Candidate − baseline latency (ms)', xlabel='Selected percentile',
               xticks=range(len(probs)), xticklabels=[f'p{100*x:g}' for x in probs])
        ax.tick_params(axis='x', rotation=45); ax.grid(axis='y', alpha=.2)
        ax.text(.02, .97, 'Above zero = slower', va='top', transform=ax.transAxes, color='#a43e30', fontsize=9)
        ax.text(.02, .03, 'Below zero = faster', va='bottom', transform=ax.transAxes, color='#116754', fontsize=9)
        effects[w] = [dict(percentile=float(100*p), delta_ms=float(d), pointwise_ci95_ms=ci.tolist()) for p, d, ci in zip(probs, delta, intervals)]
    fig.suptitle(names[args.candidate]+' vs '+names[args.baseline]+'\nLOWER IS BETTER · zero means no change', fontsize=16)
    fig.supxlabel(f'Pointwise 95% paired-process bootstrap intervals · {nblocks} pairs · selected percentiles evenly spaced\nExploratory; not a simultaneous confidence band or a formal significance claim.', fontsize=9)
    save(fig, 'quantile-effects')
    (args.out/'quantile-effects.json').write_text(json.dumps(effects, indent=2)+'\n')
    fig, axes = plt.subplots(1, 2, figsize=(12, 5.5), layout='constrained')
    for ax, w, title in zip(axes, WORKS, TITLES):
        for v, label, color, style in zip(VARIANTS, LABELS, COLORS, STYLES):
            a = np.sort(accepted[w][v]); n = len(pooled[w][v])
            ax.step(np.r_[0, a, peak], np.r_[0, 100*np.arange(1, len(a)+1)/n, 100*len(a)/n], where='post', color=color, ls=style, lw=2, label=label)
        ax.set(xlim=(0, peak), ylim=(0, 103), title=title, xlabel='Planning deadline (ms)', ylabel='Requests with an accepted plan (%)')
        ax.grid(alpha=.2); ax.legend(loc='lower right', fontsize=9)
    fig.suptitle('Accepted plans by a deadline\nHIGHER IS BETTER · linear axes · failures stay in the denominator', fontsize=15)
    fig.supxlabel('Warm planning replay; audit classifies plans but its duration is outside the timer.', fontsize=9)
    save(fig, 'accepted-by-time')
    print('Rendered percentile, difference, and accepted-by-deadline charts.')


if __name__ == '__main__':
    main()
