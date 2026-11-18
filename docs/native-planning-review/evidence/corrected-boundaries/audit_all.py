#!/usr/bin/env python3
"""Audit every comparator and candidate with frozen original geometry/arithmetic after timing."""
import argparse
from collections import Counter
import hashlib
import json
import os
from pathlib import Path
import subprocess
import time


def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument('--runs', type=Path, required=True)
    ap.add_argument('--out', type=Path, required=True)
    ap.add_argument('--protocol', type=Path, required=True)
    ap.add_argument('--manifest', type=Path, default=Path(__file__).resolve().parent/'identity.json')
    ap.add_argument('--corpus', type=Path, default=Path('/tmp/rdk-motion-scenes'))
    ap.add_argument('--include-profile', action='store_true')
    ap.add_argument('--block', type=int)
    args = ap.parse_args()
    manifest = json.loads(args.manifest.read_text())
    binary = manifest['binaries']['frozen_audit']
    assert digest(Path(binary['path'])) == binary['sha256'], 'Audit binary changed'
    args.out.mkdir(parents=True, exist_ok=True)
    variants = json.loads(args.protocol.read_text())['variants']
    runs = sorted((p for p in args.runs.iterdir() if p.is_dir() and p.name.rsplit('-', 1)[-1] in variants), key=lambda p: (int(p.name.split('-')[0]), variants.index(p.name.rsplit('-', 1)[-1])))
    if args.block is not None:
        runs = [p for p in runs if p.name.startswith(str(args.block)+'-')]
        variants = json.loads(args.protocol.read_text())['variants']
        assert len(runs) == len(variants), 'Expected every declared variant in each completed block'
    if args.include_profile:
        runs.append(args.runs.parent/'profile')
    assert runs, 'No benchmark processes found'
    catalog = {}
    for run in runs:
        journal = args.out / (run.name + '.jsonl')
        assert not journal.exists(), 'Refuse to overwrite an existing audit'
        env = {k: v for k, v in os.environ.items() if not k.startswith(('RDK_', 'MOTION_'))}
        env.update(RDK_PLANNER_SUITE=str(args.corpus), RDK_OFFLINE_AUDIT_ROOT=str(run),
                   RDK_OFFLINE_AUDIT_OUTPUT=str(journal), RDK_OFFLINE_AUDIT_PROTOCOL=str(args.protocol),
                   MOTION_ROADMAP_CACHE_DIR='')
        cmd = [binary['path'], '-test.run', '^TestFrozenArithmeticOfflineAudit$', '-test.v', '-test.timeout=90m']
        metadata = dict(command=cmd, binary_sha256=binary['sha256'],
                        environment={k: v for k, v in env.items() if k.startswith(('RDK_', 'MOTION_'))},
                        started_utc=time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime()))
        config = args.out / (run.name + '-config.json')
        config.write_text(json.dumps(metadata, indent=2)+'\n')
        print('Auditing', run, flush=True)
        with (args.out/(run.name+'.log')).open('w') as log:
            proc = subprocess.run(cmd, env=env, stdout=log, stderr=subprocess.STDOUT)
        metadata['exit'] = proc.returncode
        config.write_text(json.dumps(metadata, indent=2)+'\n')
        events = [json.loads(line) for line in journal.read_text().splitlines()] if journal.exists() else []
        active, counts, outcomes = {}, Counter(), Counter()
        for event in events:
            key = (event['file'], event['run'], event['warmup'])
            counts[event['event']] += 1
            if event['event'] == 'started':
                assert key not in active
                active[key] = event
            else:
                assert key in active
                del active[key]
                outcomes['accepted' if event['accepted'] else 'rejected_or_no_trajectory'] += 1
                outcomes['warmup' if event['warmup'] else 'measured'] += 1
        expected = sum(len(d.get('results', []))+len(d.get('warmup_results', []))
                       for p in run.rglob('*.json') for d in [json.loads(p.read_text())])
        newly_rejected = sum(e['event'] == 'finished' and e.get('original_audit_accepted', False) and not e.get('accepted', False) and e.get('has_trajectory', False) for e in events)
        completed = not active and counts['finished'] == expected and counts['started'] == expected
        catalog[run.name] = dict(completed=completed, newly_rejected=newly_rejected, expected=expected, events=dict(counts), outcomes=dict(outcomes),
                                 unfinished=list(active.values()), exit=proc.returncode,
                                 journal=str(journal), sha256=digest(journal) if journal.exists() else None,
                                 audit_s=sum(e.get('audit_ns', 0) for e in events)/1e9)
        (args.out/'accounting.json').write_text(json.dumps(catalog, indent=2)+'\n')
        assert completed and (proc.returncode == 0 or (proc.returncode == 1 and newly_rejected > 0)), 'Incomplete audit or unexplained process failure'
        print('Finished', run.name, dict(outcomes), flush=True)


if __name__ == '__main__':
    main()
