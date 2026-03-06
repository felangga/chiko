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
		// MacOS: IOPlatformUUID
		cmd = exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
		output, err = cmd.Output()
		if err != nil {
			return "chiko-error-fallback-id-darwin"
		}
		// Extract UUID from ioreg output
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
		// Linux: machine-id
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
		// Windows: MachineGuid
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

// getEncryptionKey derives a 32-byte key from the machine ID
func getEncryptionKey() []byte {
	id := getMachineID()
	hash := sha256.Sum256([]byte(id + "chiko-salt"))
	return hash[:]
}

// Encrypt encrypts data using AES-GCM with a machine-specific key
func Encrypt(data []byte) ([]byte, error) {
	key := getEncryptionKey()
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

// Decrypt decrypts data using AES-GCM with a machine-specific key
func Decrypt(data []byte) ([]byte, error) {
	key := getEncryptionKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
