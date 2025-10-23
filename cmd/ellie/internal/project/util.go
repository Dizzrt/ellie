package project

import (
	"os"

	"golang.org/x/mod/modfile"
)

// modulePath returns go module path.
func modulePath(filename string) (string, error) {
	modBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return modfile.ModulePath(modBytes), nil
}
