package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func getMachineID() string {
	var cmd *exec.Cmd
	var output []byte
	var err error

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
		output, err = cmd.Output()
		if err != nil {
			return "chiko-error-fallback-id-darwin"
		}
		for _, line := range strings.Split(string(output), "\n") {
			if strings.Contains(line, "IOPlatformUUID") {
				parts := strings.Split(line, "=")
				if len(parts) > 1 {
					return strings.Trim(parts[1], " \"")
				}
			}
		}
		return "chiko-error-fallback-id-darwin-uuid-not-found"
	case "linux":
		idPath := "/etc/machine-id"
		if _, err := os.Stat(idPath); os.IsNotExist(err) {
			idPath = "/var/lib/dbus/machine-id"
		}
		content, err := os.ReadFile(idPath)
		if err == nil {
			return strings.TrimSpace(string(content))
		}
		hostname, err := os.Hostname()
		if err == nil {
			return hostname
		}
		return "chiko-error-fallback-id-linux"
	case "windows":
		cmd = exec.Command("powershell", "(Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Cryptography').MachineGuid")
		output, err = cmd.Output()
		if err != nil {
			return "chiko-error-fallback-id-windows"
		}
		return strings.TrimSpace(string(output))
	default:
		return "chiko-default-fallback-id"
	}
}

// getEncryptionKey derives a 32-byte key from the given ID string.
func getEncryptionKey(id string) []byte {
	hash := sha256.Sum256([]byte(id + "chiko-salt"))
	return hash[:]
}

// encryptionKeys returns the list of keys to try, in order:
// first the current machine-id key, then fallback keys for history
// that was encrypted on a different machine or before reinstall.
func encryptionKeys() [][]byte {
	id := getMachineID()
	keys := [][]byte{getEncryptionKey(id)}

	// Add well-known fallback keys for history encrypted without a
	// stable machine-id (e.g. after OS reinstall, VM clone, or when
	// /etc/machine-id was missing at encrypt time).
	fallbacks := []string{
		"chiko-error-fallback-id-linux",
		"chiko-error-fallback-id-darwin",
		"chiko-error-fallback-id-windows",
		"chiko-default-fallback-id",
	}
	for _, fb := range fallbacks {
		if fb != id {
			keys = append(keys, getEncryptionKey(fb))
		}
	}
	return keys
}

// Encrypt encrypts data using AES-GCM with the current machine-specific key.
func Encrypt(data []byte) ([]byte, error) {
	key := getEncryptionKey(getMachineID())
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// Decrypt decrypts data using AES-GCM. It tries the current machine-id
// key first, then well-known fallback keys, so history remains readable
// after OS reinstall or machine-id change.
func Decrypt(data []byte) ([]byte, error) {
	if len(data) < 12 { // minimum nonce size for AES-GCM
		return nil, errors.New("ciphertext too short")
	}

	var lastErr error
	for _, key := range encryptionKeys() {
		block, err := aes.NewCipher(key)
		if err != nil {
			lastErr = err
			continue
		}
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			lastErr = err
			continue
		}
		nonceSize := gcm.NonceSize()
		if len(data) < nonceSize {
			continue
		}
		nonce, ciphertext := data[:nonceSize], data[nonceSize:]
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err == nil {
			return plaintext, nil
		}
		lastErr = err
	}
	return nil, lastErr
}
