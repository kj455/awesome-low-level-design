package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// Difficulty represents the difficulty level of a problem.
type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
)

func (d Difficulty) String() string {
	return [...]string{"Easy", "Medium", "Hard"}[d]
}

// Problem represents a LLD problem.
type Problem struct {
	ID         string
	Difficulty Difficulty
	Name       string
}

const (
	umlSectionHeader = "## UML Class Diagram"
	mainGoFileName   = "main.go"
	problemMdName    = "problem.md"
	problemsDir      = "problems"
	workspaceSubdir  = "workspace/go"
	goModName        = "go.mod"
	moduleName       = "workspace"
)

const mainGoTemplate = `package main

import "fmt"

func main() {
	fmt.Println("TODO: Implement solution")
}
`

var selectionPrefixRe = regexp.MustCompile(`^\[.*?\]\s*`)

var defaultProblems = []Problem{
	// Easy
	{"parking-lot", Easy, "Parking Lot"},
	{"stack-overflow", Easy, "Stack Overflow"},
	{"vending-machine", Easy, "Vending Machine"},
	{"logging-framework", Easy, "Logging Framework"},
	{"traffic-signal", Easy, "Traffic Signal"},
	{"coffee-vending-machine", Easy, "Coffee Vending Machine"},
	{"task-management-system", Easy, "Task Management System"},
	// Medium
	{"atm", Medium, "ATM"},
	{"linkedin", Medium, "LinkedIn"},
	{"lru-cache", Medium, "LRU Cache"},
	{"tic-tac-toe", Medium, "Tic Tac Toe"},
	{"pub-sub-system", Medium, "Pub Sub System"},
	{"elevator-system", Medium, "Elevator System"},
	{"car-rental-system", Medium, "Car Rental System"},
	{"online-auction-system", Medium, "Online Auction System"},
	{"hotel-management-system", Medium, "Hotel Management System"},
	{"digital-wallet-service", Medium, "Digital Wallet Service"},
	{"airline-management-system", Medium, "Airline Management System"},
	{"library-management-system", Medium, "Library Management System"},
	{"social-networking-service", Medium, "Social Networking Service"},
	{"restaurant-management-system", Medium, "Restaurant Management System"},
	{"concert-ticket-booking-system", Medium, "Concert Ticket Booking System"},
	// Hard
	{"cricinfo", Hard, "Cricinfo"},
	{"splitwise", Hard, "Splitwise"},
	{"chess-game", Hard, "Chess Game"},
	{"snake-and-ladder", Hard, "Snake and Ladder"},
	{"ride-sharing-service", Hard, "Ride Sharing Service"},
	{"course-registration-system", Hard, "Course Registration System"},
	{"movie-ticket-booking-system", Hard, "Movie Ticket Booking System"},
	{"online-shopping-service", Hard, "Online Shopping Service"},
	{"online-stock-brokerage-system", Hard, "Online Stock Brokerage System"},
	{"music-streaming-service", Hard, "Music Streaming Service"},
	{"food-delivery-service", Hard, "Food Delivery Service"},
}

// App holds the application state and dependencies.
type App struct {
	rootDir      string
	workspaceDir string
	problems     []Problem
}

// NewApp creates a new App instance with detected root directory.
func NewApp() (*App, error) {
	rootDir, err := detectRootDir()
	if err != nil {
		return nil, fmt.Errorf("detect root directory: %w", err)
	}

	return &App{
		rootDir:      rootDir,
		workspaceDir: filepath.Join(rootDir, workspaceSubdir),
		problems:     defaultProblems,
	}, nil
}

func detectRootDir() (string, error) {
	execPath, err := os.Executable()
	if err == nil && !strings.Contains(execPath, "go-build") {
		// Running as compiled binary: scripts/lld/lld -> root
		return filepath.Dir(filepath.Dir(filepath.Dir(execPath))), nil
	}

	// Running via "go run": find .git directory to locate project root
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root, fallback to cwd
			return cwd, nil
		}
		dir = parent
	}
}

func (a *App) IsSolved(id string) bool {
	info, err := os.Stat(filepath.Join(a.workspaceDir, id))
	return err == nil && info.IsDir()
}

func (a *App) FindProblemByName(name string) *Problem {
	for i := range a.problems {
		if a.problems[i].Name == name {
			return &a.problems[i]
		}
	}
	return nil
}

func (a *App) CountSolved() (solved, total int) {
	total = len(a.problems)
	for _, p := range a.problems {
		if a.IsSolved(p.ID) {
			solved++
		}
	}
	return solved, total
}

func (a *App) GenerateList(solvedOnly bool) []string {
	var lines []string
	for _, p := range a.problems {
		solved := a.IsSolved(p.ID)
		if solvedOnly && !solved {
			continue
		}
		mark := ""
		if solved && !solvedOnly {
			mark = " ☑"
		}
		lines = append(lines, fmt.Sprintf("[%-6s] %s%s", p.Difficulty, p.Name, mark))
	}
	// sort.Strings(lines)
	return lines
}

func RunFzf(lines []string) (string, error) {
	if _, err := exec.LookPath("fzf"); err != nil {
		return "", errors.New("fzf is not installed. Install it with: brew install fzf")
	}

	cmd := exec.Command("fzf", "--height=40%", "--reverse")
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("create stdin pipe: %w", err)
	}

	errCh := make(chan error, 1)
	go func() {
		defer stdin.Close()
		for _, line := range lines {
			if _, err := fmt.Fprintln(stdin, line); err != nil {
				errCh <- err
				return
			}
		}
		errCh <- nil
	}()

	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 130 {
			// User cancelled with Ctrl+C
			return "", nil
		}
		return "", nil // fzf returns non-zero on no selection
	}

	if writeErr := <-errCh; writeErr != nil {
		return "", fmt.Errorf("write to fzf: %w", writeErr)
	}

	return strings.TrimSpace(string(output)), nil
}

func ExtractNameFromSelection(selection string) string {
	name := selectionPrefixRe.ReplaceAllString(selection, "")
	return strings.TrimSuffix(name, " ☑")
}

func (a *App) CopyProblemMd(problemID, destDir string) (err error) {
	srcPath := filepath.Join(a.rootDir, problemsDir, problemID+".md")
	destPath := filepath.Join(destDir, problemMdName)

	srcFile, err := os.Open(srcPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("open source file: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create destination file: %w", err)
	}
	defer func() {
		if cerr := destFile.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close destination file: %w", cerr)
		}
	}()

	scanner := bufio.NewScanner(srcFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, umlSectionHeader) {
			break
		}
		if _, err := fmt.Fprintln(destFile, line); err != nil {
			return fmt.Errorf("write line: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan source file: %w", err)
	}

	return nil
}

func WriteMainGo(path string) error {
	if err := os.WriteFile(path, []byte(mainGoTemplate), 0644); err != nil {
		return fmt.Errorf("write %s: %w", mainGoFileName, err)
	}
	return nil
}

func (a *App) InitGoMod() error {
	goModPath := filepath.Join(a.workspaceDir, goModName)
	if _, err := os.Stat(goModPath); err == nil {
		return nil
	}

	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = a.workspaceDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod init: %w", err)
	}
	return nil
}

func (a *App) SetupWorkspace(p *Problem) error {
	dir := filepath.Join(a.workspaceDir, p.ID)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	if err := a.CopyProblemMd(p.ID, dir); err != nil {
		return fmt.Errorf("copy problem.md: %w", err)
	}

	mainGoPath := filepath.Join(dir, mainGoFileName)
	if _, err := os.Stat(mainGoPath); os.IsNotExist(err) {
		if err := WriteMainGo(mainGoPath); err != nil {
			return err
		}
	}

	if err := a.InitGoMod(); err != nil {
		return err
	}

	fmt.Printf("\nCreated: %s\n\n", dir)
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", dir)
	fmt.Printf("  go run %s\n", mainGoFileName)

	return nil
}

func (a *App) ResetWorkspace(p *Problem) error {
	dir := filepath.Join(a.workspaceDir, p.ID)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("workspace not found: %s", dir)
	}

	mainGoPath := filepath.Join(dir, mainGoFileName)
	if err := WriteMainGo(mainGoPath); err != nil {
		return err
	}

	fmt.Printf("\nReset: %s\n", mainGoPath)
	return nil
}

var rootCmd = &cobra.Command{
	Use:   "lld",
	Short: "Low Level Design practice CLI",
	Long:  "A CLI tool to help practice Low Level Design problems.",
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior: run start command
		startCmd.Run(cmd, args)
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(resetCmd)
}
