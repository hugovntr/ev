package cmd

import (
	"ev/crypto"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the Environment Vault content in your default editor",
	RunE:  runEdit,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func runEdit(cmd *cobra.Command, args []string) error {

	plaintext, password, err := loadVault()
	if err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp("", "ev-*.yaml")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	defer os.Remove(tmpPath)

	if _, err := tmpFile.Write(plaintext); err != nil {
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	editorEnv := os.Getenv("EDITOR")
	if editorEnv == "" {
		editorEnv = "vi"
		fmt.Println("No $EDITOR set, defaulting to `vi`")
	}

	parts := strings.Fields(editorEnv)
	executable := parts[0]
	editorArgs := parts[1:]
	editorArgs = append(editorArgs, tmpPath)

	cmdEdit := exec.Command(executable, editorArgs...)

	// Connect the editor to the user's terminal
	cmdEdit.Stdin = os.Stdin
	cmdEdit.Stdout = os.Stdout
	cmdEdit.Stderr = os.Stderr

	if err := cmdEdit.Run(); err != nil {
		return fmt.Errorf("editor execution failed: %w", err)
	}

	newData, err := os.ReadFile(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to read back temp file: %w", err)
	}

	encryptedNewData, err := crypto.Encrypt(newData, password)
	if err != nil {
		return fmt.Errorf("failed to encrypt new content: %w", err)
	}
	if err := os.WriteFile(vaultFile, encryptedNewData, 0600); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	fmt.Println("✔ Environment Vault updated successfully.")
	return nil
}
