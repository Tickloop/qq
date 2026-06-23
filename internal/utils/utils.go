package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func ResolveToAbsPath(path string) (string, error) {
	// This function does three resolutions:
	// 1. ~ to /home/<user>
	// 2. Expanding $USER or $VARS present in the input
	// 3. Converting to an absolute Path

	// 1. ~ to /home/<user>
	if path == "~" || strings.HasPrefix(path, "~/") {
		base, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		path = filepath.Join(base, strings.TrimPrefix(path, "~"))
	}

	// 2. Expand $VARS in input
	path = os.ExpandEnv(path)

	// 3. Convert to abs path
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return path, nil
}
