# Sampling-policy comparison: detailed distributions

**Environmental interference was recorded. Retain this entire cohort as contaminated; do not describe it as a clean confirmation.**

Stage: **confirmation**. Pilot results use one independent process per variant; this declared experiment uses 3.

ORDINARY PRODUCTION POLICY: all four variants use normal automatic startup calibration and default workers. Settings are recorded, not filtered. This measures the combined ordinary-policy effect; it does not hold calibrated budgets constant to estimate pure cost per check. Separate from every prior experiment; do not pool.

Original native remains visible first. Variants share the declared request/seed schedule. The sampling-only contrast removes the added angular cap; the literal PRs 2–4 variant also removes PR1’s other changes. Calibration follows ordinary production policy, as stated above.

Planning replay measures the actual PlanMotion call, including IK, search, required validation, finishing and synchronous cleanup. Original native does not join all losing search workers; its complete shutdown time is unavailable. The full and coarse_full variants join search workers; original and pr234 do not, and their asynchronous shutdown time is unavailable. The shipping branches omit benchmark-only diagnostics; absent phase and cache counters are unavailable. Snapshot preparation, RPC, hardware execution, warmup and the independent audit are excluded. These are CPU results on a shared workstation.

Scene aggregates weight the 11 distinct files equally. Scene 08 duplicates scene 02 and remains in the ledger with zero headline weight. The complete ordered workload is separate. Each process runs scenes, then one complete sequence warmup and three measured passes; its internal cache history is preserved.

All completion latencies include failures; accepted-plan latency is conditional on independent audit success. Median averages the two middle observations when needed; p95/p99 use nearest rank on pooled raw observations. Timeout latency includes cleanup and is not a successful time-to-solution.

## Scenes: all-outcome completion latency

| Variant / process | N | Accepted | Total s | Median ms | p95 ms | p99 ms | Max ms | ≥1/5/30 s |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| Original native | 330 | 311 (94.24%) | 28.793 | 63.32 | 224.49 | 291.59 | 362.14 | 0/0/0 |
| Full stack · fine grid | 330 | 330 (100.00%) | 27.371 | 53.85 | 237.20 | 267.87 | 285.63 | 0/0/0 |
| Full stack · original grid | 330 | 326 (98.79%) | 20.318 | 50.56 | 162.57 | 183.38 | 213.13 | 0/0/0 |
| PRs 2–4 only | 330 | 302 (91.52%) | 20.614 | 49.82 | 168.81 | 189.28 | 216.24 | 0/0/0 |

Full stack · fine grid against Original native: total-time ratio 0.951, exploratory 95% interval [0.917, 0.984]; mean difference -4.31 ms/request [-7.37, -1.35].

Full stack · original grid against Original native: total-time ratio 0.706, exploratory 95% interval [0.692, 0.733]; mean difference -25.68 ms/request [-27.29, -22.84].

PRs 2–4 only against Original native: total-time ratio 0.716, exploratory 95% interval [0.694, 0.747]; mean difference -24.78 ms/request [-26.77, -21.67].

Full stack · original grid against Full stack · fine grid: total-time ratio 0.742, exploratory 95% interval [0.727, 0.755]; mean difference -21.37 ms/request [-22.71, -19.92].

PRs 2–4 only against Full stack · fine grid: total-time ratio 0.753, exploratory 95% interval [0.729, 0.772]; mean difference -20.48 ms/request [-22.57, -18.54].

PRs 2–4 only against Full stack · original grid: total-time ratio 1.015, exploratory 95% interval [1.002, 1.023]; mean difference 0.90 ms/request [0.15, 1.38].

Accepted-only percentiles (N, accepted, total and threshold columns still describe all attempts):

| Variant / process | N | Accepted | Total s | Median ms | p95 ms | p99 ms | Max ms | ≥1/5/30 s |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| Original native | 330 | 311 (94.24%) | 28.793 | 64.07 | 224.49 | 291.59 | 362.14 | 0/0/0 |
| Full stack · fine grid | 330 | 330 (100.00%) | 27.371 | 53.85 | 237.20 | 267.87 | 285.63 | 0/0/0 |
| Full stack · original grid | 330 | 326 (98.79%) | 20.318 | 50.56 | 162.57 | 183.38 | 213.13 | 0/0/0 |
| PRs 2–4 only | 330 | 302 (91.52%) | 20.614 | 56.29 | 168.98 | 189.28 | 216.24 | 0/0/0 |

Each independent process:

| Variant / process | N | Accepted | Total s | Median ms | p95 ms | p99 ms | Max ms | ≥1/5/30 s |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| Original native / 1 | 110 | 100 (90.91%) | 9.758 | 66.45 | 220.41 | 248.96 | 295.35 | 0/0/0 |
| Original native / 2 | 110 | 102 (92.73%) | 9.405 | 62.54 | 219.03 | 291.59 | 362.14 | 0/0/0 |
| Original native / 3 | 110 | 109 (99.09%) | 9.630 | 62.01 | 228.48 | 281.77 | 308.34 | 0/0/0 |
| Full stack · fine grid / 1 | 110 | 110 (100.00%) | 8.947 | 53.20 | 227.88 | 267.39 | 285.63 | 0/0/0 |
| Full stack · fine grid / 2 | 110 | 110 (100.00%) | 9.256 | 54.05 | 244.85 | 267.87 | 283.24 | 0/0/0 |
| Full stack · fine grid / 3 | 110 | 110 (100.00%) | 9.168 | 54.51 | 227.64 | 255.81 | 269.94 | 0/0/0 |
| Full stack · original grid / 1 | 110 | 110 (100.00%) | 6.756 | 50.73 | 160.39 | 181.54 | 189.84 | 0/0/0 |
| Full stack · original grid / 2 | 110 | 107 (97.27%) | 6.893 | 52.12 | 171.59 | 183.38 | 213.13 | 0/0/0 |
| Full stack · original grid / 3 | 110 | 109 (99.09%) | 6.669 | 48.88 | 164.27 | 178.76 | 197.59 | 0/0/0 |
| PRs 2–4 only / 1 | 110 | 110 (100.00%) | 6.908 | 52.07 | 168.18 | 178.28 | 187.65 | 0/0/0 |
| PRs 2–4 only / 2 | 110 | 95 (86.36%) | 7.021 | 48.66 | 179.73 | 199.26 | 216.24 | 0/0/0 |
| PRs 2–4 only / 3 | 110 | 97 (88.18%) | 6.685 | 52.72 | 164.20 | 171.98 | 174.92 | 0/0/0 |

Per-input distributions and path lengths:

| Input | Variant | N / accepted | Total s | Median / p95 / p99 / max ms | Path length median / p95 / max | Waypoints median / p95 |
|---|---|---:|---:|---|---|---|
| scene01.json | Original native | 30 / 30 | 0.083 | 2.64 / 4.05 / 6.64 / 6.64 | 1.3125 / 1.3125 / 1.3125 | 2 / 2 |
| scene01.json | Full stack · fine grid | 30 / 30 | 0.078 | 2.74 / 3.23 / 3.87 / 3.87 | 1.3125 / 1.3125 / 1.3125 | 2 / 2 |
| scene01.json | Full stack · original grid | 30 / 30 | 0.062 | 2.00 / 2.66 / 3.46 / 3.46 | 1.3125 / 1.3125 / 1.3125 | 2 / 2 |
| scene01.json | PRs 2–4 only | 30 / 30 | 0.065 | 2.08 / 2.79 / 4.39 / 4.39 | 1.3125 / 1.3125 / 1.3125 | 2 / 2 |
| scene02.json | Original native | 30 / 30 | 0.256 | 8.22 / 10.62 / 10.97 / 10.97 | 0.9448 / 1.0164 / 1.0164 | 2 / 2 |
| scene02.json | Full stack · fine grid | 30 / 30 | 0.266 | 8.70 / 10.55 / 10.94 / 10.94 | 0.9448 / 1.0164 / 1.0164 | 2 / 2 |
| scene02.json | Full stack · original grid | 30 / 30 | 0.237 | 7.87 / 8.56 / 9.02 / 9.02 | 0.9806 / 1.0164 / 1.0164 | 2 / 2 |
| scene02.json | PRs 2–4 only | 30 / 30 | 0.241 | 7.92 / 8.74 / 10.11 / 10.11 | 0.9448 / 1.0164 / 1.0164 | 2 / 2 |
| scene03.json | Original native | 30 / 30 | 0.106 | 3.47 / 4.72 / 5.00 / 5.00 | 6.1254 / 6.1254 / 6.1254 | 2 / 2 |
| scene03.json | Full stack · fine grid | 30 / 30 | 0.417 | 13.47 / 15.78 / 16.15 / 16.15 | 6.1254 / 6.1254 / 6.1254 | 2 / 2 |
| scene03.json | Full stack · original grid | 30 / 30 | 0.101 | 3.38 / 4.32 / 4.42 / 4.42 | 6.1254 / 6.1254 / 6.1254 | 2 / 2 |
| scene03.json | PRs 2–4 only | 30 / 30 | 0.101 | 3.28 / 4.48 / 5.10 / 5.10 | 6.1254 / 6.1254 / 6.1254 | 2 / 2 |
| scene04.json | Original native | 30 / 23 | 1.133 | 35.24 / 69.12 / 75.21 / 75.21 | 1.8448 / 1.9454 / 2.1509 | 5 / 6 |
| scene04.json | Full stack · fine grid | 30 / 30 | 1.522 | 47.25 / 78.44 / 79.71 / 79.71 | 1.8271 / 2.3016 / 2.4601 | 6 / 8 |
| scene04.json | Full stack · original grid | 30 / 28 | 0.658 | 19.15 / 40.87 / 41.29 / 41.29 | 1.4309 / 2.1509 / 2.1509 | 5 / 6 |
| scene04.json | PRs 2–4 only | 30 / 22 | 0.706 | 22.54 / 44.54 / 45.05 / 45.05 | 1.8650 / 2.2551 / 2.5737 | 5 / 9 |
| scene05.json | Original native | 30 / 30 | 6.070 | 201.29 / 247.55 / 291.59 / 291.59 | 6.1995 / 9.1883 / 9.2165 | 8 / 15 |
| scene05.json | Full stack · fine grid | 30 / 30 | 6.025 | 215.46 / 255.81 / 285.63 / 285.63 | 6.0790 / 8.0682 / 9.2436 | 11 / 25 |
| scene05.json | Full stack · original grid | 30 / 29 | 4.375 | 152.77 / 197.59 / 213.13 / 213.13 | 5.9512 / 8.4413 / 9.2802 | 9 / 12 |
| scene05.json | PRs 2–4 only | 30 / 30 | 4.353 | 148.76 / 185.05 / 199.26 / 199.26 | 6.3038 / 8.3447 / 14.4275 | 8 / 20 |
| scene06.json | Original native | 30 / 28 | 6.345 | 198.41 / 308.34 / 362.14 / 362.14 | 6.1415 / 8.9506 / 9.0695 | 8 / 27 |
| scene06.json | Full stack · fine grid | 30 / 30 | 6.350 | 202.65 / 269.94 / 283.24 / 283.24 | 6.3872 / 9.0670 / 9.1801 | 14 / 166 |
| scene06.json | Full stack · original grid | 30 / 29 | 4.292 | 147.48 / 183.38 / 189.84 / 189.84 | 6.1708 / 8.9834 / 9.2405 | 9 / 35 |
| scene06.json | PRs 2–4 only | 30 / 30 | 4.468 | 152.13 / 189.56 / 216.24 / 216.24 | 6.1758 / 8.6750 / 9.2623 | 9 / 35 |
| scene07.json | Original native | 30 / 20 | 1.639 | 53.34 / 60.54 / 66.75 / 66.75 | 1.7816 / 1.8621 / 1.8621 | 6 / 6 |
| scene07.json | Full stack · fine grid | 30 / 30 | 1.577 | 51.33 / 63.96 / 64.21 / 64.21 | 1.3608 / 1.5677 / 1.5677 | 7 / 7 |
| scene07.json | Full stack · original grid | 30 / 30 | 1.319 | 41.84 / 50.59 / 73.06 / 73.06 | 1.4540 / 1.4540 / 1.4540 | 5 / 5 |
| scene07.json | PRs 2–4 only | 30 / 10 | 1.314 | 43.34 / 49.93 / 51.41 / 51.41 | 1.8621 / 1.8621 / 1.8621 | 6 / 6 |
| scene09.json | Original native | 30 / 30 | 2.326 | 68.62 / 121.29 / 147.46 / 147.46 | 4.3642 / 5.9288 / 5.9288 | 4 / 5 |
| scene09.json | Full stack · fine grid | 30 / 30 | 3.228 | 92.85 / 227.64 / 244.85 / 244.85 | 4.5423 / 7.1564 / 7.1564 | 5 / 6 |
| scene09.json | Full stack · original grid | 30 / 30 | 2.525 | 82.45 / 146.24 / 149.50 / 149.50 | 4.6245 / 6.4407 / 6.4407 | 4.5 / 8 |
| scene09.json | PRs 2–4 only | 30 / 30 | 2.515 | 80.42 / 149.48 / 154.66 / 154.66 | 4.6319 / 6.4407 / 6.4407 | 4 / 8 |
| scene10.json | Original native | 30 / 30 | 5.781 | 193.76 / 237.00 / 252.89 / 252.89 | 12.6680 / 13.3344 / 13.3344 | 22 / 39 |
| scene10.json | Full stack · fine grid | 30 / 30 | 3.894 | 101.33 / 245.38 / 246.25 / 246.25 | 12.6661 / 13.0315 / 13.0315 | 47 / 77 |
| scene10.json | Full stack · original grid | 30 / 30 | 2.755 | 78.93 / 160.39 / 161.98 / 161.98 | 12.6661 / 13.4803 / 13.4803 | 22.5 / 27 |
| scene10.json | PRs 2–4 only | 30 / 30 | 2.753 | 79.29 / 160.27 / 176.33 / 176.33 | 12.6661 / 13.4803 / 13.4803 | 22.5 / 27 |
| scene11.json | Original native | 30 / 30 | 3.046 | 99.20 / 119.63 / 124.82 / 124.82 | 0.6235 / 0.6235 / 0.6235 | 196 / 196 |
| scene11.json | Full stack · fine grid | 30 / 30 | 2.427 | 80.11 / 89.76 / 90.62 / 90.62 | 0.6235 / 0.6235 / 0.6235 | 196 / 196 |
| scene11.json | Full stack · original grid | 30 / 30 | 2.433 | 80.50 / 90.09 / 91.71 / 91.71 | 0.6235 / 0.6235 / 0.6235 | 196 / 196 |
| scene11.json | PRs 2–4 only | 30 / 30 | 2.424 | 80.09 / 89.36 / 90.60 / 90.60 | 0.6235 / 0.6235 / 0.6235 | 196 / 196 |
| scene12.json | Original native | 30 / 30 | 2.008 | 64.29 / 86.82 / 96.89 / 96.89 | 0.8142 / 0.8142 / 0.8142 | 93 / 93 |
| scene12.json | Full stack · fine grid | 30 / 30 | 1.587 | 52.74 / 54.19 / 54.58 / 54.58 | 0.8142 / 0.8142 / 0.8142 | 93 / 93 |
| scene12.json | Full stack · original grid | 30 / 30 | 1.562 | 51.00 / 59.55 / 64.06 / 64.06 | 0.8142 / 0.8142 / 0.8142 | 93 / 93 |
| scene12.json | PRs 2–4 only | 30 / 30 | 1.673 | 57.06 / 68.01 / 70.56 / 70.56 | 0.8142 / 0.8142 / 0.8142 | 93 / 93 |

## Sequence: all-outcome completion latency

| Variant / process | N | Accepted | Total s | Median ms | p95 ms | p99 ms | Max ms | ≥1/5/30 s |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| Original native | 882 | 881 (99.89%) | 68.748 | 51.54 | 333.01 | 494.72 | 792.28 | 0/0/0 |
| Full stack · fine grid | 882 | 882 (100.00%) | 51.512 | 34.26 | 252.43 | 543.50 | 753.42 | 0/0/0 |
| Full stack · original grid | 882 | 881 (99.89%) | 42.460 | 28.16 | 208.45 | 374.92 | 600.55 | 0/0/0 |
| PRs 2–4 only | 882 | 881 (99.89%) | 42.998 | 28.25 | 225.45 | 396.96 | 872.00 | 0/0/0 |

Full stack · fine grid against Original native: total-time ratio 0.749, exploratory 95% interval [0.740, 0.761]; mean difference -19.54 ms/request [-20.75, -18.05].

Full stack · original grid against Original native: total-time ratio 0.618, exploratory 95% interval [0.606, 0.639]; mean difference -29.80 ms/request [-31.47, -27.30].

PRs 2–4 only against Original native: total-time ratio 0.625, exploratory 95% interval [0.608, 0.643]; mean difference -29.20 ms/request [-31.30, -27.01].

Full stack · original grid against Full stack · fine grid: total-time ratio 0.824, exploratory 95% interval [0.816, 0.839]; mean difference -10.26 ms/request [-10.82, -9.25].

PRs 2–4 only against Full stack · fine grid: total-time ratio 0.835, exploratory 95% interval [0.821, 0.844]; mean difference -9.65 ms/request [-10.56, -8.96].

PRs 2–4 only against Full stack · original grid: total-time ratio 1.013, exploratory 95% interval [1.003, 1.029]; mean difference 0.61 ms/request [0.17, 1.38].

Accepted-only percentiles (N, accepted, total and threshold columns still describe all attempts):

| Variant / process | N | Accepted | Total s | Median ms | p95 ms | p99 ms | Max ms | ≥1/5/30 s |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| Original native | 882 | 881 (99.89%) | 68.748 | 51.53 | 333.01 | 494.72 | 792.28 | 0/0/0 |
| Full stack · fine grid | 882 | 882 (100.00%) | 51.512 | 34.26 | 252.43 | 543.50 | 753.42 | 0/0/0 |
| Full stack · original grid | 882 | 881 (99.89%) | 42.460 | 28.15 | 208.45 | 374.92 | 600.55 | 0/0/0 |
| PRs 2–4 only | 882 | 881 (99.89%) | 42.998 | 28.22 | 225.45 | 396.96 | 872.00 | 0/0/0 |

Each independent process:

| Variant / process | N | Accepted | Total s | Median ms | p95 ms | p99 ms | Max ms | ≥1/5/30 s |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| Original native / 1 | 294 | 294 (100.00%) | 23.456 | 52.60 | 386.35 | 502.91 | 792.28 | 0/0/0 |
| Original native / 2 | 294 | 294 (100.00%) | 22.220 | 51.18 | 333.01 | 457.39 | 625.74 | 0/0/0 |
| Original native / 3 | 294 | 293 (99.66%) | 23.072 | 51.64 | 354.84 | 636.36 | 683.67 | 0/0/0 |
| Full stack · fine grid / 1 | 294 | 294 (100.00%) | 17.356 | 33.59 | 261.51 | 629.64 | 746.50 | 0/0/0 |
| Full stack · fine grid / 2 | 294 | 294 (100.00%) | 16.912 | 34.39 | 252.43 | 529.40 | 668.54 | 0/0/0 |
| Full stack · fine grid / 3 | 294 | 294 (100.00%) | 17.243 | 34.26 | 251.16 | 551.57 | 753.42 | 0/0/0 |
| Full stack · original grid / 1 | 294 | 294 (100.00%) | 14.204 | 28.26 | 216.25 | 374.81 | 600.55 | 0/0/0 |
| Full stack · original grid / 2 | 294 | 293 (99.66%) | 14.194 | 28.08 | 209.55 | 496.25 | 575.14 | 0/0/0 |
| Full stack · original grid / 3 | 294 | 294 (100.00%) | 14.062 | 28.08 | 207.72 | 488.70 | 543.95 | 0/0/0 |
| PRs 2–4 only / 1 | 294 | 294 (100.00%) | 14.253 | 28.21 | 257.03 | 425.01 | 429.43 | 0/0/0 |
| PRs 2–4 only / 2 | 294 | 293 (99.66%) | 14.278 | 28.42 | 219.50 | 396.96 | 547.37 | 0/0/0 |
| PRs 2–4 only / 3 | 294 | 294 (100.00%) | 14.467 | 28.29 | 243.75 | 451.22 | 872.00 | 0/0/0 |

Per-input distributions and path lengths:

| Input | Variant | N / accepted | Total s | Median / p95 / p99 / max ms | Path length median / p95 / max | Waypoints median / p95 |
|---|---|---:|---:|---|---|---|
| 000-153126.083-grinding-move.json | Original native | 9 / 9 | 0.127 | 14.19 / 15.77 / 15.77 / 15.77 | 2.0083 / 2.0083 / 2.0083 | 2 / 2 |
| 000-153126.083-grinding-move.json | Full stack · fine grid | 9 / 9 | 0.140 | 15.47 / 16.21 / 16.21 / 16.21 | 2.0083 / 2.0083 / 2.0083 | 2 / 2 |
| 000-153126.083-grinding-move.json | Full stack · original grid | 9 / 9 | 0.086 | 9.68 / 9.83 / 9.83 / 9.83 | 2.0083 / 2.0083 / 2.0083 | 2 / 2 |
| 000-153126.083-grinding-move.json | PRs 2–4 only | 9 / 9 | 0.087 | 9.75 / 10.31 / 10.31 / 10.31 | 2.0083 / 2.0083 / 2.0083 | 2 / 2 |
| 001-153129.698-grinding-move.json | Original native | 9 / 9 | 0.088 | 9.79 / 10.14 / 10.14 / 10.14 | 0.2171 / 0.2171 / 0.2171 | 8 / 8 |
| 001-153129.698-grinding-move.json | Full stack · fine grid | 9 / 9 | 0.079 | 8.76 / 9.56 / 9.56 / 9.56 | 0.2171 / 0.2171 / 0.2171 | 8 / 8 |
| 001-153129.698-grinding-move.json | Full stack · original grid | 9 / 9 | 0.076 | 8.48 / 8.72 / 8.72 / 8.72 | 0.2171 / 0.2171 / 0.2171 | 8 / 8 |
| 001-153129.698-grinding-move.json | PRs 2–4 only | 9 / 9 | 0.074 | 8.17 / 8.61 / 8.61 / 8.61 | 0.2171 / 0.2171 / 0.2171 | 8 / 8 |
| 002-153130.689-grinding-move.json | Original native | 9 / 9 | 0.080 | 8.87 / 9.28 / 9.28 / 9.28 | 0.2169 / 0.2169 / 0.2169 | 8 / 8 |
| 002-153130.689-grinding-move.json | Full stack · fine grid | 9 / 9 | 0.072 | 7.89 / 8.21 / 8.21 / 8.21 | 0.2169 / 0.2169 / 0.2169 | 8 / 8 |
| 002-153130.689-grinding-move.json | Full stack · original grid | 9 / 9 | 0.068 | 7.42 / 8.07 / 8.07 / 8.07 | 0.2169 / 0.2169 / 0.2169 | 8 / 8 |
| 002-153130.689-grinding-move.json | PRs 2–4 only | 9 / 9 | 0.067 | 7.36 / 7.80 / 7.80 / 7.80 | 0.2169 / 0.2169 / 0.2169 | 8 / 8 |
| 003-153131.847-grinding-circular.json | Original native | 9 / 9 | 0.187 | 20.57 / 22.48 / 22.48 / 22.48 | 0.2088 / 0.2088 / 0.2088 | 26 / 26 |
| 003-153131.847-grinding-circular.json | Full stack · fine grid | 9 / 9 | 0.155 | 17.41 / 18.16 / 18.16 / 18.16 | 0.2088 / 0.2088 / 0.2088 | 26 / 26 |
| 003-153131.847-grinding-circular.json | Full stack · original grid | 9 / 9 | 0.148 | 16.58 / 17.68 / 17.68 / 17.68 | 0.2088 / 0.2088 / 0.2088 | 26 / 26 |
| 003-153131.847-grinding-circular.json | PRs 2–4 only | 9 / 9 | 0.149 | 16.57 / 17.46 / 17.46 / 17.46 | 0.2088 / 0.2088 / 0.2088 | 26 / 26 |
| 004-153143.958-tamping-move.json | Original native | 9 / 9 | 1.121 | 123.99 / 133.90 / 133.90 / 133.90 | 1.3788 / 1.5598 / 1.5598 | 110 / 126 |
| 004-153143.958-tamping-move.json | Full stack · fine grid | 9 / 9 | 1.784 | 200.60 / 208.74 / 208.74 / 208.74 | 1.4037 / 1.6054 / 1.6054 | 615 / 664 |
| 004-153143.958-tamping-move.json | Full stack · original grid | 9 / 9 | 1.461 | 161.91 / 165.45 / 165.45 / 165.45 | 1.3788 / 1.4076 / 1.4076 | 110 / 120 |
| 004-153143.958-tamping-move.json | PRs 2–4 only | 9 / 9 | 1.544 | 171.17 / 185.87 / 185.87 / 185.87 | 1.3788 / 1.5598 / 1.5598 | 110 / 126 |
| 005-153147.607-tamping-move.json | Original native | 9 / 9 | 0.395 | 43.82 / 46.59 / 46.59 / 46.59 | 0.4603 / 0.4603 / 0.4603 | 52 / 52 |
| 005-153147.607-tamping-move.json | Full stack · fine grid | 9 / 9 | 0.329 | 36.18 / 37.96 / 37.96 / 37.96 | 0.4603 / 0.4603 / 0.4603 | 52 / 52 |
| 005-153147.607-tamping-move.json | Full stack · original grid | 9 / 9 | 0.313 | 34.95 / 36.51 / 36.51 / 36.51 | 0.4603 / 0.4603 / 0.4603 | 52 / 52 |
| 005-153147.607-tamping-move.json | PRs 2–4 only | 9 / 9 | 0.318 | 35.20 / 37.14 / 37.14 / 37.14 | 0.4603 / 0.4603 / 0.4603 | 52 / 52 |
| 006-153152.592-tamping-move.json | Original native | 9 / 9 | 0.404 | 44.27 / 46.88 / 46.88 / 46.88 | 0.4602 / 0.4602 / 0.4602 | 52 / 52 |
| 006-153152.592-tamping-move.json | Full stack · fine grid | 9 / 9 | 0.333 | 36.29 / 42.12 / 42.12 / 42.12 | 0.4602 / 0.4602 / 0.4602 | 52 / 52 |
| 006-153152.592-tamping-move.json | Full stack · original grid | 9 / 9 | 0.333 | 34.90 / 49.00 / 49.00 / 49.00 | 0.4602 / 0.4602 / 0.4602 | 52 / 52 |
| 006-153152.592-tamping-move.json | PRs 2–4 only | 9 / 9 | 0.312 | 34.48 / 36.40 / 36.40 / 36.40 | 0.4602 / 0.4602 / 0.4602 | 52 / 52 |
| 007-153209.673-locking_portafilter-move.json | Original native | 9 / 9 | 3.429 | 382.70 / 416.21 / 416.21 / 416.21 | 2.9016 / 3.1632 / 3.1632 | 276 / 318 |
| 007-153209.673-locking_portafilter-move.json | Full stack · fine grid | 9 / 9 | 3.259 | 356.88 / 401.85 / 401.85 / 401.85 | 2.7480 / 2.7918 / 2.7918 | 1010 / 1085 |
| 007-153209.673-locking_portafilter-move.json | Full stack · original grid | 9 / 9 | 2.712 | 292.90 / 333.93 / 333.93 / 333.93 | 2.9016 / 3.1632 / 3.1632 | 276 / 292 |
| 007-153209.673-locking_portafilter-move.json | PRs 2–4 only | 9 / 9 | 2.645 | 282.24 / 341.24 / 341.24 / 341.24 | 2.7500 / 3.1632 / 3.1632 | 276 / 333 |
| 008-153214.999-locking_portafilter-move.json | Original native | 9 / 9 | 0.192 | 21.25 / 21.77 / 21.77 / 21.77 | 0.1406 / 0.1406 / 0.1406 | 23 / 23 |
| 008-153214.999-locking_portafilter-move.json | Full stack · fine grid | 9 / 9 | 0.162 | 17.95 / 20.10 / 20.10 / 20.10 | 0.1406 / 0.1406 / 0.1406 | 23 / 23 |
| 008-153214.999-locking_portafilter-move.json | Full stack · original grid | 9 / 9 | 0.159 | 17.32 / 19.37 / 19.37 / 19.37 | 0.1406 / 0.1406 / 0.1406 | 23 / 23 |
| 008-153214.999-locking_portafilter-move.json | PRs 2–4 only | 9 / 9 | 0.160 | 17.88 / 18.21 / 18.21 / 18.21 | 0.1406 / 0.1406 / 0.1406 | 23 / 23 |
| 009-153216.395-locking_portafilter-pivot.json | Original native | 9 / 9 | 0.388 | 42.57 / 45.44 / 45.44 / 45.44 | 2.8613 / 2.8613 / 2.8613 | 40 / 40 |
| 009-153216.395-locking_portafilter-pivot.json | Full stack · fine grid | 9 / 9 | 0.387 | 42.83 / 46.13 / 46.13 / 46.13 | 2.8613 / 2.8613 / 2.8613 | 40 / 40 |
| 009-153216.395-locking_portafilter-pivot.json | Full stack · original grid | 9 / 9 | 0.304 | 33.39 / 36.31 / 36.31 / 36.31 | 2.8613 / 2.8613 / 2.8613 | 40 / 40 |
| 009-153216.395-locking_portafilter-pivot.json | PRs 2–4 only | 9 / 9 | 0.295 | 32.77 / 33.34 / 33.34 / 33.34 | 2.8613 / 2.8613 / 2.8613 | 40 / 40 |
| 010-153224.639-releasing_filter-move.json | Original native | 9 / 9 | 1.050 | 114.61 / 127.59 / 127.59 / 127.59 | 0.6268 / 0.6268 / 0.6268 | 62 / 62 |
| 010-153224.639-releasing_filter-move.json | Full stack · fine grid | 9 / 9 | 0.328 | 36.11 / 39.90 / 39.90 / 39.90 | 0.6268 / 0.6268 / 0.6268 | 62 / 62 |
| 010-153224.639-releasing_filter-move.json | Full stack · original grid | 9 / 9 | 0.319 | 35.11 / 37.15 / 37.15 / 37.15 | 0.6268 / 0.6268 / 0.6268 | 62 / 62 |
| 010-153224.639-releasing_filter-move.json | PRs 2–4 only | 9 / 9 | 0.320 | 35.43 / 38.04 / 38.04 / 38.04 | 0.6268 / 0.6268 / 0.6268 | 62 / 62 |
| 011-153227.696-placing_cup-move.json | Original native | 9 / 9 | 0.297 | 32.99 / 34.21 / 34.21 / 34.21 | 3.2841 / 3.2841 / 3.2841 | 2 / 2 |
| 011-153227.696-placing_cup-move.json | Full stack · fine grid | 9 / 9 | 0.350 | 38.73 / 40.73 / 40.73 / 40.73 | 3.2841 / 3.2841 / 3.2841 | 2 / 2 |
| 011-153227.696-placing_cup-move.json | Full stack · original grid | 9 / 9 | 0.255 | 28.18 / 29.75 / 29.75 / 29.75 | 3.2841 / 3.2841 / 3.2841 | 2 / 2 |
| 011-153227.696-placing_cup-move.json | PRs 2–4 only | 9 / 9 | 0.253 | 28.21 / 28.96 / 28.96 / 28.96 | 3.2841 / 3.2841 / 3.2841 | 2 / 2 |
| 012-153235.149-placing_cup-move.json | Original native | 9 / 9 | 0.201 | 22.33 / 23.13 / 23.13 / 23.13 | 3.7835 / 3.7835 / 3.7835 | 2 / 2 |
| 012-153235.149-placing_cup-move.json | Full stack · fine grid | 9 / 9 | 0.249 | 27.70 / 28.80 / 28.80 / 28.80 | 3.7835 / 3.7835 / 3.7835 | 2 / 2 |
| 012-153235.149-placing_cup-move.json | Full stack · original grid | 9 / 9 | 0.174 | 19.03 / 21.57 / 21.57 / 21.57 | 3.7835 / 3.7835 / 3.7835 | 2 / 2 |
| 012-153235.149-placing_cup-move.json | PRs 2–4 only | 9 / 9 | 0.167 | 18.46 / 19.33 / 19.33 / 19.33 | 3.7835 / 3.7835 / 3.7835 | 2 / 2 |
| 013-153235.711-placing_cup-move.json | Original native | 9 / 9 | 0.481 | 53.98 / 55.55 / 55.55 / 55.55 | 0.6327 / 0.6327 / 0.6327 | 69 / 69 |
| 013-153235.711-placing_cup-move.json | Full stack · fine grid | 9 / 9 | 0.385 | 42.41 / 46.33 / 46.33 / 46.33 | 0.6327 / 0.6327 / 0.6327 | 69 / 69 |
| 013-153235.711-placing_cup-move.json | Full stack · original grid | 9 / 9 | 0.373 | 41.23 / 42.53 / 42.53 / 42.53 | 0.6327 / 0.6327 / 0.6327 | 69 / 69 |
| 013-153235.711-placing_cup-move.json | PRs 2–4 only | 9 / 9 | 0.374 | 41.42 / 43.30 / 43.30 / 43.30 | 0.6327 / 0.6327 / 0.6327 | 69 / 69 |
| 014-153244.672-placing_cup-move.json | Original native | 9 / 9 | 0.514 | 56.82 / 59.64 / 59.64 / 59.64 | 0.6326 / 0.6326 / 0.6326 | 69 / 69 |
| 014-153244.672-placing_cup-move.json | Full stack · fine grid | 9 / 9 | 0.411 | 45.32 / 49.30 / 49.30 / 49.30 | 0.6326 / 0.6326 / 0.6326 | 69 / 69 |
| 014-153244.672-placing_cup-move.json | Full stack · original grid | 9 / 9 | 0.402 | 43.61 / 52.20 / 52.20 / 52.20 | 0.6326 / 0.6326 / 0.6326 | 69 / 69 |
| 014-153244.672-placing_cup-move.json | PRs 2–4 only | 9 / 9 | 0.402 | 44.24 / 47.32 / 47.32 / 47.32 | 0.6326 / 0.6326 / 0.6326 | 69 / 69 |
| 015-153257.279-placing_cup-move.json | Original native | 9 / 9 | 3.530 | 406.42 / 448.32 / 448.32 / 448.32 | 5.2702 / 5.7154 / 5.7154 | 699 / 747 |
| 015-153257.279-placing_cup-move.json | Full stack · fine grid | 9 / 9 | 3.483 | 383.02 / 551.57 / 551.57 / 551.57 | 5.1572 / 5.6196 / 5.6196 | 2100 / 2282 |
| 015-153257.279-placing_cup-move.json | Full stack · original grid | 9 / 8 | 2.110 | 223.21 / 284.51 / 284.51 / 284.51 | 5.0597 / 5.7142 / 5.7142 | 693.5 / 738 |
| 015-153257.279-placing_cup-move.json | PRs 2–4 only | 9 / 8 | 2.260 | 257.03 / 279.31 / 279.31 / 279.31 | 5.3274 / 5.5512 / 5.5512 | 675 / 699 |
| 016-153310.154-placing_cup-move.json | Original native | 9 / 9 | 0.641 | 71.28 / 72.90 / 72.90 / 72.90 | 0.6053 / 0.6053 / 0.6053 | 92 / 92 |
| 016-153310.154-placing_cup-move.json | Full stack · fine grid | 9 / 9 | 0.511 | 56.55 / 60.54 / 60.54 / 60.54 | 0.6053 / 0.6053 / 0.6053 | 92 / 92 |
| 016-153310.154-placing_cup-move.json | Full stack · original grid | 9 / 9 | 0.503 | 55.63 / 58.68 / 58.68 / 58.68 | 0.6053 / 0.6053 / 0.6053 | 92 / 92 |
| 016-153310.154-placing_cup-move.json | PRs 2–4 only | 9 / 9 | 0.494 | 54.64 / 56.62 / 56.62 / 56.62 | 0.6053 / 0.6053 / 0.6053 | 92 / 92 |
| 017-153314.165-placing_cup-move.json | Original native | 9 / 9 | 0.581 | 64.57 / 65.87 / 65.87 / 65.87 | 0.6051 / 0.6051 / 0.6051 | 91 / 91 |
| 017-153314.165-placing_cup-move.json | Full stack · fine grid | 9 / 9 | 0.460 | 50.06 / 57.65 / 57.65 / 57.65 | 0.6051 / 0.6051 / 0.6051 | 91 / 91 |
| 017-153314.165-placing_cup-move.json | Full stack · original grid | 9 / 9 | 0.462 | 49.53 / 63.36 / 63.36 / 63.36 | 0.6051 / 0.6051 / 0.6051 | 91 / 91 |
| 017-153314.165-placing_cup-move.json | PRs 2–4 only | 9 / 9 | 0.449 | 49.77 / 53.47 / 53.47 / 53.47 | 0.6051 / 0.6051 / 0.6051 | 91 / 91 |
| 018-153332.655-brewing-move.json | Original native | 9 / 9 | 1.008 | 108.86 / 133.48 / 133.48 / 133.48 | 1.6383 / 1.6922 / 1.6922 | 195 / 200 |
| 018-153332.655-brewing-move.json | Full stack · fine grid | 9 / 9 | 0.902 | 100.94 / 105.04 / 105.04 / 105.04 | 1.6413 / 1.7484 / 1.7484 | 589 / 592 |
| 018-153332.655-brewing-move.json | Full stack · original grid | 9 / 9 | 0.639 | 71.17 / 75.54 / 75.54 / 75.54 | 1.6383 / 1.6922 / 1.6922 | 195 / 200 |
| 018-153332.655-brewing-move.json | PRs 2–4 only | 9 / 9 | 0.641 | 69.98 / 76.61 / 76.61 / 76.61 | 1.6383 / 1.6922 / 1.6922 | 195 / 200 |
| 019-153336.155-brewing-move.json | Original native | 9 / 9 | 0.130 | 14.42 / 16.39 / 16.39 / 16.39 | 0.1464 / 0.1464 / 0.1464 | 15 / 15 |
| 019-153336.155-brewing-move.json | Full stack · fine grid | 9 / 9 | 0.108 | 11.90 / 12.84 / 12.84 / 12.84 | 0.1464 / 0.1464 / 0.1464 | 15 / 15 |
| 019-153336.155-brewing-move.json | Full stack · original grid | 9 / 9 | 0.107 | 11.85 / 12.23 / 12.23 / 12.23 | 0.1464 / 0.1464 / 0.1464 | 15 / 15 |
| 019-153336.155-brewing-move.json | PRs 2–4 only | 9 / 9 | 0.104 | 11.55 / 11.87 / 11.87 / 11.87 | 0.1464 / 0.1464 / 0.1464 | 15 / 15 |
| 020-153337.336-brewing-move.json | Original native | 9 / 9 | 0.128 | 14.15 / 14.55 / 14.55 / 14.55 | 0.1464 / 0.1464 / 0.1464 | 15 / 15 |
| 020-153337.336-brewing-move.json | Full stack · fine grid | 9 / 9 | 0.108 | 11.96 / 13.59 / 13.59 / 13.59 | 0.1464 / 0.1464 / 0.1464 | 15 / 15 |
| 020-153337.336-brewing-move.json | Full stack · original grid | 9 / 9 | 0.106 | 11.59 / 12.75 / 12.75 / 12.75 | 0.1464 / 0.1464 / 0.1464 | 15 / 15 |
| 020-153337.336-brewing-move.json | PRs 2–4 only | 9 / 9 | 0.104 | 11.51 / 11.77 / 11.77 / 11.77 | 0.1464 / 0.1464 / 0.1464 | 15 / 15 |
| 021-153338.286-brewing-move.json | Original native | 9 / 9 | 0.174 | 19.11 / 20.41 / 20.41 / 20.41 | 3.8350 / 3.8350 / 3.8350 | 2 / 2 |
| 021-153338.286-brewing-move.json | Full stack · fine grid | 9 / 9 | 0.252 | 27.61 / 29.67 / 29.67 / 29.67 | 3.8350 / 3.8350 / 3.8350 | 2 / 2 |
| 021-153338.286-brewing-move.json | Full stack · original grid | 9 / 9 | 0.139 | 15.40 / 16.40 / 16.40 / 16.40 | 3.8350 / 3.8350 / 3.8350 | 2 / 2 |
| 021-153338.286-brewing-move.json | PRs 2–4 only | 9 / 9 | 0.140 | 15.67 / 16.61 / 16.61 / 16.61 | 3.8350 / 3.8350 / 3.8350 | 2 / 2 |
| 022-153346.202-brewing-move.json | Original native | 9 / 9 | 0.089 | 9.87 / 10.59 / 10.59 / 10.59 | 1.0812 / 1.0812 / 1.0812 | 2 / 2 |
| 022-153346.202-brewing-move.json | Full stack · fine grid | 9 / 9 | 0.085 | 9.39 / 9.86 / 9.86 / 9.86 | 1.0812 / 1.0812 / 1.0812 | 2 / 2 |
| 022-153346.202-brewing-move.json | Full stack · original grid | 9 / 9 | 0.079 | 8.90 / 9.22 / 9.22 / 9.22 | 1.0812 / 1.0812 / 1.0812 | 2 / 2 |
| 022-153346.202-brewing-move.json | PRs 2–4 only | 9 / 9 | 0.081 | 9.08 / 9.79 / 9.79 / 9.79 | 1.0812 / 1.0812 / 1.0812 | 2 / 2 |
| 023-153346.645-brewing-move.json | Original native | 9 / 9 | 0.385 | 42.38 / 44.15 / 44.15 / 44.15 | 0.5861 / 0.5861 / 0.5861 | 61 / 61 |
| 023-153346.645-brewing-move.json | Full stack · fine grid | 9 / 9 | 0.313 | 34.52 / 36.05 / 36.05 / 36.05 | 0.5861 / 0.5861 / 0.5861 | 61 / 61 |
| 023-153346.645-brewing-move.json | Full stack · original grid | 9 / 9 | 0.299 | 33.05 / 35.40 / 35.40 / 35.40 | 0.5861 / 0.5861 / 0.5861 | 61 / 61 |
| 023-153346.645-brewing-move.json | PRs 2–4 only | 9 / 9 | 0.305 | 33.48 / 37.01 / 37.01 / 37.01 | 0.5861 / 0.5861 / 0.5861 | 61 / 61 |
| 024-153352.938-brewing-move.json | Original native | 9 / 9 | 0.427 | 47.56 / 48.63 / 48.63 / 48.63 | 0.5859 / 0.5859 / 0.5859 | 61 / 61 |
| 024-153352.938-brewing-move.json | Full stack · fine grid | 9 / 9 | 0.337 | 37.62 / 38.28 / 38.28 / 38.28 | 0.5859 / 0.5859 / 0.5859 | 61 / 61 |
| 024-153352.938-brewing-move.json | Full stack · original grid | 9 / 9 | 0.334 | 36.78 / 38.28 / 38.28 / 38.28 | 0.5859 / 0.5859 / 0.5859 | 61 / 61 |
| 024-153352.938-brewing-move.json | PRs 2–4 only | 9 / 9 | 0.331 | 36.76 / 37.65 / 37.65 / 37.65 | 0.5859 / 0.5859 / 0.5859 | 61 / 61 |
| 025-153404.593-brewing-move.json | Original native | 9 / 9 | 4.727 | 456.56 / 792.28 / 792.28 / 792.28 | 3.5655 / 8.2199 / 8.2199 | 585 / 802 |
| 025-153404.593-brewing-move.json | Full stack · fine grid | 9 / 9 | 4.850 | 486.19 / 753.42 / 753.42 / 753.42 | 3.5149 / 8.3303 / 8.3303 | 1688 / 3606 |
| 025-153404.593-brewing-move.json | Full stack · original grid | 9 / 9 | 4.571 | 506.32 / 600.55 / 600.55 / 600.55 | 8.3303 / 8.3303 / 8.3303 | 837 / 837 |
| 025-153404.593-brewing-move.json | PRs 2–4 only | 9 / 9 | 4.175 | 425.93 / 872.00 / 872.00 / 872.00 | 3.5306 / 8.3303 / 8.3303 | 618 / 841 |
| 026-153415.753-brewing-move.json | Original native | 9 / 9 | 0.401 | 44.66 / 46.18 / 46.18 / 46.18 | 0.4528 / 0.4528 / 0.4528 | 61 / 61 |
| 026-153415.753-brewing-move.json | Full stack · fine grid | 9 / 9 | 0.320 | 35.58 / 36.76 / 36.76 / 36.76 | 0.4528 / 0.4528 / 0.4528 | 61 / 61 |
| 026-153415.753-brewing-move.json | Full stack · original grid | 9 / 9 | 0.323 | 34.78 / 42.55 / 42.55 / 42.55 | 0.4528 / 0.4528 / 0.4528 | 61 / 61 |
| 026-153415.753-brewing-move.json | PRs 2–4 only | 9 / 9 | 0.320 | 35.31 / 37.33 / 37.33 / 37.33 | 0.4528 / 0.4528 / 0.4528 | 61 / 61 |
| 027-153433.046-brewing-move.json | Original native | 9 / 9 | 0.405 | 45.30 / 45.92 / 45.92 / 45.92 | 0.4526 / 0.4526 / 0.4526 | 61 / 61 |
| 027-153433.046-brewing-move.json | Full stack · fine grid | 9 / 9 | 0.330 | 35.57 / 47.64 / 47.64 / 47.64 | 0.4526 / 0.4526 / 0.4526 | 61 / 61 |
| 027-153433.046-brewing-move.json | Full stack · original grid | 9 / 9 | 0.314 | 34.88 / 36.08 / 36.08 / 36.08 | 0.4526 / 0.4526 / 0.4526 | 61 / 61 |
| 027-153433.046-brewing-move.json | PRs 2–4 only | 9 / 9 | 0.318 | 34.88 / 38.06 / 38.06 / 38.06 | 0.4526 / 0.4526 / 0.4526 | 61 / 61 |
| 028-153435.099-brewing-carry.json | Original native | 9 / 9 | 0.267 | 27.25 / 36.58 / 36.58 / 36.58 | 2.6357 / 2.6357 / 2.6357 | 6 / 6 |
| 028-153435.099-brewing-carry.json | Full stack · fine grid | 9 / 9 | 0.326 | 29.01 / 54.16 / 54.16 / 54.16 | 2.6357 / 2.6357 / 2.6357 | 6 / 6 |
| 028-153435.099-brewing-carry.json | Full stack · original grid | 9 / 9 | 0.208 | 22.64 / 29.21 / 29.21 / 29.21 | 2.6357 / 2.6357 / 2.6357 | 6 / 6 |
| 028-153435.099-brewing-carry.json | PRs 2–4 only | 9 / 9 | 0.205 | 21.82 / 27.08 / 27.08 / 27.08 | 2.6357 / 2.6357 / 2.6357 | 6 / 6 |
| 029-153439.206-brewing-move.json | Original native | 9 / 9 | 0.486 | 53.39 / 56.83 / 56.83 / 56.83 | 0.7636 / 0.7636 / 0.7636 | 69 / 69 |
| 029-153439.206-brewing-move.json | Full stack · fine grid | 9 / 9 | 0.396 | 43.26 / 50.61 / 50.61 / 50.61 | 0.7636 / 0.7636 / 0.7636 | 69 / 69 |
| 029-153439.206-brewing-move.json | Full stack · original grid | 9 / 9 | 0.369 | 41.22 / 42.23 / 42.23 / 42.23 | 0.7636 / 0.7636 / 0.7636 | 69 / 69 |
| 029-153439.206-brewing-move.json | PRs 2–4 only | 9 / 9 | 0.373 | 41.19 / 42.61 / 42.61 / 42.61 | 0.7636 / 0.7636 / 0.7636 | 69 / 69 |
| 030-153443.810-brewing-move.json | Original native | 9 / 9 | 1.088 | 121.17 / 124.15 / 124.15 / 124.15 | 0.7635 / 0.7635 / 0.7635 | 69 / 69 |
| 030-153443.810-brewing-move.json | Full stack · fine grid | 9 / 9 | 0.364 | 39.61 / 43.00 / 43.00 / 43.00 | 0.7635 / 0.7635 / 0.7635 | 69 / 69 |
| 030-153443.810-brewing-move.json | Full stack · original grid | 9 / 9 | 0.357 | 38.47 / 43.75 / 43.75 / 43.75 | 0.7635 / 0.7635 / 0.7635 | 69 / 69 |
| 030-153443.810-brewing-move.json | PRs 2–4 only | 9 / 9 | 0.346 | 38.33 / 39.58 / 39.58 / 39.58 | 0.7635 / 0.7635 / 0.7635 | 69 / 69 |
| 031-153448.198-opening_fridge-move.json | Original native | 9 / 8 | 2.484 | 242.38 / 419.10 / 419.10 / 419.10 | 6.8989 / 7.1921 / 7.1921 | 608.5 / 722 |
| 031-153448.198-opening_fridge-move.json | Full stack · fine grid | 9 / 9 | 1.864 | 191.94 / 372.64 / 372.64 / 372.64 | 3.1061 / 7.3591 / 7.3591 | 1354 / 3188 |
| 031-153448.198-opening_fridge-move.json | Full stack · original grid | 9 / 9 | 1.245 | 144.71 / 219.66 / 219.66 / 219.66 | 2.9487 / 3.4281 / 3.4281 | 528 / 624 |
| 031-153448.198-opening_fridge-move.json | PRs 2–4 only | 9 / 9 | 1.459 | 164.46 / 284.69 / 284.69 / 284.69 | 3.3930 / 7.0646 / 7.0646 | 614 / 707 |
| 032-153503.048-opening_fridge-move.json | Original native | 9 / 9 | 0.311 | 34.33 / 35.65 / 35.65 / 35.65 | 0.3425 / 0.3425 / 0.3425 | 46 / 46 |
| 032-153503.048-opening_fridge-move.json | Full stack · fine grid | 9 / 9 | 0.254 | 28.21 / 29.53 / 29.53 / 29.53 | 0.3425 / 0.3425 / 0.3425 | 46 / 46 |
| 032-153503.048-opening_fridge-move.json | Full stack · original grid | 9 / 9 | 0.247 | 27.26 / 28.57 / 28.57 / 28.57 | 0.3425 / 0.3425 / 0.3425 | 46 / 46 |
| 032-153503.048-opening_fridge-move.json | PRs 2–4 only | 9 / 9 | 0.249 | 27.43 / 29.49 / 29.49 / 29.49 | 0.3425 / 0.3425 / 0.3425 | 46 / 46 |
| 033-153505.743-opening_fridge-open_door.json | Original native | 9 / 9 | 1.636 | 181.89 / 184.25 / 184.25 / 184.25 | 0.6133 / 0.6576 / 0.6576 | 6 / 7 |
| 033-153505.743-opening_fridge-open_door.json | Full stack · fine grid | 9 / 9 | 1.005 | 112.08 / 113.66 / 113.66 / 113.66 | 0.6595 / 0.6750 / 0.6750 | 8 / 15 |
| 033-153505.743-opening_fridge-open_door.json | Full stack · original grid | 9 / 9 | 0.966 | 107.15 / 107.74 / 107.74 / 107.74 | 0.6133 / 0.6576 / 0.6576 | 6 / 9 |
| 033-153505.743-opening_fridge-open_door.json | PRs 2–4 only | 9 / 9 | 0.967 | 107.50 / 108.16 / 108.16 / 108.16 | 0.6133 / 0.6576 / 0.6576 | 6 / 9 |
| 034-153505.957-opening_fridge-open_door.json | Original native | 9 / 9 | 0.740 | 82.12 / 85.28 / 85.28 / 85.28 | 0.4535 / 0.4535 / 0.4535 | 2 / 2 |
| 034-153505.957-opening_fridge-open_door.json | Full stack · fine grid | 9 / 9 | 0.078 | 8.46 / 9.37 / 9.37 / 9.37 | 0.4535 / 0.4535 / 0.4535 | 2 / 2 |
| 034-153505.957-opening_fridge-open_door.json | Full stack · original grid | 9 / 9 | 0.076 | 8.45 / 9.35 / 9.35 / 9.35 | 0.4535 / 0.4535 / 0.4535 | 2 / 2 |
| 034-153505.957-opening_fridge-open_door.json | PRs 2–4 only | 9 / 9 | 0.081 | 8.90 / 10.03 / 10.03 / 10.03 | 0.4535 / 0.4535 / 0.4535 | 2 / 2 |
| 035-153506.192-opening_fridge-open_door.json | Original native | 9 / 9 | 0.707 | 78.66 / 80.40 / 80.40 / 80.40 | 0.4109 / 0.4109 / 0.4109 | 2 / 2 |
| 035-153506.192-opening_fridge-open_door.json | Full stack · fine grid | 9 / 9 | 0.062 | 6.76 / 7.34 / 7.34 / 7.34 | 0.4109 / 0.4109 / 0.4109 | 2 / 2 |
| 035-153506.192-opening_fridge-open_door.json | Full stack · original grid | 9 / 9 | 0.053 | 5.89 / 6.08 / 6.08 / 6.08 | 0.4109 / 0.4109 / 0.4109 | 2 / 2 |
| 035-153506.192-opening_fridge-open_door.json | PRs 2–4 only | 9 / 9 | 0.052 | 5.78 / 5.92 / 5.92 / 5.92 | 0.4109 / 0.4109 / 0.4109 | 2 / 2 |
| 036-153506.394-opening_fridge-open_door.json | Original native | 9 / 9 | 0.704 | 77.88 / 82.12 / 82.12 / 82.12 | 0.3246 / 0.3246 / 0.3246 | 2 / 2 |
| 036-153506.394-opening_fridge-open_door.json | Full stack · fine grid | 9 / 9 | 0.051 | 5.63 / 5.86 / 5.86 / 5.86 | 0.3246 / 0.3246 / 0.3246 | 2 / 2 |
| 036-153506.394-opening_fridge-open_door.json | Full stack · original grid | 9 / 9 | 0.042 | 4.65 / 4.98 / 4.98 / 4.98 | 0.3246 / 0.3246 / 0.3246 | 2 / 2 |
| 036-153506.394-opening_fridge-open_door.json | PRs 2–4 only | 9 / 9 | 0.041 | 4.52 / 4.73 / 4.73 / 4.73 | 0.3246 / 0.3246 / 0.3246 | 2 / 2 |
| 037-153506.598-opening_fridge-open_door.json | Original native | 9 / 9 | 0.701 | 77.68 / 79.86 / 79.86 / 79.86 | 0.2257 / 0.2257 / 0.2257 | 2 / 2 |
| 037-153506.598-opening_fridge-open_door.json | Full stack · fine grid | 9 / 9 | 0.051 | 5.63 / 6.01 / 6.01 / 6.01 | 0.2257 / 0.2257 / 0.2257 | 2 / 2 |
| 037-153506.598-opening_fridge-open_door.json | Full stack · original grid | 9 / 9 | 0.049 | 5.08 / 7.42 / 7.42 / 7.42 | 0.2257 / 0.2257 / 0.2257 | 2 / 2 |
| 037-153506.598-opening_fridge-open_door.json | PRs 2–4 only | 9 / 9 | 0.045 | 5.01 / 5.19 / 5.19 / 5.19 | 0.2257 / 0.2257 / 0.2257 | 2 / 2 |
| 038-153506.794-opening_fridge-open_door.json | Original native | 9 / 9 | 0.703 | 78.04 / 80.91 / 80.91 / 80.91 | 0.2438 / 0.2438 / 0.2438 | 2 / 2 |
| 038-153506.794-opening_fridge-open_door.json | Full stack · fine grid | 9 / 9 | 0.056 | 6.11 / 6.51 / 6.51 / 6.51 | 0.2438 / 0.2438 / 0.2438 | 2 / 2 |
| 038-153506.794-opening_fridge-open_door.json | Full stack · original grid | 9 / 9 | 0.051 | 5.70 / 6.13 / 6.13 / 6.13 | 0.2438 / 0.2438 / 0.2438 | 2 / 2 |
| 038-153506.794-opening_fridge-open_door.json | PRs 2–4 only | 9 / 9 | 0.050 | 5.50 / 5.79 / 5.79 / 5.79 | 0.2438 / 0.2438 / 0.2438 | 2 / 2 |
| 039-153507.001-opening_fridge-open_door.json | Original native | 9 / 9 | 0.697 | 76.79 / 81.09 / 81.09 / 81.09 | 0.3226 / 0.3226 / 0.3226 | 2 / 2 |
| 039-153507.001-opening_fridge-open_door.json | Full stack · fine grid | 9 / 9 | 0.056 | 6.15 / 6.81 / 6.81 / 6.81 | 0.3226 / 0.3226 / 0.3226 | 2 / 2 |
| 039-153507.001-opening_fridge-open_door.json | Full stack · original grid | 9 / 9 | 0.052 | 5.71 / 6.08 / 6.08 / 6.08 | 0.3226 / 0.3226 / 0.3226 | 2 / 2 |
| 039-153507.001-opening_fridge-open_door.json | PRs 2–4 only | 9 / 9 | 0.052 | 5.72 / 6.18 / 6.18 / 6.18 | 0.3226 / 0.3226 / 0.3226 | 2 / 2 |
| 040-153507.196-opening_fridge-open_door.json | Original native | 9 / 9 | 0.737 | 81.74 / 84.27 / 84.27 / 84.27 | 0.4046 / 0.4046 / 0.4046 | 2 / 2 |
| 040-153507.196-opening_fridge-open_door.json | Full stack · fine grid | 9 / 9 | 0.076 | 8.48 / 9.00 / 9.00 / 9.00 | 0.4046 / 0.4046 / 0.4046 | 2 / 2 |
| 040-153507.196-opening_fridge-open_door.json | Full stack · original grid | 9 / 9 | 0.078 | 8.65 / 9.25 / 9.25 / 9.25 | 0.4046 / 0.4046 / 0.4046 | 2 / 2 |
| 040-153507.196-opening_fridge-open_door.json | PRs 2–4 only | 9 / 9 | 0.078 | 8.63 / 9.05 / 9.05 / 9.05 | 0.4046 / 0.4046 / 0.4046 | 2 / 2 |
| 041-153514.279-opening_fridge-move.json | Original native | 9 / 9 | 0.334 | 37.03 / 38.78 / 38.78 / 38.78 | 0.4413 / 0.4413 / 0.4413 | 46 / 46 |
| 041-153514.279-opening_fridge-move.json | Full stack · fine grid | 9 / 9 | 0.265 | 29.19 / 30.93 / 30.93 / 30.93 | 0.4413 / 0.4413 / 0.4413 | 46 / 46 |
| 041-153514.279-opening_fridge-move.json | Full stack · original grid | 9 / 9 | 0.264 | 29.14 / 30.84 / 30.84 / 30.84 | 0.4413 / 0.4413 / 0.4413 | 46 / 46 |
| 041-153514.279-opening_fridge-move.json | PRs 2–4 only | 9 / 9 | 0.267 | 29.36 / 31.85 / 31.85 / 31.85 | 0.4413 / 0.4413 / 0.4413 | 46 / 46 |
| 042-153516.784-adding_milk-move.json | Original native | 9 / 9 | 0.651 | 67.50 / 104.52 / 104.52 / 104.52 | 0.7259 / 0.7522 / 0.7522 | 95 / 107 |
| 042-153516.784-adding_milk-move.json | Full stack · fine grid | 9 / 9 | 0.572 | 63.42 / 67.52 / 67.52 / 67.52 | 0.7138 / 0.7207 / 0.7207 | 291 / 304 |
| 042-153516.784-adding_milk-move.json | Full stack · original grid | 9 / 9 | 0.427 | 48.18 / 48.63 / 48.63 / 48.63 | 0.7282 / 0.7522 / 0.7522 | 95 / 107 |
| 042-153516.784-adding_milk-move.json | PRs 2–4 only | 9 / 9 | 0.436 | 48.11 / 51.58 / 51.58 / 51.58 | 0.7282 / 0.7522 / 0.7522 | 95 / 107 |
| 043-153521.282-adding_milk-move.json | Original native | 9 / 9 | 0.096 | 10.40 / 12.26 / 12.26 / 12.26 | 1.3431 / 1.3431 / 1.3431 | 2 / 2 |
| 043-153521.282-adding_milk-move.json | Full stack · fine grid | 9 / 9 | 0.092 | 10.04 / 12.06 / 12.06 / 12.06 | 1.3431 / 1.3431 / 1.3431 | 2 / 2 |
| 043-153521.282-adding_milk-move.json | Full stack · original grid | 9 / 9 | 0.080 | 9.03 / 9.73 / 9.73 / 9.73 | 1.3431 / 1.3431 / 1.3431 | 2 / 2 |
| 043-153521.282-adding_milk-move.json | PRs 2–4 only | 9 / 9 | 0.080 | 9.27 / 9.96 / 9.96 / 9.96 | 1.3431 / 1.3431 / 1.3431 | 2 / 2 |
| 044-153521.781-adding_milk-move.json | Original native | 9 / 9 | 0.616 | 67.93 / 73.30 / 73.30 / 73.30 | 1.0808 / 1.0808 / 1.0808 | 95 / 95 |
| 044-153521.781-adding_milk-move.json | Full stack · fine grid | 9 / 9 | 0.496 | 53.63 / 58.51 / 58.51 / 58.51 | 1.0808 / 1.0808 / 1.0808 | 95 / 95 |
| 044-153521.781-adding_milk-move.json | Full stack · original grid | 9 / 9 | 0.464 | 51.25 / 53.75 / 53.75 / 53.75 | 1.0808 / 1.0808 / 1.0808 | 95 / 95 |
| 044-153521.781-adding_milk-move.json | PRs 2–4 only | 9 / 9 | 0.468 | 52.14 / 53.81 / 53.81 / 53.81 | 1.0808 / 1.0808 / 1.0808 | 95 / 95 |
| 045-153529.604-adding_milk-move.json | Original native | 9 / 9 | 0.678 | 74.06 / 81.67 / 81.67 / 81.67 | 1.0806 / 1.0806 / 1.0806 | 95 / 95 |
| 045-153529.604-adding_milk-move.json | Full stack · fine grid | 9 / 9 | 0.541 | 58.56 / 68.67 / 68.67 / 68.67 | 1.0806 / 1.0806 / 1.0806 | 95 / 95 |
| 045-153529.604-adding_milk-move.json | Full stack · original grid | 9 / 9 | 0.504 | 55.37 / 59.03 / 59.03 / 59.03 | 1.0806 / 1.0806 / 1.0806 | 95 / 95 |
| 045-153529.604-adding_milk-move.json | PRs 2–4 only | 9 / 9 | 0.510 | 56.83 / 58.00 / 58.00 / 58.00 | 1.0806 / 1.0806 / 1.0806 | 95 / 95 |
| 046-153535.297-adding_milk-move.json | Original native | 9 / 9 | 2.433 | 244.13 / 443.10 / 443.10 / 443.10 | 4.4649 / 5.8272 / 5.8272 | 572 / 736 |
| 046-153535.297-adding_milk-move.json | Full stack · fine grid | 9 / 9 | 2.214 | 231.76 / 418.29 / 418.29 / 418.29 | 4.1769 / 4.6260 / 4.6260 | 1406 / 1856 |
| 046-153535.297-adding_milk-move.json | Full stack · original grid | 9 / 9 | 1.320 | 151.30 / 209.55 / 209.55 / 209.55 | 4.6258 / 5.8272 / 5.8272 | 588 / 736 |
| 046-153535.297-adding_milk-move.json | PRs 2–4 only | 9 / 9 | 1.320 | 151.39 / 209.17 / 209.17 / 209.17 | 4.6258 / 5.8272 / 5.8272 | 598 / 736 |
| 047-153542.859-adding_milk-pivot.json | Original native | 9 / 9 | 0.334 | 36.75 / 40.48 / 40.48 / 40.48 | 3.1415 / 3.1415 / 3.1415 | 37 / 37 |
| 047-153542.859-adding_milk-pivot.json | Full stack · fine grid | 9 / 9 | 0.398 | 40.22 / 70.44 / 70.44 / 70.44 | 3.1415 / 3.1415 / 3.1415 | 37 / 37 |
| 047-153542.859-adding_milk-pivot.json | Full stack · original grid | 9 / 9 | 0.265 | 29.44 / 31.01 / 31.01 / 31.01 | 3.1415 / 3.1415 / 3.1415 | 37 / 37 |
| 047-153542.859-adding_milk-pivot.json | PRs 2–4 only | 9 / 9 | 0.269 | 29.30 / 35.41 / 35.41 / 35.41 | 3.1415 / 3.1415 / 3.1415 | 37 / 37 |
| 048-153548.126-adding_milk-pivot.json | Original native | 9 / 9 | 0.331 | 36.94 / 38.88 / 38.88 / 38.88 | 3.1398 / 3.1398 / 3.1398 | 37 / 37 |
| 048-153548.126-adding_milk-pivot.json | Full stack · fine grid | 9 / 9 | 0.361 | 39.84 / 41.33 / 41.33 / 41.33 | 3.1398 / 3.1398 / 3.1398 | 37 / 37 |
| 048-153548.126-adding_milk-pivot.json | Full stack · original grid | 9 / 9 | 0.262 | 29.41 / 29.73 / 29.73 / 29.73 | 3.1398 / 3.1398 / 3.1398 | 37 / 37 |
| 048-153548.126-adding_milk-pivot.json | PRs 2–4 only | 9 / 9 | 0.264 | 29.14 / 31.14 / 31.14 / 31.14 | 3.1398 / 3.1398 / 3.1398 | 37 / 37 |
| 049-153550.610-adding_milk-move.json | Original native | 9 / 9 | 2.027 | 219.78 / 297.75 / 297.75 / 297.75 | 4.0581 / 4.5400 / 4.5400 | 578 / 651 |
| 049-153550.610-adding_milk-move.json | Full stack · fine grid | 9 / 9 | 1.666 | 163.19 / 255.47 / 255.47 / 255.47 | 4.3757 / 4.5399 / 4.5399 | 1596 / 1821 |
| 049-153550.610-adding_milk-move.json | Full stack · original grid | 9 / 9 | 1.586 | 174.61 / 285.24 / 285.24 / 285.24 | 4.5332 / 4.5565 / 4.5565 | 495 / 629 |
| 049-153550.610-adding_milk-move.json | PRs 2–4 only | 9 / 9 | 1.661 | 183.00 / 297.66 / 297.66 / 297.66 | 4.5332 / 5.2274 / 5.2274 | 585 / 861 |
| 050-153558.469-adding_milk-move.json | Original native | 9 / 9 | 0.674 | 74.46 / 77.25 / 77.25 / 77.25 | 1.0808 / 1.0808 / 1.0808 | 95 / 95 |
| 050-153558.469-adding_milk-move.json | Full stack · fine grid | 9 / 9 | 0.547 | 60.70 / 63.42 / 63.42 / 63.42 | 1.0808 / 1.0808 / 1.0808 | 95 / 95 |
| 050-153558.469-adding_milk-move.json | Full stack · original grid | 9 / 9 | 0.513 | 56.68 / 62.24 / 62.24 / 62.24 | 1.0808 / 1.0808 / 1.0808 | 95 / 95 |
| 050-153558.469-adding_milk-move.json | PRs 2–4 only | 9 / 9 | 0.515 | 57.42 / 59.20 / 59.20 / 59.20 | 1.0808 / 1.0808 / 1.0808 | 95 / 95 |
| 051-153603.443-adding_milk-move.json | Original native | 9 / 9 | 0.615 | 67.97 / 70.88 / 70.88 / 70.88 | 1.0806 / 1.0806 / 1.0806 | 95 / 95 |
| 051-153603.443-adding_milk-move.json | Full stack · fine grid | 9 / 9 | 0.493 | 54.40 / 57.20 / 57.20 / 57.20 | 1.0806 / 1.0806 / 1.0806 | 95 / 95 |
| 051-153603.443-adding_milk-move.json | Full stack · original grid | 9 / 9 | 0.468 | 51.70 / 54.47 / 54.47 / 54.47 | 1.0806 / 1.0806 / 1.0806 | 95 / 95 |
| 051-153603.443-adding_milk-move.json | PRs 2–4 only | 9 / 9 | 0.472 | 52.56 / 53.44 / 53.44 / 53.44 | 1.0806 / 1.0806 / 1.0806 | 95 / 95 |
| 052-153608.679-closing_fridge-move.json | Original native | 9 / 9 | 1.426 | 157.78 / 178.08 / 178.08 / 178.08 | 2.0207 / 2.2371 / 2.2371 | 261 / 299 |
| 052-153608.679-closing_fridge-move.json | Full stack · fine grid | 9 / 9 | 1.095 | 118.86 / 157.73 / 157.73 / 157.73 | 1.9797 / 2.1451 / 2.1451 | 716 / 827 |
| 052-153608.679-closing_fridge-move.json | Full stack · original grid | 9 / 9 | 0.791 | 87.33 / 93.48 / 93.48 / 93.48 | 2.1223 / 2.2785 / 2.2785 | 299 / 353 |
| 052-153608.679-closing_fridge-move.json | PRs 2–4 only | 9 / 9 | 0.803 | 88.59 / 96.41 / 96.41 / 96.41 | 2.1223 / 2.2785 / 2.2785 | 299 / 353 |
| 053-153613.932-closing_fridge-move.json | Original native | 9 / 9 | 0.317 | 35.26 / 36.99 / 36.99 / 36.99 | 0.4416 / 0.4416 / 0.4416 | 46 / 46 |
| 053-153613.932-closing_fridge-move.json | Full stack · fine grid | 9 / 9 | 0.257 | 28.67 / 29.37 / 29.37 / 29.37 | 0.4416 / 0.4416 / 0.4416 | 46 / 46 |
| 053-153613.932-closing_fridge-move.json | Full stack · original grid | 9 / 9 | 0.249 | 27.36 / 30.97 / 30.97 / 30.97 | 0.4416 / 0.4416 / 0.4416 | 46 / 46 |
| 053-153613.932-closing_fridge-move.json | PRs 2–4 only | 9 / 9 | 0.249 | 27.60 / 29.29 / 29.29 / 29.29 | 0.4416 / 0.4416 / 0.4416 | 46 / 46 |
| 054-153616.186-closing_fridge-close_door.json | Original native | 9 / 9 | 0.723 | 80.16 / 82.82 / 82.82 / 82.82 | 0.4045 / 0.4045 / 0.4045 | 2 / 2 |
| 054-153616.186-closing_fridge-close_door.json | Full stack · fine grid | 9 / 9 | 0.071 | 7.90 / 8.44 / 8.44 / 8.44 | 0.4045 / 0.4045 / 0.4045 | 2 / 2 |
| 054-153616.186-closing_fridge-close_door.json | Full stack · original grid | 9 / 9 | 0.067 | 7.41 / 7.82 / 7.82 / 7.82 | 0.4045 / 0.4045 / 0.4045 | 2 / 2 |
| 054-153616.186-closing_fridge-close_door.json | PRs 2–4 only | 9 / 9 | 0.067 | 7.54 / 7.75 / 7.75 / 7.75 | 0.4045 / 0.4045 / 0.4045 | 2 / 2 |
| 055-153616.384-closing_fridge-close_door.json | Original native | 9 / 9 | 0.727 | 80.52 / 82.81 / 82.81 / 82.81 | 0.3226 / 0.3226 / 0.3226 | 2 / 2 |
| 055-153616.384-closing_fridge-close_door.json | Full stack · fine grid | 9 / 9 | 0.076 | 8.40 / 8.94 / 8.94 / 8.94 | 0.3226 / 0.3226 / 0.3226 | 2 / 2 |
| 055-153616.384-closing_fridge-close_door.json | Full stack · original grid | 9 / 9 | 0.074 | 8.25 / 9.10 / 9.10 / 9.10 | 0.3226 / 0.3226 / 0.3226 | 2 / 2 |
| 055-153616.384-closing_fridge-close_door.json | PRs 2–4 only | 9 / 9 | 0.071 | 7.90 / 8.39 / 8.39 / 8.39 | 0.3226 / 0.3226 / 0.3226 | 2 / 2 |
| 056-153616.607-closing_fridge-close_door.json | Original native | 9 / 9 | 0.720 | 80.06 / 82.52 / 82.52 / 82.52 | 0.2438 / 0.2438 / 0.2438 | 2 / 2 |
| 056-153616.607-closing_fridge-close_door.json | Full stack · fine grid | 9 / 9 | 0.073 | 8.08 / 8.71 / 8.71 / 8.71 | 0.2438 / 0.2438 / 0.2438 | 2 / 2 |
| 056-153616.607-closing_fridge-close_door.json | Full stack · original grid | 9 / 9 | 0.067 | 7.30 / 8.66 / 8.66 / 8.66 | 0.2438 / 0.2438 / 0.2438 | 2 / 2 |
| 056-153616.607-closing_fridge-close_door.json | PRs 2–4 only | 9 / 9 | 0.065 | 7.24 / 7.73 / 7.73 / 7.73 | 0.2438 / 0.2438 / 0.2438 | 2 / 2 |
| 057-153616.814-closing_fridge-close_door.json | Original native | 9 / 9 | 0.720 | 80.19 / 81.93 / 81.93 / 81.93 | 0.2257 / 0.2257 / 0.2257 | 2 / 2 |
| 057-153616.814-closing_fridge-close_door.json | Full stack · fine grid | 9 / 9 | 0.071 | 7.87 / 8.30 / 8.30 / 8.30 | 0.2257 / 0.2257 / 0.2257 | 2 / 2 |
| 057-153616.814-closing_fridge-close_door.json | Full stack · original grid | 9 / 9 | 0.069 | 7.79 / 8.27 / 8.27 / 8.27 | 0.2257 / 0.2257 / 0.2257 | 2 / 2 |
| 057-153616.814-closing_fridge-close_door.json | PRs 2–4 only | 9 / 9 | 0.069 | 7.66 / 8.08 / 8.08 / 8.08 | 0.2257 / 0.2257 / 0.2257 | 2 / 2 |
| 058-153617.011-closing_fridge-close_door.json | Original native | 9 / 9 | 0.708 | 78.84 / 79.81 / 79.81 / 79.81 | 0.3246 / 0.3246 / 0.3246 | 2 / 2 |
| 058-153617.011-closing_fridge-close_door.json | Full stack · fine grid | 9 / 9 | 0.071 | 7.66 / 9.23 / 9.23 / 9.23 | 0.3246 / 0.3246 / 0.3246 | 2 / 2 |
| 058-153617.011-closing_fridge-close_door.json | Full stack · original grid | 9 / 9 | 0.068 | 7.77 / 7.83 / 7.83 / 7.83 | 0.3246 / 0.3246 / 0.3246 | 2 / 2 |
| 058-153617.011-closing_fridge-close_door.json | PRs 2–4 only | 9 / 9 | 0.071 | 7.98 / 8.54 / 8.54 / 8.54 | 0.3246 / 0.3246 / 0.3246 | 2 / 2 |
| 059-153617.272-closing_fridge-close_door.json | Original native | 9 / 9 | 0.722 | 80.37 / 81.95 / 81.95 / 81.95 | 0.4109 / 0.4109 / 0.4109 | 2 / 2 |
| 059-153617.272-closing_fridge-close_door.json | Full stack · fine grid | 9 / 9 | 0.076 | 8.39 / 9.73 / 9.73 / 9.73 | 0.4109 / 0.4109 / 0.4109 | 2 / 2 |
| 059-153617.272-closing_fridge-close_door.json | Full stack · original grid | 9 / 9 | 0.081 | 8.87 / 9.60 / 9.60 / 9.60 | 0.4109 / 0.4109 / 0.4109 | 2 / 2 |
| 059-153617.272-closing_fridge-close_door.json | PRs 2–4 only | 9 / 9 | 0.080 | 8.88 / 9.47 / 9.47 / 9.47 | 0.4109 / 0.4109 / 0.4109 | 2 / 2 |
| 060-153617.462-closing_fridge-close_door.json | Original native | 9 / 9 | 0.712 | 79.04 / 81.05 / 81.05 / 81.05 | 0.4535 / 0.4535 / 0.4535 | 2 / 2 |
| 060-153617.462-closing_fridge-close_door.json | Full stack · fine grid | 9 / 9 | 0.075 | 8.09 / 9.89 / 9.89 / 9.89 | 0.4535 / 0.4535 / 0.4535 | 2 / 2 |
| 060-153617.462-closing_fridge-close_door.json | Full stack · original grid | 9 / 9 | 0.069 | 7.44 / 9.93 / 9.93 / 9.93 | 0.4535 / 0.4535 / 0.4535 | 2 / 2 |
| 060-153617.462-closing_fridge-close_door.json | PRs 2–4 only | 9 / 9 | 0.067 | 7.31 / 8.62 / 8.62 / 8.62 | 0.4535 / 0.4535 / 0.4535 | 2 / 2 |
| 061-153617.649-closing_fridge-close_door.json | Original native | 9 / 9 | 0.713 | 79.18 / 81.64 / 81.64 / 81.64 | 0.5457 / 0.5457 / 0.5457 | 2 / 2 |
| 061-153617.649-closing_fridge-close_door.json | Full stack · fine grid | 9 / 9 | 0.080 | 8.73 / 10.68 / 10.68 / 10.68 | 0.5457 / 0.5457 / 0.5457 | 2 / 2 |
| 061-153617.649-closing_fridge-close_door.json | Full stack · original grid | 9 / 9 | 0.077 | 7.36 / 16.25 / 16.25 / 16.25 | 0.5457 / 0.5457 / 0.5457 | 2 / 2 |
| 061-153617.649-closing_fridge-close_door.json | PRs 2–4 only | 9 / 9 | 0.074 | 7.43 / 14.79 / 14.79 / 14.79 | 0.5457 / 0.5457 / 0.5457 | 2 / 2 |
| 062-153624.673-closing_fridge-move.json | Original native | 9 / 9 | 0.324 | 34.90 / 44.88 / 44.88 / 44.88 | 0.3425 / 0.3425 / 0.3425 | 46 / 46 |
| 062-153624.673-closing_fridge-move.json | Full stack · fine grid | 9 / 9 | 0.264 | 28.46 / 36.53 / 36.53 / 36.53 | 0.3425 / 0.3425 / 0.3425 | 46 / 46 |
| 062-153624.673-closing_fridge-move.json | Full stack · original grid | 9 / 9 | 0.249 | 27.66 / 28.58 / 28.58 / 28.58 | 0.3425 / 0.3425 / 0.3425 | 46 / 46 |
| 062-153624.673-closing_fridge-move.json | PRs 2–4 only | 9 / 9 | 0.261 | 27.80 / 38.78 / 38.78 / 38.78 | 0.3425 / 0.3425 / 0.3425 | 46 / 46 |
| 063-153628.349-serving-move.json | Original native | 9 / 9 | 3.808 | 490.12 / 636.36 / 636.36 / 636.36 | 4.0388 / 5.3421 / 5.3421 | 552 / 656 |
| 063-153628.349-serving-move.json | Full stack · fine grid | 9 / 9 | 2.611 | 129.28 / 668.54 / 668.54 / 668.54 | 4.0880 / 4.3925 / 4.3925 | 1563 / 1816 |
| 063-153628.349-serving-move.json | Full stack · original grid | 9 / 9 | 2.172 | 313.67 / 359.14 / 359.14 / 359.14 | 4.1682 / 4.1682 / 4.1682 | 609 / 609 |
| 063-153628.349-serving-move.json | PRs 2–4 only | 9 / 9 | 2.583 | 333.53 / 486.63 / 486.63 / 486.63 | 4.1682 / 5.3421 / 5.3421 | 578 / 653 |
| 064-153638.882-serving-move.json | Original native | 9 / 9 | 0.623 | 66.44 / 89.60 / 89.60 / 89.60 | 0.6053 / 0.6053 / 0.6053 | 92 / 92 |
| 064-153638.882-serving-move.json | Full stack · fine grid | 9 / 9 | 0.485 | 51.99 / 61.02 / 61.02 / 61.02 | 0.6053 / 0.6053 / 0.6053 | 92 / 92 |
| 064-153638.882-serving-move.json | Full stack · original grid | 9 / 9 | 0.454 | 50.42 / 51.49 / 51.49 / 51.49 | 0.6053 / 0.6053 / 0.6053 | 92 / 92 |
| 064-153638.882-serving-move.json | PRs 2–4 only | 9 / 9 | 0.481 | 51.06 / 68.62 / 68.62 / 68.62 | 0.6053 / 0.6053 / 0.6053 | 92 / 92 |
| 065-153642.369-serving-move.json | Original native | 9 / 9 | 0.655 | 73.06 / 76.52 / 76.52 / 76.52 | 0.6051 / 0.6051 / 0.6051 | 91 / 91 |
| 065-153642.369-serving-move.json | Full stack · fine grid | 9 / 9 | 0.501 | 55.31 / 57.23 / 57.23 / 57.23 | 0.6051 / 0.6051 / 0.6051 | 91 / 91 |
| 065-153642.369-serving-move.json | Full stack · original grid | 9 / 9 | 0.484 | 53.16 / 55.87 / 55.87 / 55.87 | 0.6051 / 0.6051 / 0.6051 | 91 / 91 |
| 065-153642.369-serving-move.json | PRs 2–4 only | 9 / 9 | 0.492 | 53.85 / 61.07 / 61.07 / 61.07 | 0.6051 / 0.6051 / 0.6051 | 91 / 91 |
| 066-153644.748-serving-carry.json | Original native | 9 / 9 | 0.247 | 27.26 / 31.60 / 31.60 / 31.60 | 1.8185 / 1.8253 / 1.8253 | 5 / 5 |
| 066-153644.748-serving-carry.json | Full stack · fine grid | 9 / 9 | 0.227 | 24.88 / 33.61 / 33.61 / 33.61 | 1.8253 / 1.8253 / 1.8253 | 5 / 5 |
| 066-153644.748-serving-carry.json | Full stack · original grid | 9 / 9 | 0.212 | 23.45 / 25.72 / 25.72 / 25.72 | 1.6958 / 1.8300 / 1.8300 | 5 / 5 |
| 066-153644.748-serving-carry.json | PRs 2–4 only | 9 / 9 | 0.211 | 23.09 / 29.86 / 29.86 / 29.86 | 1.6958 / 1.8300 / 1.8300 | 5 / 5 |
| 067-153647.477-serving-pivot.json | Original native | 9 / 9 | 0.295 | 32.71 / 33.87 / 33.87 / 33.87 | 3.4735 / 3.4735 / 3.4735 | 22 / 22 |
| 067-153647.477-serving-pivot.json | Full stack · fine grid | 9 / 9 | 0.312 | 34.49 / 35.70 / 35.70 / 35.70 | 3.4735 / 3.4735 / 3.4735 | 22 / 22 |
| 067-153647.477-serving-pivot.json | Full stack · original grid | 9 / 9 | 0.250 | 27.50 / 29.25 / 29.25 / 29.25 | 3.4735 / 3.4735 / 3.4735 | 22 / 22 |
| 067-153647.477-serving-pivot.json | PRs 2–4 only | 9 / 9 | 0.251 | 27.70 / 29.36 / 29.36 / 29.36 | 3.4735 / 3.4735 / 3.4735 | 22 / 22 |
| 068-153651.758-serving-pivot.json | Original native | 9 / 9 | 0.296 | 32.79 / 33.65 / 33.65 / 33.65 | 3.4698 / 3.4698 / 3.4698 | 22 / 22 |
| 068-153651.758-serving-pivot.json | Full stack · fine grid | 9 / 9 | 0.318 | 35.27 / 37.06 / 37.06 / 37.06 | 3.4698 / 3.4698 / 3.4698 | 22 / 22 |
| 068-153651.758-serving-pivot.json | Full stack · original grid | 9 / 9 | 0.258 | 28.40 / 31.52 / 31.52 / 31.52 | 3.4698 / 3.4698 / 3.4698 | 22 / 22 |
| 068-153651.758-serving-pivot.json | PRs 2–4 only | 9 / 9 | 0.256 | 28.52 / 29.61 / 29.61 / 29.61 | 3.4698 / 3.4698 / 3.4698 | 22 / 22 |
| 069-153653.214-serving-carry.json | Original native | 9 / 9 | 0.487 | 53.69 / 57.50 / 57.50 / 57.50 | 5.9608 / 6.1924 / 6.1924 | 7 / 7 |
| 069-153653.214-serving-carry.json | Full stack · fine grid | 9 / 9 | 0.564 | 60.98 / 69.15 / 69.15 / 69.15 | 6.0968 / 6.1924 / 6.1924 | 7 / 7 |
| 069-153653.214-serving-carry.json | Full stack · original grid | 9 / 9 | 0.406 | 44.78 / 54.69 / 54.69 / 54.69 | 6.0596 / 6.1259 / 6.1259 | 7 / 7 |
| 069-153653.214-serving-carry.json | PRs 2–4 only | 9 / 9 | 0.383 | 42.69 / 46.44 / 46.44 / 46.44 | 5.9399 / 6.1108 / 6.1108 | 7 / 7 |
| 070-153659.568-serving-move.json | Original native | 9 / 9 | 0.520 | 57.53 / 61.03 / 61.03 / 61.03 | 0.8009 / 0.8009 / 0.8009 | 72 / 72 |
| 070-153659.568-serving-move.json | Full stack · fine grid | 9 / 9 | 0.423 | 46.67 / 49.44 / 49.44 / 49.44 | 0.8009 / 0.8009 / 0.8009 | 72 / 72 |
| 070-153659.568-serving-move.json | Full stack · original grid | 9 / 9 | 0.404 | 45.12 / 46.01 / 46.01 / 46.01 | 0.8009 / 0.8009 / 0.8009 | 72 / 72 |
| 070-153659.568-serving-move.json | PRs 2–4 only | 9 / 9 | 0.415 | 44.75 / 52.05 / 52.05 / 52.05 | 0.8009 / 0.8009 / 0.8009 | 72 / 72 |
| 071-153703.339-serving-move.json | Original native | 9 / 9 | 0.478 | 52.90 / 56.85 / 56.85 / 56.85 | 0.8008 / 0.8008 / 0.8008 | 72 / 72 |
| 071-153703.339-serving-move.json | Full stack · fine grid | 9 / 9 | 0.379 | 42.20 / 43.23 / 43.23 / 43.23 | 0.8008 / 0.8008 / 0.8008 | 72 / 72 |
| 071-153703.339-serving-move.json | Full stack · original grid | 9 / 9 | 0.361 | 39.84 / 41.74 / 41.74 / 41.74 | 0.8008 / 0.8008 / 0.8008 | 72 / 72 |
| 071-153703.339-serving-move.json | PRs 2–4 only | 9 / 9 | 0.365 | 40.16 / 43.78 / 43.78 / 43.78 | 0.8008 / 0.8008 / 0.8008 | 72 / 72 |
| 072-153706.529-serving-move.json | Original native | 9 / 9 | 0.143 | 14.82 / 18.97 / 18.97 / 18.97 | 2.7309 / 2.7309 / 2.7309 | 2 / 2 |
| 072-153706.529-serving-move.json | Full stack · fine grid | 9 / 9 | 0.207 | 19.83 / 29.61 / 29.61 / 29.61 | 2.7309 / 2.7309 / 2.7309 | 2 / 2 |
| 072-153706.529-serving-move.json | Full stack · original grid | 9 / 9 | 0.105 | 12.07 / 12.44 / 12.44 / 12.44 | 2.7309 / 2.7309 / 2.7309 | 2 / 2 |
| 072-153706.529-serving-move.json | PRs 2–4 only | 9 / 9 | 0.102 | 10.62 / 15.90 / 15.90 / 15.90 | 2.7309 / 2.7309 / 2.7309 | 2 / 2 |
| 073-153712.604-serving-move.json | Original native | 9 / 9 | 0.488 | 51.21 / 79.06 / 79.06 / 79.06 | 0.7636 / 0.7636 / 0.7636 | 69 / 69 |
| 073-153712.604-serving-move.json | Full stack · fine grid | 9 / 9 | 0.377 | 41.65 / 43.96 / 43.96 / 43.96 | 0.7636 / 0.7636 / 0.7636 | 69 / 69 |
| 073-153712.604-serving-move.json | Full stack · original grid | 9 / 9 | 0.349 | 38.59 / 41.06 / 41.06 / 41.06 | 0.7636 / 0.7636 / 0.7636 | 69 / 69 |
| 073-153712.604-serving-move.json | PRs 2–4 only | 9 / 9 | 0.358 | 38.73 / 45.83 / 45.83 / 45.83 | 0.7636 / 0.7636 / 0.7636 | 69 / 69 |
| 074-153716.800-serving-move.json | Original native | 9 / 9 | 1.115 | 123.62 / 127.30 / 127.30 / 127.30 | 0.7635 / 0.7635 / 0.7635 | 69 / 69 |
| 074-153716.800-serving-move.json | Full stack · fine grid | 9 / 9 | 0.394 | 44.44 / 46.19 / 46.19 / 46.19 | 0.7635 / 0.7635 / 0.7635 | 69 / 69 |
| 074-153716.800-serving-move.json | Full stack · original grid | 9 / 9 | 0.368 | 40.65 / 44.02 / 44.02 / 44.02 | 0.7635 / 0.7635 / 0.7635 | 69 / 69 |
| 074-153716.800-serving-move.json | PRs 2–4 only | 9 / 9 | 0.377 | 41.18 / 49.12 / 49.12 / 49.12 | 0.7635 / 0.7635 / 0.7635 | 69 / 69 |
| 075-153733.046-serving-carry.json | Original native | 9 / 9 | 0.536 | 60.66 / 63.30 / 63.30 / 63.30 | 3.6201 / 3.6201 / 3.6201 | 7 / 7 |
| 075-153733.046-serving-carry.json | Full stack · fine grid | 9 / 9 | 0.750 | 77.64 / 108.80 / 108.80 / 108.80 | 3.6764 / 3.8484 / 3.8484 | 7 / 7 |
| 075-153733.046-serving-carry.json | Full stack · original grid | 9 / 9 | 0.632 | 62.25 / 94.69 / 94.69 / 94.69 | 3.6655 / 3.7644 / 3.7644 | 7 / 7 |
| 075-153733.046-serving-carry.json | PRs 2–4 only | 9 / 9 | 0.644 | 64.66 / 95.36 / 95.36 / 95.36 | 3.6792 / 3.8485 / 3.8485 | 7 / 7 |
| 076-153746.084-serving-move.json | Original native | 9 / 9 | 0.498 | 54.86 / 60.43 / 60.43 / 60.43 | 0.6677 / 0.6677 / 0.6677 | 72 / 72 |
| 076-153746.084-serving-move.json | Full stack · fine grid | 9 / 9 | 0.409 | 44.88 / 48.84 / 48.84 / 48.84 | 0.6677 / 0.6677 / 0.6677 | 72 / 72 |
| 076-153746.084-serving-move.json | Full stack · original grid | 9 / 9 | 0.382 | 42.28 / 44.02 / 44.02 / 44.02 | 0.6677 / 0.6677 / 0.6677 | 72 / 72 |
| 076-153746.084-serving-move.json | PRs 2–4 only | 9 / 9 | 0.385 | 42.76 / 44.66 / 44.66 / 44.66 | 0.6677 / 0.6677 / 0.6677 | 72 / 72 |
| 077-153749.872-serving-move.json | Original native | 9 / 9 | 0.456 | 50.97 / 51.54 / 51.54 / 51.54 | 0.6674 / 0.6674 / 0.6674 | 72 / 72 |
| 077-153749.872-serving-move.json | Full stack · fine grid | 9 / 9 | 0.364 | 40.69 / 42.21 / 42.21 / 42.21 | 0.6674 / 0.6674 / 0.6674 | 72 / 72 |
| 077-153749.872-serving-move.json | Full stack · original grid | 9 / 9 | 0.349 | 38.93 / 39.85 / 39.85 / 39.85 | 0.6674 / 0.6674 / 0.6674 | 72 / 72 |
| 077-153749.872-serving-move.json | PRs 2–4 only | 9 / 9 | 0.353 | 38.74 / 42.13 / 42.13 / 42.13 | 0.6674 / 0.6674 / 0.6674 | 72 / 72 |
| 078-153807.814-grabbing_filter-move.json | Original native | 9 / 9 | 4.223 | 468.11 / 494.72 / 494.72 / 494.72 | 6.1930 / 6.8201 / 6.8201 | 505 / 505 |
| 078-153807.814-grabbing_filter-move.json | Full stack · fine grid | 9 / 9 | 4.715 | 518.60 / 543.50 / 543.50 / 543.50 | 6.1930 / 6.1930 / 6.1930 | 2144 / 2144 |
| 078-153807.814-grabbing_filter-move.json | Full stack · original grid | 9 / 9 | 3.304 | 366.91 / 374.92 / 374.92 / 374.92 | 6.1930 / 6.1930 / 6.1930 | 505 / 505 |
| 078-153807.814-grabbing_filter-move.json | PRs 2–4 only | 9 / 9 | 3.304 | 369.46 / 376.10 / 376.10 / 376.10 | 6.1930 / 6.8201 / 6.8201 | 505 / 505 |
| 079-153815.779-grabbing_filter-move.json | Original native | 9 / 9 | 0.397 | 43.72 / 47.54 / 47.54 / 47.54 | 0.6309 / 0.6309 / 0.6309 | 62 / 62 |
| 079-153815.779-grabbing_filter-move.json | Full stack · fine grid | 9 / 9 | 0.325 | 36.44 / 37.84 / 37.84 / 37.84 | 0.6309 / 0.6309 / 0.6309 | 62 / 62 |
| 079-153815.779-grabbing_filter-move.json | Full stack · original grid | 9 / 9 | 0.308 | 34.12 / 35.70 / 35.70 / 35.70 | 0.6309 / 0.6309 / 0.6309 | 62 / 62 |
| 079-153815.779-grabbing_filter-move.json | PRs 2–4 only | 9 / 9 | 0.315 | 34.10 / 37.84 / 37.84 / 37.84 | 0.6309 / 0.6309 / 0.6309 | 62 / 62 |
| 080-153819.582-unlocking_portafilter-pivot.json | Original native | 9 / 9 | 0.981 | 109.05 / 111.69 / 111.69 / 111.69 | 2.4096 / 2.4096 / 2.4096 | 34 / 34 |
| 080-153819.582-unlocking_portafilter-pivot.json | Full stack · fine grid | 9 / 9 | 0.352 | 38.84 / 40.41 / 40.41 / 40.41 | 2.4096 / 2.4096 / 2.4096 | 34 / 34 |
| 080-153819.582-unlocking_portafilter-pivot.json | Full stack · original grid | 9 / 9 | 0.277 | 30.49 / 32.35 / 32.35 / 32.35 | 2.4096 / 2.4096 / 2.4096 | 34 / 34 |
| 080-153819.582-unlocking_portafilter-pivot.json | PRs 2–4 only | 9 / 9 | 0.276 | 30.73 / 31.54 / 31.54 / 31.54 | 2.4096 / 2.4096 / 2.4096 | 34 / 34 |
| 081-153824.736-unlocking_portafilter-move.json | Original native | 9 / 9 | 0.091 | 10.04 / 10.64 / 10.64 / 10.64 | 0.0435 / 0.0435 / 0.0435 | 8 / 8 |
| 081-153824.736-unlocking_portafilter-move.json | Full stack · fine grid | 9 / 9 | 0.080 | 8.94 / 9.16 / 9.16 / 9.16 | 0.0435 / 0.0435 / 0.0435 | 8 / 8 |
| 081-153824.736-unlocking_portafilter-move.json | Full stack · original grid | 9 / 9 | 0.078 | 8.73 / 9.04 / 9.04 / 9.04 | 0.0435 / 0.0435 / 0.0435 | 8 / 8 |
| 081-153824.736-unlocking_portafilter-move.json | PRs 2–4 only | 9 / 9 | 0.078 | 8.66 / 8.97 / 8.97 / 8.97 | 0.0435 / 0.0435 / 0.0435 | 8 / 8 |
| 082-153825.312-unlocking_portafilter-circular.json | Original native | 9 / 9 | 0.097 | 10.86 / 11.20 / 11.20 / 11.20 | 0.0276 / 0.0276 / 0.0276 | 9 / 9 |
| 082-153825.312-unlocking_portafilter-circular.json | Full stack · fine grid | 9 / 9 | 0.084 | 9.27 / 10.03 / 10.03 / 10.03 | 0.0276 / 0.0276 / 0.0276 | 9 / 9 |
| 082-153825.312-unlocking_portafilter-circular.json | Full stack · original grid | 9 / 9 | 0.082 | 9.05 / 9.37 / 9.37 / 9.37 | 0.0276 / 0.0276 / 0.0276 | 9 / 9 |
| 082-153825.312-unlocking_portafilter-circular.json | PRs 2–4 only | 9 / 9 | 0.082 | 9.01 / 9.37 / 9.37 / 9.37 | 0.0276 / 0.0276 / 0.0276 | 9 / 9 |
| 083-153826.144-unlocking_portafilter-move.json | Original native | 9 / 9 | 0.092 | 10.19 / 11.48 / 11.48 / 11.48 | 0.0442 / 0.0442 / 0.0442 | 8 / 8 |
| 083-153826.144-unlocking_portafilter-move.json | Full stack · fine grid | 9 / 9 | 0.080 | 9.01 / 9.20 / 9.20 / 9.20 | 0.0442 / 0.0442 / 0.0442 | 8 / 8 |
| 083-153826.144-unlocking_portafilter-move.json | Full stack · original grid | 9 / 9 | 0.078 | 8.50 / 9.06 / 9.06 / 9.06 | 0.0442 / 0.0442 / 0.0442 | 8 / 8 |
| 083-153826.144-unlocking_portafilter-move.json | PRs 2–4 only | 9 / 9 | 0.078 | 8.55 / 9.31 / 9.31 / 9.31 | 0.0442 / 0.0442 / 0.0442 | 8 / 8 |
| 084-153826.731-unlocking_portafilter-move.json | Original native | 9 / 9 | 0.093 | 10.29 / 10.84 / 10.84 / 10.84 | 0.1557 / 0.1557 / 0.1557 | 8 / 8 |
| 084-153826.731-unlocking_portafilter-move.json | Full stack · fine grid | 9 / 9 | 0.091 | 9.50 / 14.18 / 14.18 / 14.18 | 0.1557 / 0.1557 / 0.1557 | 8 / 8 |
| 084-153826.731-unlocking_portafilter-move.json | Full stack · original grid | 9 / 9 | 0.080 | 8.79 / 9.28 / 9.28 / 9.28 | 0.1557 / 0.1557 / 0.1557 | 8 / 8 |
| 084-153826.731-unlocking_portafilter-move.json | PRs 2–4 only | 9 / 9 | 0.079 | 8.73 / 9.15 / 9.15 / 9.15 | 0.1557 / 0.1557 / 0.1557 | 8 / 8 |
| 085-153827.773-unlocking_portafilter-circular.json | Original native | 9 / 9 | 0.096 | 10.56 / 11.27 / 11.27 / 11.27 | 0.0276 / 0.0276 / 0.0276 | 9 / 9 |
| 085-153827.773-unlocking_portafilter-circular.json | Full stack · fine grid | 9 / 9 | 0.084 | 9.21 / 9.79 / 9.79 / 9.79 | 0.0276 / 0.0276 / 0.0276 | 9 / 9 |
| 085-153827.773-unlocking_portafilter-circular.json | Full stack · original grid | 9 / 9 | 0.083 | 9.20 / 9.47 / 9.47 / 9.47 | 0.0276 / 0.0276 / 0.0276 | 9 / 9 |
| 085-153827.773-unlocking_portafilter-circular.json | PRs 2–4 only | 9 / 9 | 0.084 | 9.33 / 9.94 / 9.94 / 9.94 | 0.0276 / 0.0276 / 0.0276 | 9 / 9 |
| 086-153828.689-unlocking_portafilter-move.json | Original native | 9 / 9 | 0.154 | 16.90 / 18.10 / 18.10 / 18.10 | 0.2375 / 0.2375 / 0.2375 | 16 / 16 |
| 086-153828.689-unlocking_portafilter-move.json | Full stack · fine grid | 9 / 9 | 0.135 | 14.90 / 15.64 / 15.64 / 15.64 | 0.2375 / 0.2375 / 0.2375 | 16 / 16 |
| 086-153828.689-unlocking_portafilter-move.json | Full stack · original grid | 9 / 9 | 0.128 | 14.16 / 16.30 / 16.30 / 16.30 | 0.2375 / 0.2375 / 0.2375 | 16 / 16 |
| 086-153828.689-unlocking_portafilter-move.json | PRs 2–4 only | 9 / 9 | 0.127 | 13.83 / 15.83 / 15.83 / 15.83 | 0.2375 / 0.2375 / 0.2375 | 16 / 16 |
| 087-153833.922-cleaning-move.json | Original native | 9 / 9 | 0.617 | 67.45 / 76.25 / 76.25 / 76.25 | 1.1781 / 1.1781 / 1.1781 | 159 / 159 |
| 087-153833.922-cleaning-move.json | Full stack · fine grid | 9 / 9 | 0.568 | 62.44 / 69.64 / 69.64 / 69.64 | 1.1841 / 1.1841 / 1.1841 | 465 / 465 |
| 087-153833.922-cleaning-move.json | Full stack · original grid | 9 / 9 | 0.455 | 50.40 / 55.11 / 55.11 / 55.11 | 1.1781 / 1.1813 / 1.1813 | 159 / 159 |
| 087-153833.922-cleaning-move.json | PRs 2–4 only | 9 / 9 | 0.461 | 50.30 / 55.84 / 55.84 / 55.84 | 1.1781 / 1.1781 / 1.1781 | 159 / 159 |
| 088-153836.752-cleaning-move.json | Original native | 9 / 9 | 0.143 | 15.99 / 16.77 / 16.77 / 16.77 | 3.1101 / 3.1101 / 3.1101 | 2 / 2 |
| 088-153836.752-cleaning-move.json | Full stack · fine grid | 9 / 9 | 0.207 | 22.67 / 24.83 / 24.83 / 24.83 | 3.1101 / 3.1101 / 3.1101 | 2 / 2 |
| 088-153836.752-cleaning-move.json | Full stack · original grid | 9 / 9 | 0.121 | 12.20 / 19.16 / 19.16 / 19.16 | 3.1101 / 3.1101 / 3.1101 | 2 / 2 |
| 088-153836.752-cleaning-move.json | PRs 2–4 only | 9 / 9 | 0.112 | 11.56 / 18.99 / 18.99 / 18.99 | 3.1101 / 3.1101 / 3.1101 | 2 / 2 |
| 089-153842.065-cleaning-move.json | Original native | 9 / 9 | 0.239 | 26.26 / 28.08 / 28.08 / 28.08 | 0.2109 / 0.2109 / 0.2109 | 34 / 34 |
| 089-153842.065-cleaning-move.json | Full stack · fine grid | 9 / 9 | 0.195 | 21.13 / 24.40 / 24.40 / 24.40 | 0.2109 / 0.2109 / 0.2109 | 34 / 34 |
| 089-153842.065-cleaning-move.json | Full stack · original grid | 9 / 9 | 0.188 | 20.87 / 21.16 / 21.16 / 21.16 | 0.2109 / 0.2109 / 0.2109 | 34 / 34 |
| 089-153842.065-cleaning-move.json | PRs 2–4 only | 9 / 9 | 0.189 | 20.78 / 22.10 / 22.10 / 22.10 | 0.2109 / 0.2109 / 0.2109 | 34 / 34 |
| 090-153843.194-cleaning-circular.json | Original native | 9 / 9 | 0.086 | 9.31 / 10.69 / 10.69 / 10.69 | 0.1016 / 0.1017 / 0.1017 | 9 / 9 |
| 090-153843.194-cleaning-circular.json | Full stack · fine grid | 9 / 9 | 0.072 | 7.94 / 8.46 / 8.46 / 8.46 | 0.1016 / 0.1017 / 0.1017 | 9 / 9 |
| 090-153843.194-cleaning-circular.json | Full stack · original grid | 9 / 9 | 0.069 | 7.58 / 8.05 / 8.05 / 8.05 | 0.1017 / 0.1017 / 0.1017 | 9 / 9 |
| 090-153843.194-cleaning-circular.json | PRs 2–4 only | 9 / 9 | 0.070 | 7.67 / 8.23 / 8.23 / 8.23 | 0.1016 / 0.1017 / 0.1017 | 9 / 9 |
| 091-153846.164-cleaning-move.json | Original native | 9 / 9 | 0.243 | 26.88 / 29.03 / 29.03 / 29.03 | 0.2129 / 0.2129 / 0.2129 | 34 / 34 |
| 091-153846.164-cleaning-move.json | Full stack · fine grid | 9 / 9 | 0.197 | 21.66 / 23.77 / 23.77 / 23.77 | 0.2129 / 0.2129 / 0.2129 | 34 / 34 |
| 091-153846.164-cleaning-move.json | Full stack · original grid | 9 / 9 | 0.191 | 21.06 / 22.31 / 22.31 / 22.31 | 0.2129 / 0.2129 / 0.2129 | 34 / 34 |
| 091-153846.164-cleaning-move.json | PRs 2–4 only | 9 / 9 | 0.191 | 21.07 / 22.04 / 22.04 / 22.04 | 0.2129 / 0.2129 / 0.2129 | 34 / 34 |
| 092-153847.685-cleaning-move.json | Original native | 9 / 9 | 0.281 | 30.81 / 33.60 / 33.60 / 33.60 | 0.4668 / 0.4668 / 0.4668 | 42 / 42 |
| 092-153847.685-cleaning-move.json | Full stack · fine grid | 9 / 9 | 0.232 | 25.67 / 27.15 / 27.15 / 27.15 | 0.4668 / 0.4668 / 0.4668 | 42 / 42 |
| 092-153847.685-cleaning-move.json | Full stack · original grid | 9 / 9 | 0.221 | 24.45 / 26.25 / 26.25 / 26.25 | 0.4668 / 0.4668 / 0.4668 | 42 / 42 |
| 092-153847.685-cleaning-move.json | PRs 2–4 only | 9 / 9 | 0.225 | 24.86 / 25.85 / 25.85 / 25.85 | 0.4668 / 0.4668 / 0.4668 | 42 / 42 |
| 093-153849.454-cleaning-move.json | Original native | 9 / 9 | 0.238 | 26.52 / 27.45 / 27.45 / 27.45 | 0.2028 / 0.2028 / 0.2028 | 34 / 34 |
| 093-153849.454-cleaning-move.json | Full stack · fine grid | 9 / 9 | 0.193 | 21.39 / 22.46 / 22.46 / 22.46 | 0.2028 / 0.2028 / 0.2028 | 34 / 34 |
| 093-153849.454-cleaning-move.json | Full stack · original grid | 9 / 9 | 0.186 | 20.79 / 21.21 / 21.21 / 21.21 | 0.2028 / 0.2028 / 0.2028 | 34 / 34 |
| 093-153849.454-cleaning-move.json | PRs 2–4 only | 9 / 9 | 0.191 | 20.62 / 23.37 / 23.37 / 23.37 | 0.2028 / 0.2028 / 0.2028 | 34 / 34 |
| 094-153850.584-cleaning-circular.json | Original native | 9 / 9 | 0.080 | 8.54 / 10.84 / 10.84 / 10.84 | 0.1217 / 0.1217 / 0.1217 | 9 / 9 |
| 094-153850.584-cleaning-circular.json | Full stack · fine grid | 9 / 9 | 0.073 | 7.99 / 10.04 / 10.04 / 10.04 | 0.1217 / 0.1217 / 0.1217 | 9 / 9 |
| 094-153850.584-cleaning-circular.json | Full stack · original grid | 9 / 9 | 0.066 | 7.28 / 7.61 / 7.61 / 7.61 | 0.1217 / 0.1217 / 0.1217 | 9 / 9 |
| 094-153850.584-cleaning-circular.json | PRs 2–4 only | 9 / 9 | 0.066 | 7.24 / 7.62 / 7.62 / 7.62 | 0.1217 / 0.1217 / 0.1217 | 9 / 9 |
| 095-153853.741-cleaning-move.json | Original native | 9 / 9 | 0.237 | 26.27 / 27.38 / 27.38 / 27.38 | 0.2046 / 0.2046 / 0.2046 | 34 / 34 |
| 095-153853.741-cleaning-move.json | Full stack · fine grid | 9 / 9 | 0.199 | 21.40 / 24.48 / 24.48 / 24.48 | 0.2046 / 0.2046 / 0.2046 | 34 / 34 |
| 095-153853.741-cleaning-move.json | Full stack · original grid | 9 / 9 | 0.187 | 20.70 / 21.67 / 21.67 / 21.67 | 0.2046 / 0.2046 / 0.2046 | 34 / 34 |
| 095-153853.741-cleaning-move.json | PRs 2–4 only | 9 / 9 | 0.188 | 20.99 / 21.33 / 21.33 / 21.33 | 0.2046 / 0.2046 / 0.2046 | 34 / 34 |
| 096-153854.899-cleaning-move.json | Original native | 9 / 9 | 0.158 | 17.61 / 18.14 / 18.14 / 18.14 | 3.5455 / 3.5455 / 3.5455 | 2 / 2 |
| 096-153854.899-cleaning-move.json | Full stack · fine grid | 9 / 9 | 0.232 | 25.88 / 27.19 / 27.19 / 27.19 | 3.5455 / 3.5455 / 3.5455 | 2 / 2 |
| 096-153854.899-cleaning-move.json | Full stack · original grid | 9 / 9 | 0.118 | 13.10 / 14.24 / 14.24 / 14.24 | 3.5455 / 3.5455 / 3.5455 | 2 / 2 |
| 096-153854.899-cleaning-move.json | PRs 2–4 only | 9 / 9 | 0.117 | 12.52 / 16.19 / 16.19 / 16.19 | 3.5455 / 3.5455 / 3.5455 | 2 / 2 |
| 097-153900.281-finishing_up-move.json | Original native | 9 / 9 | 0.113 | 12.34 / 13.45 / 13.45 / 13.45 | 1.4596 / 1.4596 / 1.4596 | 2 / 2 |
| 097-153900.281-finishing_up-move.json | Full stack · fine grid | 9 / 9 | 0.105 | 11.42 / 12.77 / 12.77 / 12.77 | 1.4596 / 1.4596 / 1.4596 | 2 / 2 |
| 097-153900.281-finishing_up-move.json | Full stack · original grid | 9 / 9 | 0.083 | 9.16 / 9.93 / 9.93 / 9.93 | 1.4596 / 1.4596 / 1.4596 | 2 / 2 |
| 097-153900.281-finishing_up-move.json | PRs 2–4 only | 9 / 9 | 0.083 | 9.31 / 9.91 / 9.91 / 9.91 | 1.4596 / 1.4596 / 1.4596 | 2 / 2 |

## Duplicate scene 08 (retained, zero headline weight)

| Variant / process | N | Accepted | Total s | Median ms | p95 ms | p99 ms | Max ms | ≥1/5/30 s |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| Original native | 30 | 30 (100.00%) | 0.243 | 8.08 | 9.30 | 9.60 | 9.60 | 0/0/0 |
| Full stack · fine grid | 30 | 30 (100.00%) | 0.262 | 8.69 | 10.47 | 10.65 | 10.65 | 0/0/0 |
| Full stack · original grid | 30 | 30 (100.00%) | 0.242 | 7.91 | 9.45 | 9.81 | 9.81 | 0/0/0 |
| PRs 2–4 only | 30 | 30 (100.00%) | 0.248 | 7.86 | 9.78 | 11.45 | 11.45 | 0/0/0 |

## Limits and reproduction

The experiment has 3 independent cache initializations per variant, not thousands of independent requests. Intervals exhaust all 3^3 paired process-block bootstrap resamples and are exploratory. There is no precise p99 or formal significance claim. Earlier completed comparisons are not retroactively significant.

Calibration policy is declared above; modern-head measured/effective settings are recorded outside the timer. Default worker counts are unchanged. Missing original diagnostics are unavailable, never zero. Phase CPU/duration totals may overlap and must not be added as disjoint wall time.

Path lengths are compared within each input, conditional on audit acceptance; they are joint-space lengths, not execution times. Scene-cache fields are maxima of global retained-memory snapshots, not sums or request-local allocations. Warmup planning and audit costs are recorded separately in process accounting; startup calibration, request decoding and snapshot preparation remain outside measured replay.

Power, thermal and load observations (including unavailable readings), exact commands, outcomes, audit costs, per-process calibration, overlapping phase totals, counters and outlier identities are retained in the local JSON report. Dense audits provide evidence at their resolution, not proof of continuous collision freedom. No hardware moved.

Raw trajectories and journals remain local and checksum-indexed. Independent frozen offline audits completed and reconciled all 6,288 attempts; no additional rejection was found. Exact commands, scripts, and frozen identities remain in the local experiment. A commit of these readable reports does not back up those payloads.
