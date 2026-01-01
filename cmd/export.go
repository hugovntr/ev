package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export all variables in shell-compatible format",
	RunE:  runExport,
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	plaintext, _, err := loadVault()
	if err != nil {
		return err
	}

	data := make(map[string]any)
	if err := yaml.Unmarshal(plaintext, &data); err != nil {
		return fmt.Errorf("failed to parse vault YAML: %w", err)
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := data[k]
		val := fmt.Sprintf("%v", v)
		fmt.Fprintf(cmd.OutOrStdout(), "export %s=%q\n", k, val)
	}

	return nil
}
