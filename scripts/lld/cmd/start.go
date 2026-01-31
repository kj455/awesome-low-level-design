package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Select a problem and setup workspace",
	Long:  "Select a LLD problem using fzf and setup the workspace with problem description and main.go template.",
	Run: func(cmd *cobra.Command, args []string) {
		app, err := NewApp()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := runStart(app); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runStart(a *App) error {
	solved, total := a.CountSolved()
	fmt.Printf("  %d/%d completed\n\n", solved, total)

	lines := a.GenerateList(false)
	selected, err := RunFzf(lines)
	if err != nil {
		return err
	}
	if selected == "" {
		return nil
	}

	name := ExtractNameFromSelection(selected)
	p := a.FindProblemByName(name)
	if p == nil {
		return fmt.Errorf("problem not found: %s", name)
	}

	return a.SetupWorkspace(p)
}
