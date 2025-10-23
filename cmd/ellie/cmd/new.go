package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dizzrt/ellie/cmd/ellie/internal/project"
	"github.com/spf13/cobra"
)

var (
	p       project.Project
	timeout string
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "A brief description of your command",
	Long:  `A longer description`,
	Run:   run,
}

func init() {
	rootCmd.AddCommand(newCmd)

	timeout = "60s"
	if p.Repo = os.Getenv("ELLIE_LAYOUT_REPO"); p.Repo == "" {
		p.Repo = "https://github.com/dizzrt/ellie-layout.git"
	}

	newCmd.Flags().StringVarP(&p.Module, "module", "m", p.Module, "module name")
	newCmd.Flags().StringVarP(&p.Repo, "repo", "r", p.Repo, "layout repo")
	newCmd.Flags().StringVarP(&p.Branch, "branch", "b", p.Branch, "repo branch")
	newCmd.Flags().StringVarP(&timeout, "timeout", "t", timeout, "timeout")
}

func run(cmd *cobra.Command, args []string) {
	p.ProjectName = args[0]
	if p.ProjectName == "" {
		fmt.Fprint(os.Stderr, "\033[31mERROR: project name is required\033[m\n")
		return
	}

	if p.Module == "" {
		p.Module = p.ProjectName
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	p.Path = filepath.Join(wd, p.ProjectName)

	t, err := time.ParseDuration(timeout)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- p.Create(ctx)
	}()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			fmt.Fprint(os.Stderr, "\033[31mERROR: project creation timed out\033[m\n")
			return
		}
		fmt.Fprintf(os.Stderr, "\033[31mERROR: failed to create project(%s)\033[m\n", ctx.Err().Error())
	case err = <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: Failed to create project(%s)\033[m\n", err.Error())
		}
	}
}
