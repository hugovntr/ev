package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

const (
	saltSize   = 16
	nonceSize  = 12
	keyIter    = 100000 // PBKDF2 iterations
	keySize    = 32     // AES-256
	lineLength = 64     // Wrap lines at 64 chars for readability
)

// Encrypt plaintext using the provided password.
// It returns a byte slice containing: Salt + Nonce + Ciphertext
func Encrypt(plaintext []byte, password string) ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	key, err := deriveKey(password, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	binary := make([]byte, saltSize+nonceSize+len(ciphertext))
	copy(binary[0:saltSize], salt)
	copy(binary[saltSize:saltSize+nonceSize], nonce)
	copy(binary[saltSize+nonceSize:], ciphertext)

	return encodeWithWrapping(binary), nil
}

// Decrypt the data using the provided password.
func Decrypt(data []byte, password string) ([]byte, error) {
	binary, err := decodeWithUnwrapping(data)
	if err != nil {
		return nil, errors.New("invalid vault format")
	}

	if len(binary) < saltSize+nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	salt := binary[:saltSize]
	nonce := binary[saltSize : saltSize+nonceSize]
	ciphertext := binary[saltSize+nonceSize:]

	key, err := deriveKey(password, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("decryption failed: incorrect password or corrupted data")
	}

	return plaintext, nil
}

// deriveKey generates a 32-byte key from a password and salt
func deriveKey(password string, salt []byte) ([]byte, error) {
	return pbkdf2.Key(sha256.New, password, salt, keyIter, keySize)
}

// encodeWithWrapping encodes bytes to Base64 and inserts newlines every 64 chars
func encodeWithWrapping(input []byte) []byte {
	encoded := base64.StdEncoding.EncodeToString(input)
	var buffer bytes.Buffer

	for i, r := range encoded {
		buffer.WriteRune(r)
		if (i+1)%lineLength == 0 && (i+1) < len(encoded) {
			buffer.WriteByte('\n')
		}
	}
	buffer.WriteByte('\n')
	return buffer.Bytes()
}

// decodeWithUnwrapping strips newlines and decodes Base64
func decodeWithUnwrapping(input []byte) ([]byte, error) {
	// Filter out newlines (0x0A) and carriage returns (0x0D)
	filtered := bytes.Map(func(r rune) rune {
		if r == '\n' || r == '\r' {
			return -1 // Drop char
		}
		return r
	}, input)

	return base64.StdEncoding.DecodeString(string(filtered))
}
