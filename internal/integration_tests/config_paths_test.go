package integration_tests

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func TestConfigPathsReportsOverridesAndDefaultsWithoutWrites(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	state := filepath.Join(root, "state")
	config := filepath.Join(root, "config")
	data := filepath.Join(root, "data")
	cmd := exec.Command(contractBinary, "config", "paths",
		"--sections-dir", "first", "--sections-dir", "second",
		"--render-dir", "custom-state", "--codex-dir", "codex-custom",
		"--claude-marketplace-dir", "marketplace",
	)
	cmd.Dir = root
	cmd.Env = append(os.Environ(),
		"HOME="+home, "XDG_STATE_HOME="+state, "XDG_CONFIG_HOME="+config, "XDG_DATA_HOME="+data,
		"CODEX_HOME="+filepath.Join(root, "codex-env"),
		"CLAUDE_CONFIG_DIR="+filepath.Join(root, "claude-env"),
		"OPENCODE_CONFIG_DIR=",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("config paths: %v\n%s", err, out)
	}

	type paths struct {
		SectionsDirs   []string          `json:"sections_dirs"`
		RenderDir      string            `json:"render_dir"`
		Latest         string            `json:"latest"`
		TargetsDir     string            `json:"targets_dir"`
		RenderLock     string            `json:"render_lock"`
		ProfileDir     string            `json:"profile_dir"`
		DefaultProfile string            `json:"default_profile"`
		Destinations   map[string]string `json:"destinations"`
	}
	var got struct {
		Effective paths `json:"effective"`
		Defaults  paths `json:"defaults"`
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("invalid paths JSON: %v\n%s", err, out)
	}
	defaults := paths{
		SectionsDirs:   []string{},
		RenderDir:      filepath.Join(state, "agentbld"),
		Latest:         filepath.Join(state, "agentbld", "latest"),
		TargetsDir:     filepath.Join(state, "agentbld", ".targets"),
		RenderLock:     filepath.Join(state, "agentbld", ".agentbld.lock"),
		ProfileDir:     filepath.Join(config, "agentbld"),
		DefaultProfile: filepath.Join(config, "agentbld", "default.json"),
		Destinations: map[string]string{
			"codex_dir":              filepath.Join(root, "codex-env"),
			"codex_marketplace_home": filepath.Join(home, ".agents", "plugins"),
			"claude_dir":             filepath.Join(root, "claude-env"),
			"claude_marketplace_dir": filepath.Join(data, "agentbld", "claude-marketplace"),
			"opencode_dir":           filepath.Join(config, "opencode"),
		},
	}
	if !reflect.DeepEqual(got.Defaults, defaults) {
		t.Fatalf("defaults = %#v, want %#v", got.Defaults, defaults)
	}
	effective := defaults
	effective.SectionsDirs = []string{filepath.Join(root, "first"), filepath.Join(root, "second")}
	effective.RenderDir = filepath.Join(root, "custom-state")
	effective.Latest = filepath.Join(effective.RenderDir, "latest")
	effective.TargetsDir = filepath.Join(effective.RenderDir, ".targets")
	effective.RenderLock = filepath.Join(effective.RenderDir, ".agentbld.lock")
	effective.Destinations["codex_dir"] = filepath.Join(root, "codex-custom")
	effective.Destinations["claude_marketplace_dir"] = filepath.Join(root, "marketplace")
	if !reflect.DeepEqual(got.Effective, effective) {
		t.Fatalf("effective paths = %#v, want %#v", got.Effective, effective)
	}
	if entries, err := os.ReadDir(root); err != nil || len(entries) != 0 {
		t.Fatalf("config paths wrote to its working directory: %v, %v", entries, err)
	}
}
