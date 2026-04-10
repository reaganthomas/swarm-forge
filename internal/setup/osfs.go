package setup

import "os"

// OSFS implements FS using the real operating system.
type OSFS struct{}

// MkdirAll creates directories using os.MkdirAll.
func (OSFS) MkdirAll(path string, perm uint32) error {
	return os.MkdirAll(path, os.FileMode(perm))
}

// WriteFile writes data to a file using os.WriteFile.
func (OSFS) WriteFile(path string, data []byte, perm uint32) error {
	return os.WriteFile(path, data, os.FileMode(perm))
}

// ReadFile reads a file using os.ReadFile.
func (OSFS) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// Stat returns whether a file or directory exists.
func (OSFS) Stat(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
