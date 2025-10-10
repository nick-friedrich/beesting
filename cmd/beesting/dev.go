package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// devCmd starts an application in development mode
var devCmd = &cobra.Command{
	Use:   "dev [name]",
	Short: "Start an application in development mode",
	Long:  `Start an application in development mode by running its main.go file.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Check if app directory exists
		appDir := filepath.Join("app", name)
		if _, err := os.Stat(appDir); os.IsNotExist(err) {
			return fmt.Errorf("app '%s' does not exist. Create it with: beesting new %s", name, name)
		}

		// Check if main.go exists
		mainGoPath := filepath.Join(appDir, "main.go")
		if _, err := os.Stat(mainGoPath); os.IsNotExist(err) {
			return fmt.Errorf("main.go not found in app '%s'", name)
		}

		fmt.Printf("🐝 Starting %s in development mode...\n\n", name)

		// Run the app with go run
		runCmd := exec.Command("go", "run", mainGoPath)
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		runCmd.Stdin = os.Stdin

		if err := runCmd.Run(); err != nil {
			return fmt.Errorf("failed to run app: %w", err)
		}

		return nil
	},
}
