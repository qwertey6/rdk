Four draft PRs are published in `qwertey6/rdk`. The publishing account has read-only access to `viamrobotics/rdk`; no upstream branches were changed. The fork’s main was fast-forwarded to the validated upstream base `8ab5892be050dadb5de3b140b6ee9ea212984f51`.

| Review | Draft | Base | Validated head |
|---|---|---|---|
| 1 | [[1/4] motionplan: preserve collision validation and join canceled native workers](https://github.com/qwertey6/rdk/pull/1) | `main` | `a2bc4d795553fdd4910dde49e992e79a24fccd9e` |
| 2 | [[2/4] motionplan: prepare native FK and batch IK objective evaluations](https://github.com/qwertey6/rdk/pull/2) | `jake/native-1-correctness` | `6a1590932f85add9d1ba3e5f231ead8fbe2de82c` |
| 3 | [[3/4] motionplan: stream conservative collision checks and reuse scene preparation](https://github.com/qwertey6/rdk/pull/3) | `jake/native-2-kinematics` | `899cedca5585f59b8d18b7ca2fc9f3941fb472e3` |
| 4 | [[4/4] motionplan: bound optional finishing and reuse exact collision states](https://github.com/qwertey6/rdk/pull/4) | `jake/native-3-collision` | `76654f793a7559db83e22d81e46fee807bc84439` |

Each PR contains one commit and only its own code/test diff. The summaries retain the original-native comparison, per-PR measurements, uncertainty, acceptance and quality/tail tradeoffs. All figures are hosted on this separate review-assets branch; no report or figure lands in the production stack. No reviewers were assigned and no PR was marked ready or merged.

To land this as a stack in upstream, an upstream writer must publish the four base/head branches there and open or recreate the corresponding upstream reviews. These fork drafts do not provide an automatic merge path into viamrobotics/rdk.
