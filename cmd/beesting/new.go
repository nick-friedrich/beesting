package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	githubRepo   = "nick-friedrich/beesting"
	githubBranch = "main"
)

// newCmd creates a new application
var newCmd = &cobra.Command{
	Use:   "new [template] <name> or new <name>",
	Short: "Create a new application",
	Long: `Create a new application in the app directory.

Usage:
  beesting new <name>                 - Create with default template
  beesting new <template> <name>      - Create from GitHub template`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var templateName, projectName string

		if len(args) == 1 {
			// beesting new my-project (use default template)
			projectName = args[0]
			return createFromDefaultTemplate(projectName)
		}

		// beesting new example-api my-project (fetch from GitHub)
		templateName = args[0]
		projectName = args[1]
		return createFromGitHubTemplate(templateName, projectName)
	},
}

// createFromDefaultTemplate creates a project with the built-in template
func createFromDefaultTemplate(name string) error {
	// Create app directory
	appDir := filepath.Join("app", name)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create main.go
	mainGoPath := filepath.Join(appDir, "main.go")

	mainGoContent := fmt.Sprintf(`package main

import (
	"net/http"

	"github.com/nick-friedrich/beesting/pkg/beesting"
)

func main() {
	app := beesting.NewApp()

	// Add middleware
	app.Use(beesting.Logger())
	app.Use(beesting.Recovery())

	// Routes
	app.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from %s!"))
	})

	app.Run(":8080")
}
`, name)

	if err := os.WriteFile(mainGoPath, []byte(mainGoContent), 0644); err != nil {
		return fmt.Errorf("failed to create main.go: %w", err)
	}

	fmt.Printf("✓ Created new app: %s\n", name)
	fmt.Printf("  Location: %s\n", appDir)
	fmt.Printf("\nRun with: beesting dev %s\n", name)

	return nil
}

// createFromGitHubTemplate fetches a template from GitHub and creates a project
func createFromGitHubTemplate(templateName, projectName string) error {
	appDir := filepath.Join("app", projectName)

	// Check if directory already exists
	if _, err := os.Stat(appDir); err == nil {
		return fmt.Errorf("app '%s' already exists", projectName)
	}

	fmt.Printf("📥 Fetching template '%s' from GitHub...\n", templateName)

	// Download the tarball from GitHub
	url := fmt.Sprintf("https://github.com/%s/archive/refs/heads/%s.tar.gz", githubRepo, githubBranch)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download template: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download template: HTTP %d", resp.StatusCode)
	}

	// Extract the specific template folder
	if err := extractTemplate(resp.Body, templateName, appDir); err != nil {
		return fmt.Errorf("failed to extract template: %w", err)
	}

	fmt.Printf("✓ Created new app: %s (from template: %s)\n", projectName, templateName)
	fmt.Printf("  Location: %s\n", appDir)
	fmt.Printf("\nRun with: beesting dev %s\n", projectName)

	return nil
}

// extractTemplate extracts a specific folder from a tar.gz archive
func extractTemplate(r io.Reader, templateName, destDir string) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	templatePath := fmt.Sprintf("beesting-%s/app/%s/", githubBranch, templateName)
	found := false

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Check if this file is in our template folder
		if !strings.HasPrefix(header.Name, templatePath) {
			continue
		}

		found = true

		// Get relative path within template
		relPath := strings.TrimPrefix(header.Name, templatePath)
		if relPath == "" {
			continue
		}

		target := filepath.Join(destDir, relPath)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			// Create parent directory
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}

			// Create file
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}

	if !found {
		return fmt.Errorf("template '%s' not found in repository", templateName)
	}

	return nil
}
