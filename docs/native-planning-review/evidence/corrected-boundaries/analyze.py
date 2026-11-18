#!/usr/bin/env python3
"""Analyze completed serial experiments, preserving independent process blocks."""
import argparse
from collections import Counter, defaultdict
import gzip
import hashlib
import itertools
import json
import math
from pathlib import Path
import statistics

HERE = Path(__file__).resolve().parent
VARIANTS = ['original', 'main', 'pr1', 'pr2', 'pr3']
LABELS = dict(original='Original native', main='Current main', pr1='PR 1', pr2='PRs 1–2', pr3='PRs 1–3')
PAIRS = [('original', v) for v in VARIANTS[1:]] + [('main','pr1'), ('pr1','pr2'), ('pr2','pr3')]


def quantile(values, fraction):
    values = sorted(values)
    return values[max(0, math.ceil(fraction * len(values)) - 1)] if values else None


def dist(values):
    values = list(values)
    return dict(median=statistics.median(values), p95=quantile(values, .95),
                p99=quantile(values, .99), max=max(values)) if values else None


def length(row):
    path = row.get('trajectory', [])
    return sum(math.sqrt(sum((x-y)**2 for name in a if a[name] for x, y in zip(a[name], b[name])))
               for a, b in zip(path, path[1:]))


def accepted(row):
    return row['status'] == 'solved' and row.get('validated', False)


def summarize(rows):
    good = [r for r in rows if accepted(r)]
    diagnostics = [r['planning']['native_diagnostics'] for r in rows if r.get('planning', {}).get('native_diagnostics')]
    phases, work = Counter(), Counter()
    solvers = Counter()
    for row in rows:
        meta = row.get('planning', {})
        solvers.update({k: meta[k] for k in ['GoalsDirectSolved', 'GoalsNudgeSolved', 'GoalsRoadmapSolved', 'GoalsCBIRRTSolved'] if k in meta})
    for d in diagnostics:
        phases.update({k[:-3]: v / 1e9 for k, v in d.items() if k.endswith('_ns')})
        work.update({k: v for k, v in d.get('collision', {}).items() if k != 'exact_state_peak_key_bytes'})
    return dict(n=len(rows), accepted=len(good), accepted_rate=len(good)/len(rows) if rows else None,
                status=dict(Counter(r['status'] for r in rows)), total_s=sum(r['replay_ns'] for r in rows)/1e9,
                completion_ms=dist(r['replay_ns']/1e6 for r in rows),
                accepted_ms=dist(r['replay_ns']/1e6 for r in good),
                audit_s=sum(r.get('audit_ns', 0) for r in rows)/1e9,
                thresholds={str(s): sum(r['replay_ns'] >= s*1e9 for r in rows) for s in [1, 5, 30]},
                phase_s=dict(phases) if diagnostics else None,
                work=dict(work) if diagnostics else None,
                exact_state_peak_key_bytes=dist(d['collision']['exact_state_peak_key_bytes'] for d in diagnostics
                    if 'exact_state_peak_key_bytes' in d.get('collision', {})),
                search_lifetime={
                    'observed_requests': len(diagnostics),
                    'unavailable_requests': len(rows)-len(diagnostics),
                    'started': sum(d['search_workers_started'] for d in diagnostics),
                    'finished': sum(d['search_workers_finished'] for d in diagnostics),
                    'unfinished_requests': sum(d['search_workers_started'] != d['search_workers_finished'] for d in diagnostics),
                } if diagnostics and all('search_workers_started' in d for d in diagnostics) else None,
                search_selection={
                    key: sum(d[key] for d in diagnostics) for key in [
                        'search_solutions_completed', 'search_alternative_selections',
                        'search_first_path_cost', 'search_selected_path_cost']
                } if diagnostics and all('search_solutions_completed' in d for d in diagnostics) else None,
                goal_solver_counts=dict(solvers) or None,
                effective_calibration=dict(Counter(str(d['startup_speed_multiplier']) for d in diagnostics)) or None,
                measured_calibration=dict(Counter(str(d['measured_startup_speed_multiplier']) for d in diagnostics)) or None,
                scene_cache_snapshot_max={k: max(d.get('scene_cache', {}).get(k, 0) for d in diagnostics)
                                          for k in ['bytes', 'budget_bytes', 'entries', 'evictions']} if diagnostics else None,
                errors=dict(Counter(r['error'] for r in rows if r.get('error'))))


def bootstrap(base, candidate):
    """Exhaust all n**n paired process resamples; never resample requests IID."""
    n = len(base)
    if n < 2:
        return {'note': 'One independent process: no uncertainty interval.'}
    ratios, differences = [], []
    for selected in itertools.product(range(n), repeat=n):
        a = sum(base[i]['total_s'] for i in selected)
        b = sum(candidate[i]['total_s'] for i in selected)
        count = sum(base[i]['n'] for i in selected)
        ratios.append(b/a)
        differences.append(1000*(b-a)/count)
    a, b = sum(s['total_s'] for s in base), sum(s['total_s'] for s in candidate)
    return dict(independent_pairs=n, resamples=n**n,
                ratio=b/a, ratio_ci95=[quantile(ratios, .025), quantile(ratios, .975)],
                delta_mean_ms=1000*(b-a)/sum(s['n'] for s in base),
                delta_mean_ms_ci95=[quantile(differences, .025), quantile(differences, .975)],
                note='Exploratory paired process-block percentile bootstrap; few independent processes cannot establish precise tails or generalization.')


def table(entries, accepted_only=False):
    out = ['| Variant / process | N | Accepted | Total s | Median ms | p95 ms | p99 ms | Max ms | ≥1/5/30 s |',
           '|---|---:|---:|---:|---:|---:|---:|---:|---:|---|']
    for name, s in entries:
        d = s['accepted_ms' if accepted_only else 'completion_ms']
        if d is None:
            continue
        out.append(f"| {name} | {s['n']} | {s['accepted']} ({100*s['accepted_rate']:.2f}%) | {s['total_s']:.3f} | " +
                   ' | '.join(f'{d[k]:.2f}' for k in ['median', 'p95', 'p99', 'max']) +
                   ' | ' + '/'.join(str(s['thresholds'][str(x)]) for x in [1, 5, 30]) + ' |')
    return out


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument('--root', type=Path, default=Path('/tmp/rdk-native-batch'))
    ap.add_argument('--stage', choices=['pilot', 'confirmation'], required=True)
    ap.add_argument('--report-dir', type=Path, default=Path('/tmp/rdk-native-pr-cleanup/analysis'))
    ap.add_argument('--candidate-label', default='Cleaned native')
    ap.add_argument('--intermediate-label', default='Pre-cleanup native')
    ap.add_argument('--protocol', type=Path, default=HERE / 'protocol.json')
    args = ap.parse_args()
    protocol = json.loads(args.protocol.read_text())
    declared = {f'{b}-{v}' for b, order in enumerate(protocol['variant_order'], 1) for v in order}
    args.report_dir.mkdir(parents=True, exist_ok=True)
    LABELS['candidate'] = args.candidate_label
    LABELS['intermediate'] = args.intermediate_label
    archive = args.root / (args.stage + '-archives')
    archive.mkdir(exist_ok=True)
    rows_by_variant = defaultdict(list)
    process = {}
    journals = {}
    artifacts = {}
    configs = {}
    for config_path in sorted((args.root / args.stage).glob('*-config.json')):
        config = json.loads(config_path.read_text())
        configs[config_path.stem] = config
        variant, block = config['variant'], config['block']
        name = f'{block}-{variant}'
        run = args.root / args.stage / name
        rows, raw, events, active = [], {}, [], {}
        counts = Counter()
        warmups = []
        for path in sorted(run.rglob('*.jsonl')):
            relative = str(path.relative_to(run))
            for line in path.read_text().splitlines():
                event = json.loads(line)
                event['file'] = relative
                events.append(event)
                data = event['data']
                row = data.get('result', data)
                key = (relative, row['run'], data['warmup'])
                counts[event['event']] += 1
                if event['event'] == 'started':
                    assert key not in active, key
                    active[key] = event
                elif event['event'] == 'plan_returned':
                    assert key in active and active[key]['event'] == 'started', key
                    active[key] = event
                elif event['event'] == 'finished':
                    assert key in active and active[key]['event'] == 'plan_returned', key
                    del active[key]
                    if data['warmup']:
                        warmups.append(row)
        for path in sorted(run.rglob('*.json')):
            data = json.loads(path.read_text())
            data.pop('snapshot_request', None)
            relative = str(path.relative_to(run))
            raw[relative] = data
            rows.extend(dict(r, source=relative, input=path.name, block=block, variant=variant) for r in data['results'])
        journals[name] = dict(events=dict(counts), unfinished=list(active.values()), measured=len(rows),
                              warmup=summarize(warmups), process_exit=config.get('exit'))
        for suffix, data in [('raw', raw), ('journal', events)]:
            path = archive / f'{name}.{suffix}.json.gz'
            path.write_bytes(gzip.compress(json.dumps(data, separators=(',', ':')).encode(), mtime=0))
            artifacts[str(path)] = dict(sha256=hashlib.sha256(path.read_bytes()).hexdigest(), bytes=path.stat().st_size)
        # Persist accounting before any assertion, so incomplete runs remain visible.
        (args.report_dir / (args.stage + '-accounting.json')).write_text(json.dumps(journals, indent=2)+'\n')
        expected = 594
        assert not active and config.get('exit') == 0, f'Incomplete process {name}'
        assert len(rows) == expected and all(counts[k] == expected+110 for k in ['started', 'plan_returned', 'finished']), name
        rows_by_variant[variant].extend(rows)
        process[name] = {workload: summarize([r for r in rows if belongs(r, workload)]) for workload in ['scenes', 'sequence']}
    assert set(process) == declared, 'Incomplete declared experiment; accounting retained, no pooled inference'
    block_count = len(protocol['variant_order'])
    result = dict(stage=args.stage, processes=process, journals=journals, configurations=configs,
                  artifacts=artifacts, primary_baseline='original', workloads={})
    interference_path = args.root / args.stage / 'INTERFERENCE.json'
    result['environment_interference'] = json.loads(interference_path.read_text()) if interference_path.exists() else None
    lines = ['# Corrected PR boundaries: three-block follow-up', '', f'Stage: **{args.stage}**. Pilot results use one independent process per variant; this declared experiment uses {block_count}.', '',
             'Original native is the primary baseline. Unchanged current main is an additional comparator; each PR is cumulative and compared with its predecessor.', '',
             'Planning replay measures the actual PlanMotion call, including IK, search, required validation, finishing and synchronous cleanup. Original native does not join all losing search workers; its complete shutdown time is unavailable. PR 1 and later join search workers. The shipping branches omit benchmark-only diagnostics; absent phase and cache counters are unavailable. Snapshot preparation, RPC, hardware execution, warmup and the independent audit are excluded. These are CPU results on a shared workstation.', '',
             'Scene aggregates weight the 11 distinct files equally. Scene 08 duplicates scene 02 and remains in the ledger with zero headline weight. The complete ordered workload is separate. Each process runs scenes, then one complete sequence warmup and three measured passes; its internal cache history is preserved.', '',
             'All completion latencies include failures; accepted-plan latency is conditional on independent audit success. Median averages the two middle observations when needed; p95/p99 use nearest rank on pooled raw observations. Timeout latency includes cleanup and is not a successful time-to-solution.', '']
    if result['environment_interference']:
        lines[2:2] = ['**Environmental interference was recorded. Retain this entire cohort as contaminated; do not describe it as a clean confirmation.**', '']
    for workload in ['scenes', 'sequence']:
        pooled, per_input = {}, {}
        for variant in VARIANTS:
            rows = [r for r in rows_by_variant[variant] if belongs(r, workload)]
            pooled[variant] = summarize(rows)
            inputs = defaultdict(list)
            for row in rows:
                inputs[row['input']].append(row)
            per_input[variant] = {name: dict(summary=summarize(group),
                                            accepted_path_length=dist(length(r) for r in group if accepted(r)),
                                            accepted_waypoints=dist(len(r.get('trajectory', [])) for r in group if accepted(r)))
                                  for name, group in sorted(inputs.items())}
        blocks = sorted(set(r['block'] for r in rows_by_variant['pr3']))
        comparisons = {f'{base}_to_{candidate}': bootstrap([process[f'{b}-{base}'][workload] for b in blocks],
                                       [process[f'{b}-{candidate}'][workload] for b in blocks])
                       for base, candidate in PAIRS}
        quality = {}
        for base, candidate in PAIRS:
            quality[f'{base}_to_{candidate}'] = []
            for name, cand in per_input[candidate].items():
                old = per_input[base][name]
                if old['accepted_path_length'] and cand['accepted_path_length']:
                    quality[f'{base}_to_{candidate}'].append(dict(input=name, median_ratio=cand['accepted_path_length']['median']/max(1e-12, old['accepted_path_length']['median']),
                                              p95_ratio=cand['accepted_path_length']['p95']/max(1e-12, old['accepted_path_length']['p95'])))
        outliers = {v: sorted([r for r in rows_by_variant[v] if belongs(r, workload)], key=lambda r: r['replay_ns'], reverse=True)[:10] for v in VARIANTS}
        # Trajectories are in the raw archive; retain identities and diagnostics here.
        outliers = {v: [{k: val for k, val in r.items() if k != 'trajectory'} for r in group] for v, group in outliers.items()}
        result['workloads'][workload] = dict(pooled=pooled, comparisons=comparisons, per_input=per_input,
                                            quality_comparisons=quality, slowest=outliers)
        lines += [f'## {workload.title()}: all-outcome completion latency', ''] + table([(LABELS[v], pooled[v]) for v in VARIANTS]) + ['']
        for base, candidate in PAIRS:
            c = comparisons[f'{base}_to_{candidate}']
            if 'ratio' in c:
                lines += [f"{LABELS[candidate]} against {LABELS[base]}: total-time ratio {c['ratio']:.3f}, exploratory 95% interval [{c['ratio_ci95'][0]:.3f}, {c['ratio_ci95'][1]:.3f}]; mean difference {c['delta_mean_ms']:.2f} ms/request [{c['delta_mean_ms_ci95'][0]:.2f}, {c['delta_mean_ms_ci95'][1]:.2f}].", '']
        lines += ['Accepted-only percentiles (N, accepted, total and threshold columns still describe all attempts):', '']
        lines += table([(LABELS[v], pooled[v]) for v in VARIANTS], accepted_only=True) + ['']
        lines += ['Each independent process:', ''] + table([(f'{LABELS[v]} / {b}', process[f'{b}-{v}'][workload]) for v in VARIANTS for b in blocks]) + ['']
        lines += ['Per-input distributions and path lengths:', '',
                  '| Input | Variant | N / accepted | Total s | Median / p95 / p99 / max ms | Path length median / p95 / max | Waypoints median / p95 |',
                  '|---|---|---:|---:|---|---|---|']
        for name in per_input['original']:
            for variant in VARIANTS:
                d = per_input[variant][name]
                s = d['summary']
                latency = ' / '.join(f'{s["completion_ms"][k]:.2f}' for k in ['median', 'p95', 'p99', 'max'])
                paths = ' / '.join(f'{d["accepted_path_length"][k]:.4f}' for k in ['median', 'p95', 'max']) if d['accepted_path_length'] else 'unavailable'
                waypoints = ' / '.join(f'{d["accepted_waypoints"][k]:g}' for k in ['median', 'p95']) if d['accepted_waypoints'] else 'unavailable'
                lines.append(f"| {name} | {LABELS[variant]} | {s['n']} / {s['accepted']} | {s['total_s']:.3f} | {latency} | {paths} | {waypoints} |")
        lines += ['']
    result['duplicate_scene08'] = {v: summarize([r for r in rows_by_variant[v] if r['input'] == 'scene08.json']) for v in VARIANTS}
    lines += ['## Duplicate scene 08 (retained, zero headline weight)', '']
    lines += table([(LABELS[v], result['duplicate_scene08'][v]) for v in VARIANTS]) + ['']
    lines += ['## Limits and reproduction', '',
              f'The experiment has {block_count} independent cache initializations per variant, not thousands of independent requests. Intervals exhaust all {block_count}^{block_count} paired process-block bootstrap resamples and are exploratory. There is no precise p99 or formal significance claim. Earlier completed comparisons are not retroactively significant.', '',
              'All variants use ordinary automatic production calibration and default worker counts. Modern-head settings are recorded outside the timer in process logs. Missing original diagnostics are unavailable, never zero. Phase CPU/duration totals may overlap and must not be added as disjoint wall time.', '',
              'Path lengths are compared within each input, conditional on audit acceptance; they are joint-space lengths, not execution times. Scene-cache fields are maxima of global retained-memory snapshots, not sums or request-local allocations. Warmup planning and audit costs are recorded separately in process accounting; startup calibration, request decoding and snapshot preparation remain outside measured replay.', '',
              'Power, thermal and load observations (including unavailable readings), exact commands, outcomes, audit costs, per-process calibration, overlapping phase totals, counters and outlier identities are in the JSON report. Dense audits provide evidence at their resolution, not proof of continuous collision freedom. No hardware moved.', '',
              'Raw trajectories and journals are local under the experiment output directory, checksum-indexed in the JSON report. Independent frozen offline audit results are stored separately and must be reconciled before claims of acceptance. A commit of these reports does not back up those payloads. Reproduce with `run_observed.py --candidate <frozen-binary> --stage confirmation --out <fresh-directory>`, then run `analyze.py` after all timing processes exit.', '']
    if all(c.get('policy') == 'production' for c in configs.values()):
        lines = [line.replace('Intermediate and candidate use effective calibration 1/3.',
                             'All variants use ordinary automatic startup calibration; measured/effective settings are in the JSON.')
                 for line in lines]
    (args.report_dir / (args.stage + '-results.json')).write_text(json.dumps(result, indent=2)+'\n')
    (args.report_dir / (args.stage + '-results.md')).write_text('\n'.join(lines)+'\n')
    print('Wrote all-outcome, accepted-only, per-input and process-block results.')


def belongs(row, workload):
    return (row['source'].startswith('scene') and row['input'] != 'scene08.json') if workload == 'scenes' else row['source'].startswith('sequence/')


if __name__ == '__main__':
    main()
