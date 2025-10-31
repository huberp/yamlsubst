package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/huberp/yamlsubst/pkg/substitutor"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	yamlFile  string
	inputFile string
)

var rootCmd = &cobra.Command{
	Use:   "yamlsubst",
	Short: "Replace placeholders in input with values from YAML file",
	Long: `yamlsubst is a CLI tool similar to envsubst that replaces placeholders
in text input with values from a YAML file.

Placeholders use the format: ${.path.to.value}

Example:
  echo "Hello ${.name}" | yamlsubst --yaml values.yaml
  yamlsubst --yaml values.yaml --file template.txt`,
	RunE: run,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("yamlsubst version %s (commit: %s, built: %s)\n", version, commit, date)
	},
}

func init() {
	rootCmd.Flags().StringVar(&yamlFile, "yaml", "", "YAML file containing values for substitution (required)")
	rootCmd.Flags().StringVar(&inputFile, "file", "", "Input file containing placeholders (reads from stdin if not specified)")
	if err := rootCmd.MarkFlagRequired("yaml"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(versionCmd)
}

func run(cmd *cobra.Command, args []string) error {
	// Read YAML file
	yamlContent, err := os.ReadFile(yamlFile) // #nosec G304 -- CLI tool reads user-specified files
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Read input
	var input []byte
	if inputFile != "" {
		input, err = os.ReadFile(inputFile) // #nosec G304 -- CLI tool reads user-specified files
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}
	} else {
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
	}

	// Perform substitution
	result, err := substitutor.Substitute(string(input), string(yamlContent))
	if err != nil {
		return fmt.Errorf("substitution failed: %w", err)
	}

	// Output result
	fmt.Print(result)
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
