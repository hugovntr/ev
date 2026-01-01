package crypto

import (
	"bytes"
	"strings"
	"testing"
)

func TestEncryptDecrypt_HappyPath(t *testing.T) {
	password := "my-secret-password"
	originalText := []byte("API_KEY=123456\nDB_HOST=localhost")

	// 1. Encrypt
	encryptedData, err := Encrypt(originalText, password)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// 2. Decrypt
	decryptedText, err := Decrypt(encryptedData, password)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// 3. Verify exact match
	if !bytes.Equal(originalText, decryptedText) {
		t.Errorf("Mismatch!\nExpected: %s\nGot: %s", originalText, decryptedText)
	}
}

func TestEncryptDecrypt_WrongPassword(t *testing.T) {
	password := "correct-password"
	wrongPassword := "wrong-password"
	originalText := []byte("sensitive-data")

	encryptedData, _ := Encrypt(originalText, password)

	// Attempt decrypt with wrong password
	_, err := Decrypt(encryptedData, wrongPassword)
	if err == nil {
		t.Fatal("Decryption should have failed with wrong password, but it succeeded")
	}

	// Check if the error message is somewhat descriptive (optional)
	if !strings.Contains(err.Error(), "decryption failed") {
		t.Errorf("Expected standard decryption error, got: %v", err)
	}
}

func TestEncryptDecrypt_TamperedData(t *testing.T) {
	password := "secure-password"
	originalText := []byte("integrity-check")

	encryptedData, _ := Encrypt(originalText, password)

	// Tamper with the data: change the last byte
	encryptedData[len(encryptedData)-1] ^= 0xFF

	_, err := Decrypt(encryptedData, password)
	if err == nil {
		t.Fatal("Decryption should have failed for tampered data, but it succeeded")
	}
}
