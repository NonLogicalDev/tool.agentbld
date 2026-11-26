# agentbld

Agentbld assembles shared instructions, skills, and reference files into the configuration each coding agent
expects. You can layer in your own preferences and adapt the same setup for different teams or projects.

Start with a skill you found or wrote. Keep the base skill shareable and add instructions for the way you work.
For example, use a shared design planning skill at home and extend it at work with your team’s review and
rollout guidance.
Agentbld combines the base skill with the sections you choose when you render.

<p align="center">
  <img src="assets/mascot.png" width="256" alt="Agentbld mascot with a folder connected to coding agent icons">
</p>

`render` creates a numbered, read-only snapshot called a **generation**. You can inspect it before using
`sync <number>` to activate it in the tools you choose. To return to an earlier setup, sync its saved generation
again, even if the source files have changed.

## Try it

The bundled example combines a shared base with instructions for work. From an agentbld checkout, run this with
Nix to preview the result. If you installed the binary, use `agentbld` in place of `nix run . --`.

```sh
demo=.tmp/agentbld-demo
nix run . -- render --sections-dir examples/input-sections/shared \
  --sections-dir examples/input-sections/work --render-dir "$demo" \
  --mode=codex --mode=claude --mode=opencode
ls "$demo/latest"
```

Open a generated instruction file such as `$demo/latest/codex/AGENTS.md` to see the assembled instructions.
This command writes to the demo directory and leaves your agent configuration directories alone.

- [Example inputs](examples/input-sections/) show the shared sections and the separate home and work additions.
- [Rendered work output](examples/output/) shows the assembled instructions, skills, and bundled references.
- [Example walkthrough](examples/README.md) shows how to render both setups and compare the same skill at home
  and at work.

## Run or install

### Nix flake

```sh
nix run github:NonLogicalDev/tool.agentbld -- --help
nix profile install github:NonLogicalDev/tool.agentbld
```

From a checkout, run `nix run . -- --version` or build with `nix build .`.

### Nix without flakes

From a checkout with `nixpkgs` in `NIX_PATH`:

```sh
nix-build
./result/bin/agentbld --version
```

Both entry points use [nix/package.nix](nix/package.nix). To reuse the recipe in a Nix configuration:

```nix
pkgs.callPackage ./nix/package.nix { }
```

The recipe accepts `source`, `version`, and `revision` overrides for a pinned package.

### Go

**Prerequisite** — Go 1.24 or later. Add `$GOBIN` to `PATH`, or use `$(go env GOPATH)/bin` if `GOBIN` is unset.

| Use | Version suffix |
| --- | --- |
| Track the main branch | `@main` |
| Pin a release once tagged | `@vX.Y.Z` |

```sh
go run github.com/NonLogicalDev/tool.agentbld/cmd/agentbld@main --help
go install github.com/NonLogicalDev/tool.agentbld/cmd/agentbld@main
```

## Configuration layout

A section is a numbered directory such as `10-team`. Agentbld processes sections in name order within each source
directory. Pass `--sections-dir` more than once to add source directories in order, starting with the shared base.

Later sections can replace an instruction fragment or resource with the same path, or add new fragments to a skill.
Filter a fragment with `harnesses` in its YAML frontmatter. See the [source sections guide][source-sections].

```text
sections/
└── {nn}-{ns}-{suffix}/
    ├── AGENTS.md.d/
    │   └── {order}-{section}.md
    ├── skills/
    │   └── {skill-name}/
    │       ├── SKILL.md.d/
    │       │   ├── 10-skill.md
    │       │   └── {order}-{fragment}.md
    │       ├── references/
    │       └── scripts/
    ├── resources/
    │   └── {resource-group}/
    │       └── {resource-name}
    └── opencode/plugins/
        └── {plugin-name}.js
```

- Override the generation directory with `--render-dir`.
- `config.json` records resolved modes. Destination paths are chosen at sync time.
- The selected source files and resolved options determine whether render reuses the newest generation.

```text
${XDG_STATE_HOME:-$HOME/.local/state}/agentbld/
├── latest -> {yyyy}-{mm}-{dd}__{NN}/
└── {yyyy}-{mm}-{dd}__{NN}/
    ├── config.json
    └── <harness>/       # One directory per selected harness
        ├── <instruction-file>
        ├── skills/
        └── resources/
```

Each harness defines its instruction filename and can add files such as command rules or plugin packages.

### Skill fragments

- Put a complete skill in `skills/<name>/SKILL.md` or split it into ordered files under `SKILL.md.d/`.
- To extend a shared skill, put a fragment under `skills/<name>/SKILL.md.d/` in a later section. Use the same
  skill name. Keep work instructions in a separate source directory and include it with `--sections-dir` at work.
- The example extends the [shared design planning skill][design-base] with a [team review fragment][design-team]
  and a delivery fragment. Render combines them into [one `SKILL.md`][design-rendered].

### Reference placeholders

- **Skill path in instructions** — `{{%_skills_%}}/Code_DesignPlan/SKILL.md` resolves at sync time.
  - With `--codex-dir ~/.codex-work`, it becomes `~/.codex-work/skills/Code_DesignPlan/SKILL.md`.
- **Resource path in instructions** — `{{%_resources_%}}/engineering-practices/small-changes.md` resolves under
  the harness’s `resources/` directory at sync time.
- **Resource path in a skill** — `{{%_resources_%}}/engineering-practices/team-review.md` becomes
  `resources/engineering-practices/team-review.md` at render time. Agentbld bundles the checklist inside the skill.
- **Link to a skill** — use `[$@:VCS_Checkpoint](skill://@/VCS_Checkpoint)` in instructions or another skill.
  - Render resolves it for the harness and plugin. In the example Codex plugin, the link points to
    `skill://example-skills/VCS_Checkpoint`.
- **Embed instructions in a skill** — `{{%_agents_%}}` in `Agent_Guidance` inserts the assembled instructions.

See the [runnable examples](examples/#reference-and-skill-examples) for the source fragments and generated files.

## Operations

| Command | Use it to | Input |
| --- | --- | --- |
| `render` | Create or reuse a generation | Source sections |
| `render --sync` | Render and activate that exact generation | Source sections and destinations |
| `check` | Verify the newest generation and `latest` without writing | Source sections and render options |
| `sync <number>` | Activate a saved generation | Ordinal and destinations |
| `modes` | List harness options and examples | — |
| `prune` | Keep the five newest generations | State root with no active links to older generations |
| `reset` | Back up Codex config and remove stale managed links | Codex destination and state root |

- **Check exit codes**
  - `0` — current
  - `1` — stale
  - `2` — error
- **Numbered sync**
  - Uses the saved files and options. Omit source flags.
  - Leaves `latest` pointing to the newest generation.
  - Repeat bare `--mode` names to select harnesses.
  - Requires the current generation layout and `config.json`. Render into a fresh state directory when upgrading.

```sh
# Run from the directory that contains your sections/.
agentbld render --sections-dir "$PWD/sections"
agentbld check --sections-dir "$PWD/sections"
agentbld render --sync --sections-dir "$PWD/sections"
agentbld sync 1  # Replace 1 with the ordinal printed by render.
```

## Configuration paths

```sh
agentbld config paths
agentbld config paths --sections-dir ./sections --render-dir ~/agentbld-state --codex-dir ~/.codex-work
```

- `effective` includes path flags and the current environment.
- `defaults` uses the current environment with all path flags omitted. There is no default sections directory.
- Both list the state paths, profile directory, and every harness’s destinations, including marketplaces.
- The command reads path options only. It creates no files and loads no saved profile.

## Sync profiles

```sh
agentbld sync 1 --mode=codex --codex-dir /path/to/codex --save-profile work
agentbld sync 2 --profile work
agentbld profile list
agentbld profile show work
```

| Setting | Behavior |
| --- | --- |
| Directory | `${XDG_CONFIG_HOME:-$HOME/.config}/agentbld/` |
| Default file | `default.json`, selected with `--profile default` |
| Named file | `<name>.json`, selected with `--profile <name>` |
| List profiles | `agentbld profile list` prints names in alphabetical order |
| Inspect a profile | `agentbld profile show NAME` prints its JSON. Use `default` to inspect `default.json`. |
| Save | `--save-profile <name>` creates or replaces the file after successful activation |
| Precedence | Explicit sync flags, then the selected profile, then the normal defaults |

- **Saved fields**
  - `render_dir`
  - `modes`, a list of selected harness names
  - Destination fields for the selected harnesses
- **Portable paths**
  - Saved paths under your home directory start with `~/`. Sync expands them against the current home.
  - Paths outside your home stay absolute.
- **Selection** — omit `--profile` to use regular flags and defaults.
- **Generation** — pass its ordinal each time. Each generation supplies its own mode options.
- **Mode overrides** — supplying `--mode` replaces the profile's entire harness list for that invocation.
- **Names** — start with a letter or digit and use letters, digits, underscores, or hyphens.

Example `work.json`

```json
{
  "render_dir": "~/.local/state/agentbld",
  "modes": ["codex"],
  "codex_dir": "~/.codex-work",
  "codex_marketplace_home": "~/.agents/plugins"
}
```

Save or select the profile named `default`.

```sh
agentbld sync 2 --profile work --save-profile default
agentbld sync 3 --profile default
```

## Harnesses

Run `agentbld modes` to see the harnesses and options available in your build.

Use `--mode=<harness>[:<JSON object>]` once per harness to choose what it receives. When you sync, agentbld resolves
path placeholders for each destination and links the selected files into its configuration directory. It reports a
conflict if an unmanaged file is already at a path it needs.

Codex and Claude Code can also package source skills as a plugin. Render saves the plugin’s skills under
`plugin-skillset/`. Sync writes the local package so you can register or install it through the harness.

### Codex

| Sync flag | Default |
| --- | --- |
| `--codex-dir` | `${CODEX_HOME:-$HOME/.codex}` |
| `--codex-marketplace-home` | `$HOME/.agents/plugins` |

See [Codex examples, mode options and command rules][codex-config].

### OpenCode

| Mode option | Selects | Default |
| --- | --- | --- |
| `agents` | `AGENTS.md` and resources | On |
| `skills` | Standalone skills | On |
| `plugins` | Native plugins | On |

`--opencode-dir` overrides `$OPENCODE_CONFIG_DIR` or `${XDG_CONFIG_HOME:-$HOME/.config}/opencode`.

### Claude Code

| Plugin selection | `agents` default | `skills` default | Plugin includes `Agent_Guidance` |
| --- | --- | --- | --- |
| No plugin | On | On | — |
| `"skills"` | On | Off | No |
| `"full"` | Off | Off | Yes |

- **`agents`** selects `CLAUDE.md` and resources. **`skills`** selects standalone skills.
- **Overrides**
  - Set `agents` and `skills` explicitly to override the defaults.
  - Set `plugin_name` to choose a name. Defaults to `skillset`.
- **Standalone links**
  - Plugin-only sync leaves the Claude configuration directory untouched.
  - When syncing instructions and a plugin with `skills` off, sync removes previously managed skill links.

| Sync flag | Default |
| --- | --- |
| `--claude-dir` | `${CLAUDE_CONFIG_DIR:-$HOME/.claude}` |
| `--claude-marketplace-dir` | `${XDG_DATA_HOME:-$HOME/.local/share}/agentbld/claude-marketplace` |

**Install the generated plugin after sync**

```sh
claude plugin marketplace add "${XDG_DATA_HOME:-$HOME/.local/share}/agentbld/claude-marketplace"
claude plugin install skillset@agentbld
```

- Use your configured path and names if they differ. See the
  [Claude Code marketplace documentation](https://code.claude.com/docs/en/plugin-marketplaces).

## Development

- **Entry point and CLI** — `cmd/agentbld/main.go` and `internal/`.
- **Shared engine code** — `internal/engine/` handles source layering, generation storage, replay, and managed links.
  - Shared filesystem helpers also live here.
- **Harness adapters** — to add a harness, create a package under [`internal/harnesses/`](internal/harnesses/).
  - Implement the [`engine.Handler`](internal/engine/harness.go) operations for rendering and syncing its files.
  - Register the handler and any CLI options in [`internal/harness_registry.go`](internal/harness_registry.go).
- **Tests**
  - Unit tests live beside the code they exercise.
  - CLI integration tests build and run the binary from `internal/integration_tests/`.
  - Run `go test ./...` from the repository root.

## Releases

- Push a semantic version tag such as `v0.1.0` to trigger the GitHub release workflow.
- GoReleaser publishes macOS arm64 and Linux amd64 and arm64 archives with SHA-256 checksums.
- `agentbld --version` prints the release version and full commit hash.
- Build release artifacts locally with a snapshot.

```sh
goreleaser release --config .github/actions/goreleaser.yaml --snapshot --clean
```

## History

Agentbld was extracted from [NonLogicalDev/conf.agents][original-source] with its Git history.

[source-sections]: https://github.com/NonLogicalDev/conf.agents#source-sections
[codex-config]: https://github.com/NonLogicalDev/conf.agents#codex
[original-source]: https://github.com/NonLogicalDev/conf.agents/tree/819153a/utils/agentbld
[design-base]: examples/input-sections/shared/20-engineering/skills/Code_DesignPlan/SKILL.md
[design-team]: examples/input-sections/work/30-team-review/skills/Code_DesignPlan/SKILL.md.d/20-review-readiness.md
[design-rendered]: examples/output/claude/skills/Code_DesignPlan/SKILL.md
