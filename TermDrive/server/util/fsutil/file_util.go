package fsutil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/AtahanPoyraz/TermDrive/server/config"
	"github.com/AtahanPoyraz/TermDrive/server/internal/model"
)

// Create creates a file at the given path. If the parent directory doesn't exist, it will be created.
// Returns an error if there's any issue during directory or file creation.
func Create(path string) error {
	if filepath.Ext(path) == "" {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
		return nil
	}

	dirPath := filepath.Dir(path)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	return nil
}

// Delete removes a file or directory at the specified path. If it's a directory, it will be deleted recursively.
// Returns an error if the path cannot be accessed or deleted.
func Delete(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, &os.PathError{}) {
			return err
		}

		return fmt.Errorf("failed to access the path: %v", err)
	}

	if !info.IsDir() {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("failed to delete file: %v", err)
		}
		return nil
	}

	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to delete directory: %v", err)
	}

	return nil
}

// HumanReadableFileSize converts a byte size to a human-readable format (e.g., KB, MB, GB).
// Returns the human-readable size string (e.g., "12.34 MB").
func HumanReadableFileSize(size int64) string {
	const (
		_        = iota
		KB int64 = 1 << (10 * iota)
		MB
		GB
		TB
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	}

	return fmt.Sprintf("%d B", size)
}

// List retrieves and formats a list of files in the specified directory. It returns the file/directory permissions,
// owner, size, modification time, and name. It also returns an error if there's an issue reading the directory or
// retrieving file info.
func List(ctx context.Context, path string) ([]string, error) {
	user, ok := ctx.Value(config.UserContextKey).(*model.UserModel)
	if !ok || user == nil {
		return nil, fmt.Errorf("user not found or failed to retrieve. Context: %v", ctx.Err())
	}

	content, err := os.ReadDir(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, os.ErrNotExist
		}

		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	var contents []string
	for _, entity := range content {
		fileInfo, err := os.Stat(filepath.Join(path, entity.Name()))
		if err != nil {
			if errors.Is(err, &os.PathError{}) {
				return nil, err
			}

			return nil, fmt.Errorf("error getting file info: %v", err)
		}

		contents = append(contents,
			fmt.Sprintf("%-10s %s %s %-10s %-10s %s",
				fileInfo.Mode(),
				user.Username,
				user.Role,
				HumanReadableFileSize(fileInfo.Size()),
				fileInfo.ModTime().Format("2006-01-02 15:04:05"),
				fileInfo.Name(),
			),
		)
	}

	return contents, nil
}

// WriteFile writes the contents of a file to the specified destination path. The function supports context timeout.
// Returns an error if the file can't be written or if the operation exceeds the timeout.
func WriteFile(ctx context.Context, file multipart.File, path string, bufferSize int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*1)
	defer cancel()

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	outFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()

	go func() {
		defer pipeWriter.Close()
		buffer := make([]byte, bufferSize)
		for {
			select {
			case <-ctx.Done():
				pipeWriter.CloseWithError(fmt.Errorf("operation timed out"))
				return

			default:
				n, err := file.Read(buffer)
				if err != nil && err != io.EOF {
					pipeWriter.CloseWithError(fmt.Errorf("failed to read file: %v", err))
					return
				}

				if n == 0 {
					pipeWriter.Close()
					return
				}

				if _, err := pipeWriter.Write(buffer[:n]); err != nil {
					pipeWriter.CloseWithError(fmt.Errorf("failed to write to pipe: %v", err))
					return
				}
			}
		}
	}()

	_, err = io.Copy(outFile, pipeReader)
	if err != nil {
		return fmt.Errorf("failed to copy data to file: %v", err)
	}

	return nil
}

func ReadFile(ctx context.Context, path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, os.ErrNotExist
		}

		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return file, nil
}
