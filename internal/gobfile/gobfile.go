package gobfile

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

func getAppDataPath(appName string) (string, error) {
	if !isValidDirName(appName) {
		return "", fmt.Errorf("invalid directory name: %s", appName)
	}

	var appDataDir string
	switch runtime.GOOS {
	case "windows":
		appDataDir = os.Getenv("LOCALAPPDATA")
	case "darwin":
		appDataDir = os.Getenv("HOME")
		appDataDir = filepath.Join(appDataDir, "Library", "Application Support")
	case "linux", "freebsd", "netbsd", "openbsd", "dragonfly":
		appDataDir = os.Getenv("HOME")
		appDataDir = filepath.Join(appDataDir, ".local", "share")
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	appDataPath := filepath.Join(appDataDir, appName)
	if _, err := os.Stat(appDataPath); os.IsNotExist(err) {
		if err := os.MkdirAll(appDataPath, 0755); err != nil {
			return "", err
		}
	}

	return appDataPath, nil
}

func isValidDirName(name string) bool {
	invalidCharsRegex := regexp.MustCompile(`[<>:"/\|?*\x00-\x1F]`)
	return !invalidCharsRegex.MatchString(name)
}

func encodeFile(filePath string, data any) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func decodeFile(filePath string, data any) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(bufio.NewReader(file))
	if err := decoder.Decode(data); err != nil {
		return err
	}

	return nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func createEmptyFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func ReadData[T any](appName string, fileName string) (T, error) {
	var decodedData T

	appDataPath, err := getAppDataPath(appName)
	if err != nil {
		return decodedData, fmt.Errorf("error getting app data path: %w", err)
	}

	filePath := filepath.Join(appDataPath, fmt.Sprintf("%s.gob", fileName))

	if !fileExists(filePath) {
		if err := createEmptyFile(filePath); err != nil {
			return decodedData, fmt.Errorf("error creating file: %w", err)
		}
	}

	if err := decodeFile(filePath, &decodedData); err != nil {
		return decodedData, fmt.Errorf("error decoding file: %w", err)
	}

	return decodedData, nil
}

func WriteData[T any](appName string, fileName string, data T) error {
	appDataPath, err := getAppDataPath(appName)
	if err != nil {
		return fmt.Errorf("error getting app data path: %w", err)
	}

	filePath := filepath.Join(appDataPath, fmt.Sprintf("%s.gob", fileName))

	if err := encodeFile(filePath, data); err != nil {
		return fmt.Errorf("error encoding file: %w", err)
	}

	return nil
}
