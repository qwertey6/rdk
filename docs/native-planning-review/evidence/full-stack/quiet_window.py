"""Conservative host admission and observation for the deferred experiment.

The monitor is identical across variants and stays outside planner code. ps CPU
readings are smoothed, not instantaneous; five-second observations cannot prove
the absence of interference. Keep an operator-declared quiet window as well.
"""
import json
import os
from pathlib import Path
import re
import subprocess
import time


class QuietWindow:
    def __init__(self, root, policy):
        self.root, self.policy = root, policy
        self.power = None
        self.interference = False

    def observe(self, name, planning_pid=None):
        start = time.monotonic()
        readings = {}
        for key, command in {
            'load': ['ps', '-Ao', 'pid,ppid,pcpu,comm'],
            'power': ['pmset', '-g', 'batt'],
            'thermal': ['pmset', '-g', 'therm'],
        }.items():
            try:
                p = subprocess.run(command, capture_output=True, text=True, timeout=10)
                readings[key] = dict(exit=p.returncode, output=p.stdout, error=p.stderr)
            except (OSError, subprocess.TimeoutExpired) as err:
                readings[key] = dict(unavailable=str(err))
        reasons = [key + ' observation unavailable' for key, r in readings.items() if r.get('exit') != 0]
        power_text = readings['power'].get('output', '')
        power = ('AC' if "'AC Power'" in power_text else
                 'battery' if "'Battery Power'" in power_text else None)
        if power is None:
            reasons.append('unknown power source')
        elif self.power is not None and power != self.power:
            reasons.append('power source changed within experiment')
        thermal = readings['thermal'].get('output', '')
        if any(float(n) < 100 for n in re.findall(r'(?:CPU_Speed_Limit|CPU_Scheduler_Limit)\s*=\s*([\d.]+)', thermal)):
            reasons.append('CPU thermal limit')
        processes = []
        for line in readings['load'].get('output', '').splitlines()[1:]:
            fields = line.split(None, 3)
            if len(fields) != 4:
                reasons.append('unparsed process observation')
                continue
            pid, parent, cpu, command = fields
            processes.append((int(pid), int(parent), float(cpu), command))
        excluded = {os.getpid(), planning_pid}
        # Exclude the benchmark's descendants and our short-lived observers.
        while True:
            children = {pid for pid, parent, _, _ in processes if parent in excluded}
            if children <= excluded:
                break
            excluded |= children
        other = [(pid, cpu, cmd) for pid, _, cpu, cmd in processes if pid not in excluded]
        total = sum(cpu for _, cpu, _ in other)
        if not processes:
            reasons.append('empty process observation')
        if total > self.policy['max_external_cpu_percent']:
            reasons.append('aggregate external CPU exceeds limit')
        active_tools = [dict(pid=pid, cpu=cpu, command=cmd) for pid, cpu, cmd in other
                        if cpu >= self.policy['active_tool_cpu_percent'] and
                        re.search(self.policy['competing_tool_pattern'], Path(cmd).name, re.I)]
        if active_tools:
            reasons.append('competing build/test tool active')
        record = dict(utc=time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime()),
                      planning_pid=planning_pid, power=power, external_cpu_percent=total,
                      active_tools=active_tools, reasons=reasons, observation=readings,
                      observation_wall_ms=1000*(time.monotonic()-start))
        with (self.root/(name+'-host.jsonl')).open('a') as stream:
            stream.write(json.dumps(record)+'\n')
        return record

    def mark_interference(self, name, record):
        self.interference = True
        path = self.root/'INTERFERENCE.json'
        if not path.exists():
            path.write_text(json.dumps(dict(process=name, first_observation=record,
                action='Finish active process, stop new launches, retain whole experiment; no selective reruns.'), indent=2)+'\n')

    def admit(self, name):
        deadline = time.monotonic()+self.policy['admission_timeout_seconds']
        consecutive, last_power = 0, None
        while True:
            record = self.observe(name)
            if self.power is not None and record['reasons']:
                self.mark_interference(name, record)
                raise SystemExit('Interference between processes; experiment retained and stopped.')
            if record['power'] != last_power:
                consecutive = 0
            consecutive = 0 if record['reasons'] else consecutive+1
            last_power = record['power']
            if consecutive >= self.policy['quiet_checks']:
                self.power = record['power']
                return
            if time.monotonic() >= deadline:
                self.mark_interference(name, record)
                raise SystemExit('Quiet-window admission unavailable; no process launched for '+name)
            print('Quiet-window admission', name, record['reasons'] or 'first quiet check', flush=True)
            time.sleep(self.policy['admission_interval_seconds'])

    def run(self, name, command, env, log):
        proc = subprocess.Popen(command, env=env, stdout=log, stderr=subprocess.STDOUT)
        while proc.poll() is None:
            record = self.observe(name, proc.pid)
            if record['reasons']:
                self.mark_interference(name, record)
            try:
                proc.wait(timeout=self.policy['monitor_interval_seconds'])
            except subprocess.TimeoutExpired:
                pass
        final = self.observe(name)
        if final['reasons']:
            self.mark_interference(name, final)
        return proc
