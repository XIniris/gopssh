package file

import (
	"gopssh/log"
	"os"
)

func OpenFile(path string) (*os.File, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			log.Error("file %s not found, error: %v", path, err)
		}

		log.Error("failed to get file %s status, error: %v", path, err)
		return nil, err
	}

	return os.Open(path)
}

func SaveStringAsFile(path string, content string) error {
	f, err := os.Create(path)
	if err != nil {
		log.Error("failed to create file %s, error: %v", path, err)
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		log.Error("failed to write string to file %s, error: %v", path, err)
		return err
	}

	return nil
}

func SaveBytesAsFile(path string, content []byte) error {
	f, err := os.Create(path)
	if err != nil {
		log.Error("failed to create file %s, error: %v", path, err)
		return err
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		log.Error("failed to write string to file %s, error: %v", path, err)
		return err
	}

	return nil
}

func EnsureDirExist(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path, 0755); err != nil {
				log.Error("failed to create dir %s, error: %v", path, err)
				return err
			}
		} else {
			log.Error("failed to get dir %s status, error: %v", path, err)
			return err
		}
	}

	return nil
}