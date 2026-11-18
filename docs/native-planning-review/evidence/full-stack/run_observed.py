#!/usr/bin/env python3
"""Serial comparison under observed desktop activity; host readings never gate timing."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import subprocess
import time

HERE = Path(__file__).resolve().parent


def digest(path):
    sha = hashlib.sha256()
    with path.open('rb') as stream:
        for chunk in iter(lambda: stream.read(1024 * 1024), b''):
            sha.update(chunk)
    return sha.hexdigest()


def observe():
    result = {}
    for name, cmd in {
        'hardware': ['sysctl', '-n', 'machdep.cpu.brand_string', 'hw.memsize', 'hw.logicalcpu'],
        'os': ['sw_vers'], 'power': ['pmset', '-g', 'batt'],
        'thermal': ['pmset', '-g', 'therm'],
        'load': ['ps', '-Ao', 'pid,pcpu,pmem,etime,comm', '-r'],
    }.items():
        try:
            proc = subprocess.run(cmd, capture_output=True, text=True, timeout=10)
            result[name] = {'exit': proc.returncode, 'output': proc.stdout[:15000], 'error': proc.stderr}
        except (OSError, subprocess.TimeoutExpired) as err:
            result[name] = {'unavailable': str(err)}
    return result



from quiet_window import QuietWindow


class ObservedWindow(QuietWindow):
    """Record existing host activity without rejecting or dropping requests."""
    def admit(self, name):
        self.interference = False
        record = self.observe(name)
        if self.power is None:
            self.power = record['power']
        if record['reasons']:
            self.mark_interference(name, record)

    def mark_interference(self, name, record):
        self.interference = True
        path = self.root / 'INTERFERENCE.json'
        if not path.exists():
            path.write_text(json.dumps(dict(process=name, first_observation=record,
                action='Record existing host activity and continue all declared processes. Retain every outcome; no selective reruns.'), indent=2)+'\n')


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument('--out', type=Path, default=Path('/tmp/rdk-native-pr-cleanup'))
    ap.add_argument('--corpus', type=Path, default=Path('/tmp/rdk-motion-scenes'))
    ap.add_argument('--candidate', type=Path, required=True)
    ap.add_argument('--intermediate', type=Path, default=Path('/tmp/rdk-native-finishing-state/finishing-state-v1.test'))
    ap.add_argument('--intermediate-sha256', default='c31d375c9752bd4703a0d961f21320ac9a822d771f9c268cdcfa08a55029f865')
    ap.add_argument('--stage', choices=['pilot', 'confirmation'], required=True)
    ap.add_argument('--policy', choices=['fixed', 'production'], default='production')
    ap.add_argument('--protocol', type=Path, default=HERE / 'protocol.json')
    ap.add_argument('--dry-run', action='store_true', help='Verify hashes and print the schedule; launch no planning or host observations')
    args = ap.parse_args()
    protocol = json.loads(args.protocol.read_text())
    assert protocol.get('admission_mode') == 'record_only', 'Use the explicitly revised observed-load protocol'
    identity_path = args.protocol.with_name('identity.json')
    identity = json.loads(identity_path.read_text())
    assert digest(args.protocol) == identity['protocol_sha256'], 'Protocol changed after freeze'
    for name, expected_hash in identity.get('tool_sha256', {}).items():
        assert digest(HERE / name) == expected_hash, 'Experiment tool changed: ' + name
    auditor = identity['binaries']['frozen_audit']
    assert digest(Path(auditor['path'])) == auditor['sha256'], 'Frozen auditor changed'
    required_policy = protocol.get('run_policy', 'production')
    if args.stage != 'confirmation' or args.policy != required_policy:
        ap.error('Command must match the frozen protocol: confirmation / ' + required_policy)
    for name, expected_hash in protocol['source_sha256'].items():
        assert digest(Path(name)) == expected_hash, 'Source changed: ' + name
    for name, expected in protocol['input_sha256'].items():
        assert digest(args.corpus / name) == expected, name
    variants = {name: Path(record['path']) for name, record in identity['binaries'].items() if name != 'frozen_audit'}
    expected = {'original': '6df28e527991ccc3f14a234374039c82bbfc79dd1b198b3ca031852efd05e163'}
    hashes = {name: digest(binary) for name, binary in variants.items()}
    assert hashes == protocol['binary_sha256'], 'Frozen binary mismatch'
    for name, sha in expected.items():
        assert hashes[name] == sha, f'{name} binary changed'
    orders = [list(variants)] if args.stage == 'pilot' else protocol['variant_order']
    if args.dry_run:
        print(json.dumps(dict(verified_source_files=len(protocol['source_sha256']),
                              verified_input_files=len(protocol['input_sha256']),
                              binaries={k: dict(path=str(v), sha256=hashes[k]) for k, v in variants.items()},
                              policy=args.policy, variant_order=orders,
                              measured_requests=sum(map(len, orders))*894,
                              including_warmup=sum(map(len, orders))*1004,
                              output=str(args.out / args.stage), planning_processes_launched=0), indent=2))
        return
    root = args.out / args.stage
    root.mkdir(parents=True, exist_ok=False)  # A fresh experiment never overwrites prior evidence.
    quiet = None
    if 'quiet_window' in protocol:
        quiet = ObservedWindow(root, protocol['quiet_window'])
    (root / 'protocol.json').write_text(args.protocol.read_text())
    for block, order in enumerate(orders, 1):
        for variant in order:
            name = f'{block}-{variant}'
            if quiet:
                quiet.admit(name)
            else:
                raise SystemExit('Host observation policy missing')
            run = root / name
            run.mkdir()  # Refuse overwrite, including interrupted runs.
            env = dict(os.environ)
            for key in list(env):
                if key.startswith(('RDK_NATIVE_', 'RDK_PLANNER_', 'RDK_IK_')) or key in ['MP_NUM_THREADS', 'CBIRRT_RACE_ATTEMPTS']:
                    del env[key]
            env.update(RDK_PLANNER_SUITE=str(args.corpus), RDK_PLANNER_SUITE_BACKEND='native',
                       RDK_PLANNER_SUITE_OUTPUT=str(run), MOTION_ROADMAP_CACHE_DIR='',
                       RDK_PLANNER_BENCHMARK_RUNS='12' if args.stage == 'pilot' else '50',
                       RDK_PLANNER_SUITE_TIMEOUT='60')
            if quiet:
                for key in protocol['environment_unset']:
                    env.pop(key, None)
                env.update(protocol['environment_set'])
            if variant == 'original':
                env['RDK_NATIVE_KINEMATICS'] = 'native'
            cmd = ['/usr/bin/time', '-l', str(variants[variant]), '-test.run', '^Test(Native|VAMP)(RecordedBenchmarkSuite|OrderedBenchmarkSequence)$',
                   '-test.count=1', '-test.v', '-test.timeout=90m']
            config = dict(variant=variant, block=block, policy=args.policy, command=cmd, binary_sha256=hashes[variant],
                          environment={k: v for k, v in env.items() if k.startswith(('RDK_', 'MOTION_')) or
                                       k in ['GOGC', 'GOMEMLIMIT', 'GOMAXPROCS', 'MP_NUM_THREADS', 'CBIRRT_RACE_ATTEMPTS']},
                          before=observe(), started_utc=time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime()))
            config_path = root / (name + '-config.json')
            config_path.write_text(json.dumps(config, indent=2) + '\n')
            print(f'Starting {args.stage} {name}', flush=True)
            with (root / (name + '.log')).open('w') as log:
                if quiet:
                    proc = quiet.run(name, cmd, env, log)
                else:
                    proc = subprocess.run(cmd, env=env, stdout=log, stderr=subprocess.STDOUT)
            config.update(exit=proc.returncode, after=observe(), finished_utc=time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime()))
            if quiet:
                config['environment_interference'] = quiet.interference
            config_path.write_text(json.dumps(config, indent=2) + '\n')
            print(f'Finished {args.stage} {name}, exit={proc.returncode}', flush=True)
            # Preserve a failed run in place; don't shrink its denominator or start
            # interpreting partial output as a completed latency experiment.
            if proc.returncode:
                raise SystemExit(proc.returncode)


if __name__ == '__main__':
    main()
