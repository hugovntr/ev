package cmd

import (
	"errors"
	"ev/crypto"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Initialize a new Environment Vault",
	RunE:  runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {

	if _, err := os.Stat(vaultFile); err == nil {
		return fmt.Errorf("vault file already exists at %s. Delete it manually if you want to start over", vaultFile)
	}

	// Register password
	fmt.Print("Enter new Vault password: ")
	bytePw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	fmt.Println()
	password := string(bytePw)

	if len(password) == 0 {
		return errors.New("password can not be empty")
	}

	fmt.Print("Confirm password: ")
	bytePwConfirm, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	fmt.Println()

	if password != string(bytePwConfirm) {
		return errors.New("passwords does not match")
	}

	// Create directories if they don't exists
	if err := os.MkdirAll(filepath.Dir(vaultFile), 0700); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(passwordFile), 0700); err != nil {
		return fmt.Errorf("failed to create credentials directory: %w", err)
	}

	// Save password
	err = os.WriteFile(passwordFile, []byte(password), 0600)
	if err != nil {
		return fmt.Errorf("failed to save credentials file: %w", err)
	}
	fmt.Printf("✔ Password saved to %s\n", passwordFile)

	// Initialize the vault with an empty YAML template
	initialContent := []byte(`# ENVIRONMENT VAULT
# Add your variables here in YAML format
#
# Example:
#   KEY: value
#
# Then you can use ` + "`ev get API_KEY`" + ` to get the value back as plain-text.
# Or ` + "`ev export`" + ` to get all the values and load them into a shell session


`)

	encryptedContent, err := crypto.Encrypt(initialContent, password)
	if err != nil {
		return fmt.Errorf("failed to encrypt initial content: %w", err)
	}

	err = os.WriteFile(vaultFile, encryptedContent, 0600)
	if err != nil {
		return fmt.Errorf("failed to save vault file: %w", err)
	}

	fmt.Printf("✔ Vault initialized at %s\n", vaultFile)
	return nil
}
