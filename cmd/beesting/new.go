package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// newCmd creates a new application
var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new application",
	Long:  `Create a new application in the app directory with a basic main.go template.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Create app directory
		appDir := filepath.Join("app", name)
		if err := os.MkdirAll(appDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Create main.go
		mainGoPath := filepath.Join(appDir, "main.go")

		mainGoContent := fmt.Sprintf(`package main

import "fmt"

func main() {
	fmt.Println("Hello from %s!")
}
`, name)

		if err := os.WriteFile(mainGoPath, []byte(mainGoContent), 0644); err != nil {
			return fmt.Errorf("failed to create main.go: %w", err)
		}

		fmt.Printf("✓ Created new app: %s\n", name)
		fmt.Printf("  Location: %s\n", appDir)
		fmt.Printf("\nRun with: beesting dev %s\n", name)

		return nil
	},
}
