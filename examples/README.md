# Layered configuration example

Start with shared instructions and skills from [NonLogicalDev/conf.agents][conf-agents]. Add personal notes
for home or review and delivery guidance for work. Both setups extend the same design planning skill.

## What to explore

- **Shared**
  - [Collaboration](input-sections/shared/00-collaboration/) covers user intent, audience, and Agent_Guidance.
  - [Writing](input-sections/shared/10-writing/) brings in prose guidance and reference material.
  - [Engineering](input-sections/shared/20-engineering/) covers design plans, behavior diffs, and checkpoints.
- **Home**
  - [Personal workflow](input-sections/home/30-personal-workflow/) adds decision notes and a way to resume projects.
- **Work**
  - [Team review](input-sections/work/30-team-review/) adds reviewer guidance and a bundled review checklist.
  - [Project delivery](input-sections/work/40-project-delivery/) adds landing plans and verification guidance.

The [pinned public configuration][conf-agents] supplies the source material. The environment fragments adapt it
for home and work.

## Render both setups

Run from the repository root with `agentbld` installed. Pass the shared source first and one environment second.

```sh
render_example() {
  variant=$1
  agentbld render --sections-dir examples/input-sections/shared \
    --sections-dir "examples/input-sections/$variant" --render-dir ".tmp/example-$variant" \
    --mode='codex:{"agents":true,"skills":true,"plugin":"full","plugin_name":"example-skills"}' \
    --mode=claude --mode=opencode
}

render_example home
render_example work
```

Each render has its own `latest/` directory under `.tmp/`. [The saved output](output/) shows the work setup.
Compare the two versions of the design planning skill.

```sh
diff -u .tmp/example-home/latest/claude/skills/Code_DesignPlan/SKILL.md \
  .tmp/example-work/latest/claude/skills/Code_DesignPlan/SKILL.md
```

## Reference and skill examples

- **Extend a shared skill** — start with the [design planning skill][design-base]. At home, the
  [personal notes fragment][design-home] adds a short decision record. At work, the [review fragment][design-team]
  and [delivery fragment][design-delivery] add review and rollout guidance to the [rendered skill][design-rendered].
- **Override instructions** — [home][workflow-home] and [work][workflow-work] each replace the
  [shared working agreement][workflow-shared] using the same filename.
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
[design-base]: input-sections/shared/20-engineering/skills/Code_DesignPlan/SKILL.md
[design-home]: input-sections/home/30-personal-workflow/skills/Code_DesignPlan/SKILL.md.d/20-personal-notes.md
[design-team]: input-sections/work/30-team-review/skills/Code_DesignPlan/SKILL.md.d/20-review-readiness.md
[design-delivery]: input-sections/work/40-project-delivery/skills/Code_DesignPlan/SKILL.md.d/30-delivery.md
[design-rendered]: output/claude/skills/Code_DesignPlan/SKILL.md
[workflow-shared]: input-sections/shared/00-collaboration/AGENTS.md.d/40-working-agreement.md
[workflow-home]: input-sections/home/30-personal-workflow/AGENTS.md.d/40-working-agreement.md
[workflow-work]: input-sections/work/30-team-review/AGENTS.md.d/40-working-agreement.md
[engineering]: input-sections/shared/20-engineering/AGENTS.md.d/30-engineering-practices.md
[behavior]: input-sections/work/30-team-review/skills/Code_BehaviorDiff/SKILL.md.d/20-review-findings.md
[diff-rendered]: output/claude/skills/Code_BehaviorDiff/SKILL.md
[design-plugin]: output/codex/plugin-skillset/Code_DesignPlan/SKILL.md
[guidance]: input-sections/shared/00-collaboration/skills/Agent_Guidance/SKILL.md.d/10-skill.md
[guidance-rendered]: output/claude/skills/Agent_Guidance/SKILL.md
[codex-only]: input-sections/work/40-project-delivery/AGENTS.md.d/60-inline-review.md
