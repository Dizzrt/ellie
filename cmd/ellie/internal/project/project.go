package project

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type Project struct {
	ProjectName string
	Module      string
	Path        string
	Repo        string
	Branch      string
}

func (p *Project) Create(ctx context.Context) error {
	if err := p.Clone(ctx); err != nil {
		return err
	}

	if err := p.cleanFiles(); err != nil {
		return err
	}

	module, err := modulePath(path.Join(p.Path, "go.mod"))
	if err != nil {
		return err
	}

	if err := p.replaceModule(module); err != nil {
		return err
	}

	return nil
}

func (p *Project) replaceModule(oldModule string) error {
	return filepath.WalkDir(p.Path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// skip dir
		if d.IsDir() {
			return nil
		}

		// read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// check file content whether contains oldModule
		contentStr := string(content)
		if !strings.Contains(contentStr, oldModule) {
			return nil
		}

		// replace module
		newContent := strings.ReplaceAll(contentStr, oldModule, p.Module)
		return os.WriteFile(path, []byte(newContent), d.Type())
	})
}

// cleanFiles cleans the unnecessary files
func (p *Project) cleanFiles() error {
	// needCleanFiles := []string{".git", ".gitkeep"}
	if err := os.RemoveAll(path.Join(p.Path, ".git")); err != nil {
		return err
	}

	if err := os.RemoveAll(path.Join(p.Path, "api", "gen")); err != nil {
		return err
	}

	if err := os.RemoveAll(path.Join(p.Path, "cmd", "wire_gen.go")); err != nil {
		return err
	}

	fmt.Printf("\n🍺 Service creation succeeded %s\n", color.GreenString(p.ProjectName))
	fmt.Print("💻 Try to start the service with the following command 👇:\n\n")

	fmt.Println(color.WhiteString("$ cd %s", p.ProjectName))
	fmt.Println(color.WhiteString("$ cd api && buf generate && cd ../cmd && wire gen && cd .."))
	fmt.Println(color.WhiteString("$ go build -o ./bin/app . && cp -r ./configs ./bin/"))
	fmt.Println(color.WhiteString("$ ./bin/app run\n"))

	return nil
}

// clones the template repository
func (p *Project) Clone(ctx context.Context) error {
	// if _, err := os.Stat(r.Path()); !os.IsNotExist(err) {
	// 	return r.Pull(ctx)
	// }

	fmt.Println(p.Path)

	var cmd *exec.Cmd
	if p.Branch == "" {
		cmd = exec.CommandContext(ctx, "git", "clone", p.Repo, p.Path)
	} else {
		cmd = exec.CommandContext(ctx, "git", "clone", "-b", p.Branch, p.Repo, p.Path)
	}

	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		return err
	}

	return nil
}
