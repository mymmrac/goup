package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
)

func run(ctx *cli.Context) error {
	includeHidden := ctx.Bool("all")
	includeVendor := ctx.Bool("vendor")
	excludePatterns := ctx.StringSlice("exclude")

	paths := ctx.Args().Slice()
	if len(paths) == 0 {
		paths = []string{"."}
	}
	for i := 0; i < len(paths); i++ {
		paths[i] = filepath.Clean(paths[i])
	}
	log.Debugf("Lookup paths: %s", paths)

	log.Info("Stating update")
	for _, path := range paths {
		err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("walk dir: %w", err)
			}

			isDir := d.IsDir()

			if !includeHidden && isDir && containsHiddenDir(path) {
				log.Debugf("Skip hidden directory: %s", path)
				return filepath.SkipDir
			}

			pathBase := filepath.Base(path)

			if !includeVendor && isDir && pathBase == "vendor" {
				log.Debugf("Skip vendor directory: %s", path)
				return filepath.SkipDir
			}

			if isDir {
				for _, pattern := range excludePatterns {
					var match bool
					match, err = filepath.Match(pattern, pathBase)
					if err != nil {
						return fmt.Errorf("invalid exclude pattern: %w", err)
					}

					if match {
						log.Debugf("Skip excluded directory: %s", path)
						return filepath.SkipDir
					}
				}
			}

			if isDir || d.Name() != "go.mod" {
				return nil
			}

			pathAbs, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			pathDir := filepath.Dir(pathAbs)
			log.Infof("In: %s", pathDir)

			cmd := exec.Command("go", "list", "-f", "{{if not (or .Main .Indirect)}}{{.Path}}{{end}}", "-m", "all")
			cmd.Dir = pathDir
			cmd.Stderr = os.Stderr
			out, err := cmd.Output()
			if err != nil {
				log.Errorf("Get modules list, err: %s", err)
				return nil
			}

			modules := strings.Split(strings.TrimSpace(string(out)), "\n")
			log.Infof("Updating %d modules", len(modules))
			log.Debugf("Modules: %s", modules)

			runModTidy(pathDir)

			for _, module := range modules {
				cmd = exec.Command("go", "get", module)
				cmd.Dir = pathDir
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr
				err = cmd.Run()
				if err != nil {
					log.Errorf("Failed to update %s, err: %s", module, err)
					break
				}
			}

			runModTidy(pathDir)

			return nil
		})
		if err != nil {
			return err
		}
	}
	log.Info("Done")
	return nil
}

func runModTidy(dir string) {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Failed to tidy, err: %s", err)
	}
}

func containsHiddenDir(path string) bool {
	parts := strings.Split(path, string(filepath.Separator))
	for _, part := range parts {
		if len(part) > 1 && strings.HasPrefix(part, ".") && !strings.HasPrefix(part, "..") {
			return true
		}
	}
	return false
}
