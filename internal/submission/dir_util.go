package submission

import "os"

func ensureDirectoryExists(dir string) error {
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
