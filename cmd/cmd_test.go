package cmd

import (
	"bytes"
	"ev/crypto"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper to reset flags and buffers between tests
func executeCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)

	err := rootCmd.Execute()
	return buf.String(), err
}

func TestIntegrationFlow(t *testing.T) {
	// Setup Temporary Sandbox
	tmpDir, err := os.MkdirTemp("", "ev-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir) // Clean up after test

	// Define Paths
	testVault := filepath.Join(tmpDir, "vault")
	testCreds := filepath.Join(tmpDir, "credentials")
	testPassword := "secure-test-password"

	// Manually Create Credentials File
	// We do this to bypass the interactive terminal prompt of 'ev create'
	if err := os.WriteFile(testCreds, []byte(testPassword), 0600); err != nil {
		t.Fatalf("Failed to write creds: %v", err)
	}

	// Manually Create Initial Vault (Mocking 'ev create' internals to avoid stdin prompt)
	// We inject some test data directly
	initialData := []byte("API_KEY: 12345\nDB_HOST: localhost\n")
	encryptedData, err := crypto.Encrypt(initialData, testPassword)
	if err != nil {
		t.Fatalf("Failed to encrypt setup data: %v", err)
	}
	if err := os.WriteFile(testVault, encryptedData, 0600); err != nil {
		t.Fatalf("Failed to write vault: %v", err)
	}

	// TEST: 'ev get'
	// We must pass the flags explicitly because the global var defaults might be ~/.config/...
	// NOTE: Cobra flags persist, so we pass them in the args
	output, err := executeCommand("get", "API_KEY", "-f", testVault, "-p", testCreds)
	if err != nil {
		t.Fatalf("ev get failed: %v", err)
	}
	if output != "12345" {
		t.Errorf("Expected '12345', got '%s'", output)
	}

	// TEST: 'ev export'
	output, err = executeCommand("export", "-f", testVault, "-p", testCreds)
	if err != nil {
		t.Fatalf("ev export failed: %v", err)
	}

	expectedExport := `export API_KEY="12345"`
	if !strings.Contains(output, expectedExport) {
		t.Errorf("Export output missing API_KEY. Got:\n%s", output)
	}
}
