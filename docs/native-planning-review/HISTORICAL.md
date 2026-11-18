Earlier cohorts are separate experiments and are not pooled with the extracted-stack results.

| Cohort | Original scene total | Candidate scene total | Scene change | Original ordered total | Candidate ordered total | Ordered change |
|---|---:|---:|---:|---:|---:|---:|
| November 13 cleanup recheck, four blocks | 250.424 s | 358.925 s | **43.3% more** | 109.468 s | 126.247 s | **15.3% more** |
| November 17 cleanup comparison, six blocks | 324.881 s | 316.164 s | **2.7% less**, uncertain | 152.348 s | 117.119 s | **23.1% less** |

The unfavorable November 13 cohort had competing builds and unusually variable startup calibration. It remains evidence; the later favorable result does not erase it. The November 17 ordered reduction's exploratory interval was 19.0–27.2%; its scene interval crossed equality. Those binaries still included diagnostics/replay support and are not the four extracted PR heads. The earlier 23.1% number must not replace the extracted stack's 21.7% result. No 30-second stall was reproduced. Full historical local records are on auxiliary commit 445ffb884; raw archives are not backed up by this review-assets branch.
