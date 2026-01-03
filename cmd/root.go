package cmd

import (
	"ev/crypto"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// Global variables to store flag values
var (
	vaultFile    string
	passwordFile string
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "ev",
	Short:   "Environment Vault: Manage secure environment variables",
	Version: version,
}

// Execute is the entry point for main.go
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Define default paths
	home, _ := os.UserHomeDir()
	defaultVault := filepath.Join(home, ".config/ev/vault")
	defaultCreds := filepath.Join(home, ".config/ev/credentials")

	// PersistentFlags make these available to all subcommands (create, edit, get)
	rootCmd.PersistentFlags().StringVarP(&vaultFile, "vault-file", "f", defaultVault, "Path to the vault file")
	rootCmd.PersistentFlags().StringVarP(&passwordFile, "password-file", "p", defaultCreds, "Path to the password file")
}

// resolvePassword determines how to get the password based on your 3 rules.
// 1. Check file at 'passwordFile' path (default or flag override).
// 2. (Technically covered by 1 since we have defaults).
// 3. Prompt user if file doesn't exist.
func resolvePassword(promptUser bool) (string, error) {
	// Try reading from the password file first
	if _, err := os.Stat(passwordFile); err == nil {
		content, err := os.ReadFile(passwordFile)
		if err != nil {
			return "", fmt.Errorf("failed to read password file: %w", err)
		}
		return strings.TrimSpace(string(content)), nil
	}

	// If file missing and promptUser is false, we error out (e.g., in scripts)
	if !promptUser {
		return "", fmt.Errorf("password file not found at %s", passwordFile)
	}

	// Fallback: Ask the user manually
	fmt.Print("Enter Vault Password: ")
	bytePw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Print newline after input
	if err != nil {
		return "", err
	}
	return string(bytePw), nil
}

// loadVault handle the common workflow of checking the vault file,
// getting the password, and decrypting the content.
func loadVault() ([]byte, string, error) {
	// Check if vault exists
	if _, err := os.Stat(vaultFile); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("vault file not found at %s. Run 'ev create' first", vaultFile)
	}

	// Get Password
	// We default to prompting (true) because this is a helper for CLI commands.
	password, err := resolvePassword(true)
	if err != nil {
		return nil, "", err
	}

	// Read File
	encryptedData, err := os.ReadFile(vaultFile)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read vault: %w", err)
	}

	// Decrypt
	plaintext, err := crypto.Decrypt(encryptedData, password)
	if err != nil {
		return nil, "", err
	}

	return plaintext, password, nil
}
