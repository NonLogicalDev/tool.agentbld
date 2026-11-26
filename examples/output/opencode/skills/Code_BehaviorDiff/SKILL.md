---
name: Code_BehaviorDiff
description: Explain before/after behavior in English pseudocode. Use for behavioral diffs, request or job flow comparisons, state transitions, and PR behavior breakdowns.
---

# Code Behavior Diff

Explain actors, conditions, actions, and outcomes in English. Source references carry implementation names; retain observable interfaces, statuses, and literal values when needed.

## Compare

- Identify versions or modes. For PRs, use the intended base and reviewed head; exclude inherited stack changes.
- Inspect both sides and needed callers. State missing evidence; distinguish source behavior from execution, tests, and deployment. Label inference, illustrations, and proposals.
- Explain only unless changes, experiments, or publication were requested.

## Explain

- Briefly state the outcome and documented or inferred motivation, then show aligned English steps: `-` old, `+` new, unmarked shared context. Use a standard Markdown code fence labeled `diff` for every behavioral diff.
- Follow the affected request, job, or state transition. Keep only context needed for activation, ordering, and results. Keep each side's order; shared steps cannot move across changed work. Group equivalent paths; separate distinct flows. Respect requested formats and narrow questions.
- Preserve consequential conditions: authorization, cache hits, missing state, retries, errors, and publication. Name responsible actors; distinguish local work from network requests. Count requests only with source support.
- Show reused, recomputed, and repeated work. Reuse does not imply skipped authorization, freshness checks, or client requests. For moves or representation changes, state preserved behavior; invent none.
- Finish with material consequences, costs, exceptions, and evidence limits. For reviews, separate optional improvements from actionable defects.

## Cite

Attach `# [O1]` to old steps and `# [N1]` to new steps. Below the block, list inspected sources with clickable revision-pinned file/line links and what each proves. Keep prose and source lists in standard Markdown. On supporting Flow surfaces, source lists may be collapsed in `<details>` with a `Sources` summary. Keep old/new locations separate. Group adjacent supported steps under one marker. Without URLs, use supplied revision/file/line labels; invent no permalink.

## Illustrative shape

```diff
 A worker combines selected inputs and uploads the result.
 If publication conflicts:
-    Discard the completed work and restart the merge.              # [O1]
+    Reload the current manifest.
+    If all inputs used for the completed merge remain present:
+        Retry publication using the result already uploaded.      # [N1]
+    Otherwise:
+        Discard the result and restart the merge.                  # [N2]
```

Replace this illustration with inspected behavior and sources.

## Findings for team reviews

Use the [team review checklist](resources/engineering-practices/team-review.md) to verify the review scope.
For each actionable finding, report the affected file and line, the conditions that trigger the failure, and its impact.
Tell reviewers which behaviors you inspected in source and which you exercised in tests.
