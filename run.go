package main

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
)

func run(ctx *cli.Context) error {
	paths := ctx.Args().Slice()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	log.Info("Stating update")
	for _, path := range paths {
		err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() || info.Name() != "go.mod" {
				return nil
			}

			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			dir := filepath.Dir(absPath)
			log.Infof("Updating in: %s", dir)

			cmd := exec.Command("go", "-C", dir, "list", "-f", "{{if not (or .Main .Indirect)}}{{.Path}}{{end}}", "-m", "all")
			cmd.Stderr = os.Stderr
			out, err := cmd.Output()
			if err != nil {
				return err
			}

			modules := strings.Split(strings.TrimSpace(string(out)), "\n")
			log.Infof("Updating %d modules", len(modules))
			log.Debugf("Modules: %s", modules)

			runModTidy(dir)

			for _, module := range modules {
				cmd = exec.Command("go", "-C", dir, "get", module)
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr
				err = cmd.Run()
				if err != nil {
					log.Errorf("Failed to update %s, err: %s", module, err)
					break
				}
			}

			runModTidy(dir)

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
	cmd := exec.Command("go", "-C", dir, "mod", "tidy")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Failed to tidy, err: %s", err)
	}
}
