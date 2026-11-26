package main

import (
	"os"
	"runtime/debug"

	"github.com/NonLogicalDev/tool.agentbld/internal"
)

var (
	version = "dev"
	commit  string
)

func main() {
	os.Exit(internal.Main(buildVersion()))
}

func buildVersion() string {
	build, revision := version, commit
	if info, ok := debug.ReadBuildInfo(); ok {
		if build == "dev" && info.Main.Version != "" && info.Main.Version != "(devel)" {
			build = info.Main.Version
		}
		if revision == "" {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					revision = setting.Value
					break
				}
			}
		}
	}
	result := "agentbld " + build
	if revision != "" {
		result += " (" + revision + ")"
	}
	return result
}
