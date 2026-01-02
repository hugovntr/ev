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
	if err := os.WriteFile(testCreds, []byte(testPassword), 0600); err != nil {
		t.Fatalf("Failed to write creds: %v", err)
	}

	// Manually Create Initial Vault
	// We inject some test data: "API_KEY" (string) and "PORT" (int) to test type conversion
	initialData := []byte("API_KEY: 12345\nPORT: 8080\n")
	encryptedData, err := crypto.Encrypt(initialData, testPassword)
	if err != nil {
		t.Fatalf("Failed to encrypt setup data: %v", err)
	}
	if err := os.WriteFile(testVault, encryptedData, 0600); err != nil {
		t.Fatalf("Failed to write vault: %v", err)
	}

	// TEST: 'ev get'
	output, err := executeCommand("get", "API_KEY", "-f", testVault, "-p", testCreds)
	if err != nil {
		t.Fatalf("ev get failed: %v", err)
	}
	// Check exact value (no weird characters or newlines)
	if output != "12345" {
		t.Errorf("Expected '12345', got '%s'", output)
	}

	// TEST: 'ev env' (Default / Bash)
	output, err = executeCommand("env", "-f", testVault, "-p", testCreds)
	if err != nil {
		t.Fatalf("ev env (default) failed: %v", err)
	}
	expectedBash := `export API_KEY="12345"`
	if !strings.Contains(output, expectedBash) {
		t.Errorf("Bash env output missing API_KEY.\nExpected: %s\nGot:\n%s", expectedBash, output)
	}

	// TEST: 'ev env fish'
	output, err = executeCommand("env", "fish", "-f", testVault, "-p", testCreds)
	if err != nil {
		t.Fatalf("ev env fish failed: %v", err)
	}
	expectedFish := `set -gx API_KEY "12345"`
	if !strings.Contains(output, expectedFish) {
		t.Errorf("Fish env output incorrect.\nExpected: %s\nGot:\n%s", expectedFish, output)
	}

	// TEST: 'ev env pwsh'
	output, err = executeCommand("env", "pwsh", "-f", testVault, "-p", testCreds)
	if err != nil {
		t.Fatalf("ev env pwsh failed: %v", err)
	}
	expectedPwsh := `$Env:API_KEY = "12345"`
	if !strings.Contains(output, expectedPwsh) {
		t.Errorf("PowerShell env output incorrect.\nExpected: %s\nGot:\n%s", expectedPwsh, output)
	}
}
