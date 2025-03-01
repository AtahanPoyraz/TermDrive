package fsutil

import (
	"errors"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"time"
)

// ExtractMimeType extracts the MIME type of a file based on its extension.
// It returns the MIME type as a string or "unknown" if no valid MIME type can be determined.
// Returns an error if the file does not exist or if there's an issue checking the file.
func ExtractMimeType(path string) (string, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file not found: %v", path)
		}

		return "", fmt.Errorf("unable to check file: %v", err)
	}

	if mimeType := mime.TypeByExtension(filepath.Ext(path)); mimeType != "" {
		return mimeType, nil
	}

	return "unknown", nil
}

// ExtractFileSize returns the size of the file in bytes.
// Returns an error if the file does not exist or if there's an issue retrieving the file size.
func ExtractFileSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, errors.New("file does not exist")
		}
		return 0, err
	}
	return fileInfo.Size(), nil
}

// ExtractFilePermissions retrieves the file's permissions in a human-readable format (e.g., rwxr-xr-x).
// Returns an error if the file does not exist or if there's an issue retrieving the file permissions.
func ExtractFilePermissions(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New("file does not exist")
		}
		return "", err
	}

	perms := "-"
	if info.IsDir() {
		perms = "d"
	}

	permBits := []rune{'-', '-', '-', '-', '-', '-', '-', '-', '-'}
	permissions := []os.FileMode{0400, 0200, 0100, 0040, 0020, 0010, 0004, 0002, 0001}
	chars := []rune{'r', 'w', 'x'}

	for i, bit := range permissions {
		if info.Mode()&bit != 0 {
			permBits[i] = chars[i%3]
		}
	}

	return perms + string(permBits), nil
}

// ExtractFileLastModified retrieves the last modified time of the file.
// Returns an error if the file does not exist or if there's an issue retrieving the modification time.
func ExtractFileLastModified(path string) (time.Time, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return time.Time{}, errors.New("file does not exist")
		}
		return time.Time{}, err
	}

	return fileInfo.ModTime(), nil
}
