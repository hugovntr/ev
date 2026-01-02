package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var envCmd = &cobra.Command{
	Use:   "env [shell]",
	Short: "Export all variables in shell-compatible format (bash, zsh, fish, powershell)",
	Long: `Output shell commands to set environment variables.
Supported shells: bash, zsh, sh, fish, powershell, cmd.
Default: bash`,
	ValidArgs: []string{"bash", "zsh", "sh", "fish", "powershell", "pwsh", "cmd"},
	RunE:      runEnv,
}

func init() {
	rootCmd.AddCommand(envCmd)
}

func runEnv(cmd *cobra.Command, args []string) error {

	target := "bash"
	if len(args) > 0 {
		target = strings.ToLower(args[0])
	}

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
		qval := fmt.Sprintf("%q", val)

		switch target {
		case "fish":
			fmt.Fprintf(cmd.OutOrStdout(), "set -gx %s %s\n", k, qval)
		case "powershell", "pwsh":
			fmt.Fprintf(cmd.OutOrStdout(), "$Env:%s = %s\n", k, qval)
		case "cmd":
			fmt.Fprintf(cmd.OutOrStdout(), "set %s=%s\n", k, val)
		case "bash", "zsh", "sh":
			fmt.Fprintf(cmd.OutOrStdout(), "export %s=%s\n", k, qval)
		default:
			return fmt.Errorf("unsupported shell: %s", target)
		}
	}

	return nil
}
