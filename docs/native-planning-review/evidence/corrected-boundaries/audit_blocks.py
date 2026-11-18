#!/usr/bin/env python3
"""Audit independent blocks concurrently, strictly after all headline timing."""
import argparse
from concurrent.futures import ThreadPoolExecutor
import datetime
import json
from pathlib import Path
import subprocess
import sys

HERE = Path(__file__).resolve().parent


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument('--runs', type=Path, required=True)
    ap.add_argument('--out', type=Path, required=True)
    ap.add_argument('--protocol', type=Path, required=True)
    args = ap.parse_args()
    protocol = json.loads(args.protocol.read_text())
    variants = protocol.get('variants', ['original', 'intermediate', 'candidate'])
    blocks = range(1, len(protocol['variant_order'])+1)
    for block in blocks:
        for variant in variants:
            config = json.loads((args.runs/f'{block}-{variant}-config.json').read_text())
            assert config.get('exit') == 0 and config.get('finished_utc'), 'All timing must finish first'
    args.out.mkdir(parents=True, exist_ok=False)
    commands = {str(block): [sys.executable, str(HERE/'audit_all.py'), '--runs', str(args.runs),
                             '--out', str(args.out/f'block{block}'), '--protocol', str(args.protocol),
                             '--block', str(block)] for block in blocks}
    metadata = dict(started_utc=datetime.datetime.now(datetime.timezone.utc).isoformat(),
                    parallel_blocks=4, commands=commands,
                    note='Independent frozen audit processes only; headline timing already complete. Audit durations include concurrent audit load and are not planner latency comparisons.')
    config_path = args.out/'runner-config.json'
    config_path.write_text(json.dumps(metadata, indent=2)+'\n')

    def run(block):
        print('Starting offline audit block', block, flush=True)
        with (args.out/f'block{block}-controller.log').open('w') as log:
            code = subprocess.run(commands[block], stdout=log, stderr=subprocess.STDOUT).returncode
        print('Finished offline audit block', block, 'exit', code, flush=True)
        return block, code

    with ThreadPoolExecutor(max_workers=4) as pool:
        codes = dict(pool.map(run, commands))
    accounting = {}
    for block in commands:
        path = args.out/f'block{block}'/'accounting.json'
        if path.exists():
            records = json.loads(path.read_text())
            assert not (accounting.keys() & records.keys()), 'Duplicate audit process'
            accounting.update(records)
    (args.out/'accounting.json').write_text(json.dumps(accounting, indent=2)+'\n')
    metadata.update(finished_utc=datetime.datetime.now(datetime.timezone.utc).isoformat(), exits=codes)
    config_path.write_text(json.dumps(metadata, indent=2)+'\n')
    assert all(code == 0 for code in codes.values()) and len(accounting) == len(blocks)*len(variants), 'Audit incomplete; retain journals and controller logs'


if __name__ == '__main__':
    main()
