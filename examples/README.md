# Layered configuration example

Start with shared instructions and skills from [NonLogicalDev/conf.agents][conf-agents]. Personal notes and team
guidance adapt the same skills to home and work.

## What to explore

- [00-shared](input-sections/00-shared/) provides writing guidance, engineering practices, and reusable skills.
- [10-home](input-sections/10-home/) adds personal decision notes and guidance for resuming projects.
- [20-work](input-sections/20-work/) adds team review, verification, and delivery guidance.

For a personal setup, use the shared and home sections. Add the work section to apply team guidance.

## Render

Run from the repository root with `agentbld` installed.

```sh
agentbld render --sections-dir examples/input-sections --render-dir .tmp/example-work \
  --mode='codex:{"agents":true,"skills":true,"plugin":"full","plugin_name":"example-skills"}' \
  --mode=claude --mode=opencode
```

Inspect `.tmp/example-work/latest/` or browse the [saved work output](output/).

## Reference and skill examples

- **Extend a shared skill** — start with the [design planning skill][design-base]. At home, the
  [personal notes fragment][design-home] adds a short decision record. At work, the [review fragment][design-team]
  and [delivery fragment][design-delivery] add review and rollout guidance to the [rendered skill][design-rendered].
- **Override instructions** — the [shared working agreement][workflow-shared], [home version][workflow-home],
  and [work version][workflow-work] use the same filename. The last section wins.
- **Resolve installed paths** — [engineering instructions][engineering] use `{{%_skills_%}}` and `{{%_resources_%}}`.
  Sync resolves both placeholders against the selected harness’s destination.
- **Bundle a reference** — the [behavior diff fragment][behavior] uses the same placeholder with the team review
  checklist. Render copies the checklist into the skill and rewrites the link. See the [rendered skill][diff-rendered].
- **Link to another skill** — the design plan's delivery fragment links to `Agent_DeliveryPlan`. In the
  [rendered Codex plugin skill][design-plugin], the link uses the `example-skills` namespace.
- **Embed instructions in a skill** — the [Agent_Guidance template][guidance] uses `{{%_agents_%}}`.
  The [rendered skill][guidance-rendered] contains the assembled instructions and the resources they reference.
- **Target a harness** — [inline review][codex-only] uses `harnesses: [codex]` in its frontmatter.

Compare a fresh work render with the saved output.

```sh
(
  set -e
  for harness in codex claude opencode; do
    git diff --no-index -- "examples/output/$harness" ".tmp/example-work/latest/$harness"
  done
)
```

[conf-agents]: https://github.com/NonLogicalDev/conf.agents/tree/8515225/sections/00-personal-public
[design-base]: input-sections/00-shared/skills/Code_DesignPlan/SKILL.md
[design-home]: input-sections/10-home/skills/Code_DesignPlan/SKILL.md.d/20-planning-practices.md
[design-team]: input-sections/20-work/skills/Code_DesignPlan/SKILL.md.d/20-planning-practices.md
[design-delivery]: input-sections/20-work/skills/Code_DesignPlan/SKILL.md.d/30-delivery.md
[design-rendered]: output/claude/skills/Code_DesignPlan/SKILL.md
[workflow-shared]: input-sections/00-shared/AGENTS.md.d/40-working-agreement.md
[workflow-home]: input-sections/10-home/AGENTS.md.d/40-working-agreement.md
[workflow-work]: input-sections/20-work/AGENTS.md.d/40-working-agreement.md
[engineering]: input-sections/00-shared/AGENTS.md.d/30-engineering-practices.md
[behavior]: input-sections/20-work/skills/Code_BehaviorDiff/SKILL.md.d/20-review-findings.md
[diff-rendered]: output/claude/skills/Code_BehaviorDiff/SKILL.md
[design-plugin]: output/codex/plugin-skillset/Code_DesignPlan/SKILL.md
[guidance]: input-sections/00-shared/skills/Agent_Guidance/SKILL.md.d/10-skill.md
[guidance-rendered]: output/claude/skills/Agent_Guidance/SKILL.md
[codex-only]: input-sections/20-work/AGENTS.md.d/60-inline-review.md
