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
	Long:  `Start an application in development mode with hot-reloading (using Air if available).`,
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

		// Get absolute path to app directory
		absAppDir, err := filepath.Abs(appDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}

		// Run templ generate
		templCmd := exec.Command("templ", "generate", absAppDir)
		templCmd.Stdout = os.Stdout
		templCmd.Stderr = os.Stderr
		templCmd.Stdin = os.Stdin

		if err := templCmd.Run(); err != nil {
			return fmt.Errorf("failed to run templ generate: %w", err)
		}

		// Check if Air is installed
		if isAirInstalled() {
			return runWithAir(name, absAppDir)
		}

		// Fall back to go run
		fmt.Printf("🐝 Starting %s in development mode...\n", name)
		fmt.Println("   (Install Air for hot-reloading: go install github.com/air-verse/air@latest)")

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

// isAirInstalled checks if Air is available in the PATH
func isAirInstalled() bool {
	_, err := exec.LookPath("air")
	return err == nil
}

// runWithAir runs the app with Air for hot-reloading
func runWithAir(name, appDir string) error {
	fmt.Printf("🐝 Starting %s with hot-reloading (Air)...\n\n", name)

	// Create temporary .air.toml for this app
	airConfig := fmt.Sprintf(`root = "%s"
tmp_dir = "%s/tmp"

[build]
  bin = "%s/tmp/main"
  cmd = "templ generate && go build -o %s/tmp/main %s"
  delay = 1000
  exclude_dir = ["tmp", "node_modules"]
  exclude_regex = ["_test.go", "_templ.go"]
  include_ext = ["go", "tpl", "tmpl", "html", "templ"]
  kill_delay = "500ms"
  rerun = false
  rerun_delay = 500
  send_interrupt = true
  stop_on_error = false

[color]
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[screen]
  clear_on_rebuild = true
`, appDir, appDir, appDir, appDir, appDir)

	// Write temporary config
	tmpConfigPath := filepath.Join(appDir, ".air.toml")
	if err := os.WriteFile(tmpConfigPath, []byte(airConfig), 0644); err != nil {
		return fmt.Errorf("failed to create Air config: %w", err)
	}

	// Ensure cleanup on exit
	defer os.Remove(tmpConfigPath)

	// Run Air
	airCmd := exec.Command("air", "-c", tmpConfigPath)
	airCmd.Dir = appDir
	airCmd.Stdout = os.Stdout
	airCmd.Stderr = os.Stderr
	airCmd.Stdin = os.Stdin

	if err := airCmd.Run(); err != nil {
		return fmt.Errorf("failed to run with Air: %w", err)
	}

	return nil
}
