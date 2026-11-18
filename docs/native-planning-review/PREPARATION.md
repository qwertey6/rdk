The four production branches are extracted, each with one commit and its own tests. Published draft URLs are tracked in [PULL_REQUESTS.md](PULL_REQUESTS.md). Source IDs and line counts are in [evidence/diff-stats.json](evidence/diff-stats.json).

Open four dependent reviews in this order:

| PR | Head | Base | Description |
|---|---|---|---|
| 1 | `jake/native-1-correctness` | `main` | [Correctness and worker lifetime](01-correctness-and-worker-lifetime.md) |
| 2 | `jake/native-2-kinematics` | `jake/native-1-correctness` | [Native FK and IK](02-native-fk-and-ik.md) |
| 3 | `jake/native-3-collision` | `jake/native-2-kinematics` | [Collision checking and scene reuse](03-collision-and-scene-reuse.md) |
| 4 | `jake/native-4-finishing` | `jake/native-3-collision` | [Finishing and smoothing](04-finishing-and-smoothing.md) |

Each PR targets its predecessor so GitHub shows only that change. Keep this order through review. After a lower PR merges, update the next PR's base and rebase/cherry-pick onto the merged commit if the merge method changed its identity. Run affected checks after conflict resolution. Do not merge or cherry-pick the auxiliary cleanup/evidence branch into main.

The description's first line is the PR title. The rest is ready review text with measurements and tradeoffs. Images are embedded by relative links here; a GitHub PR body needs hosted absolute image URLs. Publish the evidence directory on a separate review-assets branch or upload the rendered PNGs to GitHub, then replace image links in the PR body. Do not add image/report commits to the production stack. Retain PDF figures and analysis code with the auxiliary evidence for reproduction.

Use the corrected-boundary cohort for PR1–3 increments and the six-block cohort for PR4 and the full-stack result. The first cohort's original PR1/2 boundaries were superseded after audit failures. Never relabel those old measurements as final PR1/2, sum incremental percentages across cohorts, or present the prior 23.1% result as the new extracted-stack measurement.

The full stack reduces ordered planning time 21.7% versus original native but worsens ordered p99 26.0%; the individual-scene total change is uncertain. PR1 can be slower because it restores necessary checking and joins workers. Review the 512 MiB scene-cache default and the scene04 path-quality flag before deployment. These are warm planning replay results, not hardware execution or request-to-first-motion measurements.

Raw trajectories, binaries and journals are local, checksum-indexed. Push/commit of the evidence directory does not back them up. Preserve those payloads separately if remote reproducibility or long-term retention is needed.
