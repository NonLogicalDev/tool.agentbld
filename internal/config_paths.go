package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/NonLogicalDev/tool.agentbld/internal/engine"
)

type configCommands struct {
	Paths struct{} `cmd:"" help:"Show effective and default paths."`
}

type configuredPaths struct {
	SectionsDirs   []string          `json:"sections_dirs"`
	RenderDir      string            `json:"render_dir"`
	Latest         string            `json:"latest"`
	TargetsDir     string            `json:"targets_dir"`
	RenderLock     string            `json:"render_lock"`
	ProfileDir     string            `json:"profile_dir"`
	DefaultProfile string            `json:"default_profile"`
	Destinations   map[string]string `json:"destinations"`
}

func printPaths(args cli) error {
	effective, err := resolvePaths(args)
	if err != nil {
		return err
	}
	defaults, err := resolvePaths(cli{})
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(struct {
		Effective configuredPaths `json:"effective"`
		Defaults  configuredPaths `json:"defaults"`
	}{Effective: effective, Defaults: defaults})
}

func resolvePaths(args cli) (configuredPaths, error) {
	var paths configuredPaths
	root, err := renderRoot(args.RenderDir)
	if err != nil {
		return paths, err
	}
	profiles, err := profileDirectory()
	if err != nil {
		return paths, err
	}
	paths.RenderDir = root
	paths.Latest = filepath.Join(root, engine.LatestAlias)
	paths.TargetsDir = filepath.Join(root, engine.TargetsDir)
	paths.RenderLock = filepath.Join(root, engine.RenderLockFile)
	paths.ProfileDir = profiles
	paths.DefaultProfile = filepath.Join(profiles, "default.json")
	paths.SectionsDirs = make([]string, 0, len(args.SectionsDirs))
	for _, dir := range args.SectionsDirs {
		absolute, err := engine.Absolute(dir)
		if err != nil {
			return paths, err
		}
		paths.SectionsDirs = append(paths.SectionsDirs, absolute)
	}
	registry := configuredHarnesses(&args)
	paths.Destinations = map[string]string{}
	for _, name := range engine.SortedKeys(registry) {
		destination, err := registry[name].ResolveDestination(nil)
		if err != nil {
			return paths, fmt.Errorf("%s paths: %w", name, err)
		}
		for key, path := range destination.Profile {
			paths.Destinations[key] = path
		}
	}
	return paths, nil
}
