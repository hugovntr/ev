package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var getCmd = &cobra.Command{
	Use:   "get [variable_name]",
	Short: "Retrieve a specific value from the Environment Vault",
	Args:  cobra.ExactArgs(1),
	RunE:  runGet,
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func runGet(cmd *cobra.Command, args []string) error {
	key := args[0]

	plaintext, _, err := loadVault()
	if err != nil {
		return err
	}

	data := make(map[string]any)
	if err := yaml.Unmarshal(plaintext, &data); err != nil {
		return fmt.Errorf("failed to parse vault YAML: %w", err)
	}

	val, ok := data[key]
	if !ok {
		return fmt.Errorf("key '%s' not found in vault", key)
	}

	strVal := fmt.Sprintf("%v", val)
	fmt.Fprint(cmd.OutOrStdout(), strings.TrimSpace(strVal))
	return nil
}
