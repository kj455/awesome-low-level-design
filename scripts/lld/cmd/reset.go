package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset main.go for a solved problem",
	Long:  "Select a solved LLD problem using fzf and reset its main.go to the template.",
	Run: func(cmd *cobra.Command, args []string) {
		app, err := NewApp()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := runReset(app); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runReset(a *App) error {
	fmt.Println("Select a problem to reset:\n")

	lines := a.GenerateList(true)
	if len(lines) == 0 {
		fmt.Println("No solved problems found.")
		return nil
	}

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

	return a.ResetWorkspace(p)
}
